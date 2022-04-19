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
	var podname []string
	var showPodName bool = true
	var showRaw bool

	clientset, err := loadConfig(kubeFlags)
	if err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			showPodName = false
		}
	}

	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}

	podList, err := getPods(clientset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	metricset, err := loadMetricConfig(kubeFlags)
	if err != nil {
		return err
	}
	podStateList, err := getMetricPods(metricset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	if cmd.Flag("raw").Value.String() == "true" {
		showRaw = true
	}

	table := Table{}
	table.SetHeader(
		"T", "PODNAME", "CONTAINER", "USED", "REQUEST", "LIMIT", "%REQ", "%LIMIT",
	)

	if len(commonFlagList.filterList) >= 1 {
		err = table.SetFilter(commonFlagList.filterList)
		if err != nil {
			return err
		}
	}

	if !showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(1)
	}

	podState := podMetrics2Hashtable(podStateList)
	for _, pod := range podList {
		// process init containers
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := statsProcessTableRow(container, podState[pod.Name][container.Name], pod.Name, "I", resourceType, showRaw)
			table.AddRow(tblOut...)
		}

		// process standard containers
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := statsProcessTableRow(container, podState[pod.Name][container.Name], pod.Name, "S", resourceType, showRaw)
			table.AddRow(tblOut...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	// at this point we have data on all containers

	// do we need to find the outliers, we have enough data to compute a range
	if commonFlagList.showOddities {
		row2Remove, err := table.ListOutOfRange(3, table.GetRows()) //3 = used column
		if err != nil {
			return err
		}
		table.HideRows(row2Remove)
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil
}

func statsProcessTableRow(container v1.Container, metrics v1.ResourceList, podName, containerType string, resource string, showRaw bool) []Cell {
	var displayValue, request, limit, percentLimit, percentRequest string
	var rawRequest, rawLimit, rawValue int64
	var rawPercentRequest, rawPercentLimit float64

	floatfmt := "%.6f"

	if resource == "cpu" {
		if metrics.Cpu() != nil {
			rawValue = metrics.Cpu().MilliValue()
			if showRaw {
				displayValue = metrics.Cpu().String()
			} else {
				displayValue = fmt.Sprintf("%d", metrics.Cpu().MilliValue())
				floatfmt = "%.2f"
			}

			limit = container.Resources.Limits.Cpu().String()
			rawLimit = container.Resources.Limits.Cpu().Value() //* 1000
			// limit = fmt.Sprintf("%d", rawLimit)
			request = container.Resources.Requests.Cpu().String()
			rawRequest = container.Resources.Requests.Cpu().Value() //* 1000
			// request = fmt.Sprintf("%d", rawRequest)

			if cpuVal := metrics.Cpu().AsApproximateFloat64(); cpuVal > 0 {
				// check cpu limits has a value
				if container.Resources.Limits.Cpu().AsApproximateFloat64() == 0 {
					percentLimit = "-"
					rawPercentLimit = 0.0
				} else {
					val := validateFloat64(cpuVal / container.Resources.Limits.Cpu().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
					rawPercentLimit = val
				}
				// check cpu requests has a value
				if container.Resources.Requests.Cpu().AsApproximateFloat64() == 0 {
					percentRequest = "-"
					rawPercentRequest = 0.0
				} else {
					val := validateFloat64(cpuVal / container.Resources.Requests.Cpu().AsApproximateFloat64() * 100)
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
				displayValue = memoryHumanReadable(metrics.Memory().Value())
				floatfmt = "%.2f"
			}

			limit = container.Resources.Limits.Memory().String()
			rawLimit = container.Resources.Limits.Memory().Value()
			// limit = fmt.Sprintf("%d", rawLimit)
			request = container.Resources.Requests.Memory().String()
			rawRequest = container.Resources.Requests.Memory().Value()
			// request = fmt.Sprintf("%d", rawRequest)

			if memVal := metrics.Memory().AsApproximateFloat64(); memVal > 0 {
				// check memory limits has a value
				if container.Resources.Limits.Memory().AsApproximateFloat64() == 0 {
					percentLimit = "-"
					rawPercentLimit = 0.0
				} else {
					val := validateFloat64(memVal / container.Resources.Limits.Memory().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
					rawPercentLimit = val
				}
				// check memory requests has a value
				if container.Resources.Requests.Memory().AsApproximateFloat64() == 0 {
					percentRequest = "-"
					rawPercentRequest = 0.0
				} else {
					val := validateFloat64(memVal / container.Resources.Requests.Memory().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
					rawPercentRequest = val
				}
			}
		}
	}

	return []Cell{
		NewCellText(containerType),
		NewCellText(podName),
		NewCellText(container.Name),
		NewCellInt(displayValue, rawValue),
		NewCellInt(request, rawRequest),
		NewCellInt(limit, rawLimit),
		NewCellFloat(percentRequest, rawPercentRequest),
		NewCellFloat(percentLimit, rawPercentLimit),
	}
}

func podMetrics2Hashtable(stateList []v1beta1.PodMetrics) map[string]map[string]v1.ResourceList {
	podState := make(map[string]map[string]v1.ResourceList)

	for _, pod := range stateList {
		podState[pod.Name] = make(map[string]v1.ResourceList)
		for _, container := range pod.Containers {
			podState[pod.Name][container.Name] = container.Usage
		}
	}
	return podState
}
