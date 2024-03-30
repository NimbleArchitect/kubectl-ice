package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	apires "k8s.io/apimachinery/pkg/api/resource"
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

	stdinChanged, err := builder.HasStdinChanged()
	if err != nil {
		return err
	}

	//only need to pull metrics info we are reading live data,
	// if we read from a file metric data wont exist
	if len(commonFlagList.inputFilename) == 0 && !stdinChanged {
		if err := connect.LoadMetricConfig(kubeFlags); err != nil {
			return err
		}
		podStateList, err := connect.GetMetricPods(args)
		if err != nil {
			log.Tell(err)
		} else {
			loopinfo.MetricsResource = loopinfo.podMetrics2Hashtable(podStateList)
		}
	}

	if cmd.Flag("size") != nil {
		if len(cmd.Flag("size").Value.String()) > 0 {
			loopinfo.BytesAs = cmd.Flag("size").Value.String()
		}
	}

	if cmd.Flag("raw").Value.String() == "true" {
		loopinfo.ShowRaw = true
		loopinfo.BytesAs = "M"
	}

	table := Table{}
	table.ColourOutput = commonFlagList.outputAsColour
	table.CustomColours = commonFlagList.useTheseColours

	builder.Table = &table
	builder.ShowTreeView = commonFlagList.showTreeView

	if err := builder.Build(&loopinfo); err != nil {
		return err
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	// do we need to find the outliers, we have enough data to compute a range
	if commonFlagList.showOddities {
		row2Remove, err := table.ListOutOfRange(builder.DefaultHeaderLen) //1 = used column
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

func (s *resource) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	rowOut := make([]Cell, 5)

	for _, r := range rows {
		// "USED", "REQUEST", "LIMIT", "%REQ", "%LIMIT",
		rowOut[0].number += r[0].number
		rowOut[1].number += r[1].number
		rowOut[2].number += r[2].number
	}

	floatfmt := "%.6f"
	typefmt := "%d"
	if s.ResourceType == "cpu" {
		typefmt = "%dm"
	}
	if !s.ShowRaw {
		floatfmt = "%.2f"
	}

	if s.ResourceType == "memory" {
		// everything is stored internally as kb so we need to * 1000 to get back to bytes
		if s.ShowRaw {
			typefmt = "%dk"
			rowOut[0].text = fmt.Sprintf(typefmt, rowOut[0].number)
		} else {
			rowOut[0].text = memoryHumanReadable(rowOut[0].number*1000, s.BytesAs)
		}
		rowOut[1].text = memoryHumanReadable(rowOut[1].number, s.BytesAs)
		rowOut[2].text = memoryHumanReadable(rowOut[2].number, s.BytesAs)
	} else {
		if s.ShowRaw {
			rowOut[0].text = fmt.Sprintf("%dn", rowOut[0].number)
		} else {
			rowOut[0].text = fmt.Sprintf(typefmt, rowOut[0].number)
		}
		rowOut[1].text = fmt.Sprintf(typefmt, rowOut[1].number)
		rowOut[2].text = fmt.Sprintf(typefmt, rowOut[2].number)
	}

	if rowOut[0].number > 0 {
		if rowOut[1].number > 0.0 {
			// calc % request
			val := validateFloat64(float64(rowOut[0].number) / float64(rowOut[1].number) * 100)
			rowOut[4].text = fmt.Sprintf(floatfmt, val)
			rowOut[4].float = val
			rowOut[4].colour = setColourValue(int(val))
		}

		if rowOut[2].number > 0.0 {
			// calc % limit
			val := validateFloat64(float64(rowOut[0].number) / float64(rowOut[2].number) * 100)
			rowOut[3].text = fmt.Sprintf(floatfmt, val)
			rowOut[3].float = val
			rowOut[3].colour = setColourValue(int(val))
		}

		usedColour := [2]int{0, 0}
		if rowOut[3].float > rowOut[4].float {
			usedColour = setColourValue(int(rowOut[4].float))
		} else {
			usedColour = setColourValue(int(rowOut[3].float))
		}

		rowOut[0].colour = usedColour
	}

	return rowOut, nil
}

func (s *resource) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	metrics := s.MetricsResource[info.PodName][info.Name]
	out := make([][]Cell, 1)
	out[0] = s.statsProcessTableRow(container.Resources, metrics, info, s.ResourceType)
	return out, nil
}

func (s *resource) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	metrics := s.MetricsResource[info.PodName][info.Name]
	out := make([][]Cell, 1)
	out[0] = s.statsProcessTableRow(container.Resources, metrics, info, s.ResourceType)
	return out, nil
}

