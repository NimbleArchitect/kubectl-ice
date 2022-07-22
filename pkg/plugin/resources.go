package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// returns a string replacing %[1] with the resourse type r
func resourceShort(r string) string {
	return fmt.Sprintf("Show configured %[1]s size, limit and %% usage of each container", r)
}

// returns a string replacing %[1] with the resourse type r
func resourceDescription(r string) string {
	return fmt.Sprintf(` Prints the current %[1]s usage along with configured requests and limits. The calculated %% fields
serve as an easy way to see how close you are to the configured sizes.  By specifying the -r 
flag you can see raw unfiltered values.  If no name is specified the container %[1]s details
of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`, r)
}

// returns a string replacing %[2] with the resourse type r
// %[1] is replaced with its self as it is needed later on
func resourceExample(r string) string {
	return fmt.Sprintf(`  # List containers %[2]s info from pods
  %[1]s %[2]s

  # List container %[2]s info from pods output in JSON format
  %[1]s %[2]s -o json

  # List container %[2]s info from a single pod
  %[1]s %[2]s my-pod-4jh36

  # List %[2]s info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s %[2]s -c web-container

  # List %[2]s info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s %[2]s -c web-container --sort '!CONTAINER'

  # List %[2]s info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s %[2]s -c web-container --sort PODNAME

  # List container %[2]s info from all pods where label app matches web
  %[1]s %[2]s -l app=web

  # List container %[2]s info from all pods where the pod label app is either web or mail
  %[1]s %[2]s -l "app in (web,mail)"`, "%[1]s", r)
}

func Resources(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string, resourceType string) error {

	log := logger{location: "Resource"}
	log.Debug("Start", resourceType)

	builder := RowBuilder{}
	builder.LoopSpec = true
	builder.ShowInitContainers = true
	builder.PodName = args

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList

	loopinfo := resource{}
	builder.Connection = &connect
	builder.SetFlagsFrom(commonFlagList)

	loopinfo.ResourceType = resourceType

	if err := connect.LoadMetricConfig(kubeFlags); err != nil {
		return err
	}
	podStateList, err := connect.GetMetricPods(args)
	if err != nil {
		return err
	}

	if cmd.Flag("raw").Value.String() == "true" {
		loopinfo.ShowRaw = true
	}

	if cmd.Flag("size") != nil {
		if len(cmd.Flag("size").Value.String()) > 0 {
			loopinfo.BytesAs = cmd.Flag("size").Value.String()
		}
	}

	table := Table{}
	builder.Table = &table
	builder.ShowTreeView = commonFlagList.showTreeView

	loopinfo.MetricsResource = loopinfo.podMetrics2Hashtable(podStateList)
	builder.Build(&loopinfo)

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	// do we need to find the outliers, we have enough data to compute a range
	if commonFlagList.showOddities {
		row2Remove, err := table.ListOutOfRange(builder.DefaultHeaderLen, table.GetRows()) //1 = used column
		if err != nil {
			return err
		}
		table.HideRows(row2Remove)
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil
}

type resource struct {
	MetricsResource map[string]map[string]v1.ResourceList
	ResourceType    string
	BytesAs         string
	ShowRaw         bool
	ShowPrevious    bool
	ShowDetails     bool
}

func (s *resource) Headers() []string {
	return []string{
		"USED", "REQUEST", "LIMIT", "%REQ", "%LIMIT",
	}
}

func (s *resource) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *resource) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s *resource) BuildBranch(info BuilderInformation, podList []v1.Pod) ([]Cell, error) {
	out := []Cell{
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}
	return out, nil
}

func (s *resource) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}, nil
}

func (s *resource) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	metrics := s.MetricsResource[info.PodName][info.ContainerName]
	out := make([][]Cell, 1)
	out[0] = s.statsProcessTableRow(container.Resources, metrics, info, s.ResourceType, s.ShowRaw, s.BytesAs)
	return out, nil
}

func (s *resource) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	metrics := s.MetricsResource[info.PodName][info.ContainerName]
	out := make([][]Cell, 1)
	out[0] = s.statsProcessTableRow(container.Resources, metrics, info, s.ResourceType, s.ShowRaw, s.BytesAs)
	return out, nil
}