func (s *resource) statsProcessTableRow(res v1.ResourceRequirements, metrics v1.ResourceList, info BuilderInformation, resource string) []Cell {
	var cellList []Cell
	var displayValue, request, limit, percentLimit, percentRequest string
	var rawRequest, rawLimit, rawValue int64
	var rawPercentRequest, rawPercentLimit float64
	var requestCell, limitCell Cell

	log := logger{location: "resources:statsProcessTableRow"}
	log.Debug("Start")

	percentRequestColour := [2]int{-1, 0}
	percentLimitColour := [2]int{-1, 0}
	floatfmt := "%.6f"

	if resource == "cpu" {
		if res.Size() >= 3 {
			if res.Limits.Cpu() != nil {
				if s.ShowRaw {
					rawLimit = res.Limits.Cpu().ScaledValue(apires.Nano)
					limit = fmt.Sprintf("%dn", rawLimit)
				} else {
					rawLimit = res.Limits.Cpu().MilliValue()
					limit = fmt.Sprintf("%dm", rawLimit)
				}
				limitCell = NewCellInt(limit, rawLimit)
			}

			if res.Requests.Cpu() != nil {
				if s.ShowRaw {
					rawRequest = res.Requests.Cpu().ScaledValue(apires.Nano)
					request = fmt.Sprintf("%dn", rawRequest)
				} else {
					rawRequest = res.Requests.Cpu().MilliValue()
					request = fmt.Sprintf("%dm", rawRequest)
				}
				requestCell = NewCellInt(request, rawRequest)
			}
		}
		if metrics.Cpu() != nil {
			if s.ShowRaw {
				// this returns nanocores as the display value when using --raw
				displayValue = metrics.Cpu().String()
				rawValue = metrics.Cpu().ScaledValue(apires.Nano)
			} else {
				floatfmt = "%.2f"
				displayValue = fmt.Sprintf("%dm", metrics.Cpu().MilliValue())
				rawValue = metrics.Cpu().MilliValue()
			}

			if cpuVal := metrics.Cpu().AsApproximateFloat64(); cpuVal > 0 {
				// check cpu limits has a value
				if res.Limits.Cpu().AsApproximateFloat64() == 0 {
					percentLimit = "-"
					rawPercentLimit = 0.0
				} else {
					val := validateFloat64(cpuVal / res.Limits.Cpu().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
					rawPercentLimit = val

					percentLimitColour = setColourValue(int(val))
				}
				// check cpu requests has a value
				if res.Requests.Cpu().AsApproximateFloat64() == 0 {
					percentRequest = "-"
					rawPercentRequest = 0.0
				} else {
					val := validateFloat64(cpuVal / res.Requests.Cpu().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
					rawPercentRequest = val

					percentRequestColour = setColourValue(int(val))
				}
			}
		}

	}

	if resource == "memory" {
		if res.Size() >= 3 {
			if res.Limits.Memory() != nil {
				limit = res.Limits.Memory().String()
				rawLimit = res.Limits.Memory().Value()
				limitCell = NewCellInt(limit, rawLimit)
			}

			if res.Requests.Memory() != nil {
				request = res.Requests.Memory().String()
				rawRequest = res.Requests.Memory().Value()
				requestCell = NewCellInt(request, rawRequest)
			}
		}
		if metrics.Memory() != nil {
			rawValue = metrics.Memory().Value() / 1000
			if s.ShowRaw {
				displayValue = fmt.Sprintf("%dk", metrics.Memory().Value())
			} else {
				displayValue = memoryHumanReadable(metrics.Memory().Value(), s.BytesAs)
				floatfmt = "%.2f"
			}

			if memVal := metrics.Memory().AsApproximateFloat64(); memVal > 0 {
				// check memory limits has a value
				if res.Limits.Memory().AsApproximateFloat64() == 0 {
					percentLimit = "-"
					rawPercentLimit = 0.0
					percentLimitColour = [2]int{-1, 0}
				} else {
					val := validateFloat64(memVal / res.Limits.Memory().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
					rawPercentLimit = val

					percentLimitColour = setColourValue(int(val))
				}
				// check memory requests has a value
				if res.Requests.Memory().AsApproximateFloat64() == 0 {
					percentRequest = "-"
					rawPercentRequest = 0.0
					percentRequestColour = [2]int{-1, 0}
				} else {
					val := validateFloat64(memVal / res.Requests.Memory().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
					rawPercentRequest = val

					percentRequestColour = setColourValue(int(val))
				}
			}
		}
	}

	usedColour := [2]int{-1, 0}
	if percentLimitColour[0] > percentRequestColour[0] {
		usedColour = percentLimitColour
	} else {
		usedColour = percentRequestColour
	}

	cellList = append(cellList,
		NewCellColourInt(usedColour, displayValue, rawValue),
		requestCell,
		limitCell,
		NewCellColourFloat(percentRequestColour, percentRequest, rawPercentRequest),
		NewCellColourFloat(percentLimitColour, percentLimit, rawPercentLimit),
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

func (s *resource) BuildPodRow(pod v1.Pod, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}