func (s *resource) Sum(rows [][]Cell) []Cell {
	rowOut := make([]Cell, 5)
	return rowOut
}

func (s *resource) statsProcessTableRow(res v1.ResourceRequirements, metrics v1.ResourceList, info BuilderInformation, resource string, showRaw bool, bytesAs string) []Cell {
	var cellList []Cell
	var displayValue, request, limit, percentLimit, percentRequest string
	var rawRequest, rawLimit, rawValue int64
	var rawPercentRequest, rawPercentLimit float64

	log := logger{location: "resources:statsProcessTableRow"}
	log.Debug("Start")

	floatfmt := "%.6f"

	if resource == "cpu" {
		if metrics.Cpu() != nil {
			rawValue = metrics.Cpu().MilliValue()
			if showRaw {
				displayValue = metrics.Cpu().String()
			} else {
				displayValue = fmt.Sprintf("%dm", metrics.Cpu().MilliValue())
				floatfmt = "%.2f"
			}

			rawLimit = res.Limits.Cpu().MilliValue()
			limit = fmt.Sprintf("%dm", rawLimit)
			rawRequest = res.Requests.Cpu().MilliValue()
			request = fmt.Sprintf("%dm", rawRequest)

			if cpuVal := metrics.Cpu().AsApproximateFloat64(); cpuVal > 0 {
				// check cpu limits has a value
				if res.Limits.Cpu().AsApproximateFloat64() == 0 {
					percentLimit = "-"
					rawPercentLimit = 0.0
				} else {
					val := validateFloat64(cpuVal / res.Limits.Cpu().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
					rawPercentLimit = val
				}
				// check cpu requests has a value
				if res.Requests.Cpu().AsApproximateFloat64() == 0 {
					percentRequest = "-"
					rawPercentRequest = 0.0
				} else {
					val := validateFloat64(cpuVal / res.Requests.Cpu().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
					rawPercentRequest = val
				}
			}
		}

	}

	if resource == "memory" {
		if metrics.Memory() != nil {
			rawValue = metrics.Memory().Value() / 1000
			if showRaw {
				displayValue = fmt.Sprintf("%d", metrics.Memory().Value())
			} else {
				displayValue = memoryHumanReadable(metrics.Memory().Value(), bytesAs)
				floatfmt = "%.2f"
			}

			limit = res.Limits.Memory().String()
			rawLimit = res.Limits.Memory().Value()
			request = res.Requests.Memory().String()
			rawRequest = res.Requests.Memory().Value()

			if memVal := metrics.Memory().AsApproximateFloat64(); memVal > 0 {
				// check memory limits has a value
				if res.Limits.Memory().AsApproximateFloat64() == 0 {
					percentLimit = "-"
					rawPercentLimit = 0.0
				} else {
					val := validateFloat64(memVal / res.Limits.Memory().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
					rawPercentLimit = val
				}
				// check memory requests has a value
				if res.Requests.Memory().AsApproximateFloat64() == 0 {
					percentRequest = "-"
					rawPercentRequest = 0.0
				} else {
					val := validateFloat64(memVal / res.Requests.Memory().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
					rawPercentRequest = val
				}
			}
		}
	}

	cellList = append(cellList,
		NewCellInt(displayValue, rawValue),
		NewCellInt(request, rawRequest),
		NewCellInt(limit, rawLimit),
		NewCellFloat(percentRequest, rawPercentRequest),
		NewCellFloat(percentLimit, rawPercentLimit),
	)

	log.Debug("cellList", cellList)
	return cellList
}

func (s *resource) podMetrics2Hashtable(stateList []v1beta1.PodMetrics) map[string]map[string]v1.ResourceList {
	podState := make(map[string]map[string]v1.ResourceList)

	for _, pod := range stateList {
		podState[pod.Name] = make(map[string]v1.ResourceList)
		for _, container := range pod.Containers {
			podState[pod.Name][container.Name] = container.Usage
		}
	}
	return podState
}
