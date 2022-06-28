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
	var columnInfo containerInfomation
	var tblHead []string
	var podname []string
	var showPodName bool = true
	var showRaw bool
	var nodeLabels map[string]map[string]string
	var podLabels map[string]map[string]string

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
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
	connect.Flags = commonFlagList

	podList, err := connect.GetPods(podname)
	if err != nil {
		return err
	}

	if err := connect.LoadMetricConfig(kubeFlags); err != nil {
		return err
	}
	podStateList, err := connect.GetMetricPods(podname)
	if err != nil {
		return err
	}

	if cmd.Flag("raw").Value.String() == "true" {
		showRaw = true
	}

	if cmd.Flag("node-label").Value.String() != "" {
		columnInfo.labelNodeName = cmd.Flag("node-label").Value.String()
		nodeLabels, err = connect.GetNodeLabels(podList)
		if err != nil {
			return err
		}
	}

	if cmd.Flag("pod-label").Value.String() != "" {
		columnInfo.labelPodName = cmd.Flag("pod-label").Value.String()
		podLabels, err = connect.GetPodLabels(podList)
		if err != nil {
			return err
		}
	}

	table := Table{}
	columnInfo.treeView = commonFlagList.showTreeView

	tblHead = columnInfo.GetDefaultHead()
	defaultHeaderLen := len(tblHead)
	if commonFlagList.showTreeView {
		// we have to control the name when displaying a tree view as the table
		//  object dosent have the extra info to be able to process it
		tblHead = append(tblHead, "NAME")
	}

	tblHead = append(tblHead, "USED", "REQUEST", "LIMIT", "%REQ", "%LIMIT")
	table.SetHeader(tblHead...)

	if len(commonFlagList.filterList) >= 1 {
		err = table.SetFilter(commonFlagList.filterList)
		if err != nil {
			return err
		}
	}

	commonFlagList.showPodName = showPodName
	columnInfo.SetVisibleColumns(table, commonFlagList)

	podState := podMetrics2Hashtable(podStateList)
	for _, pod := range podList {
		columnInfo.LoadFromPod(pod)

		if columnInfo.labelNodeName != "" {
			columnInfo.labelNodeValue = nodeLabels[pod.Spec.NodeName][columnInfo.labelNodeName]
		}
		if columnInfo.labelPodName != "" {
			columnInfo.labelPodValue = podLabels[pod.Name][columnInfo.labelPodName]
		}

		//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
		if columnInfo.treeView {
			tblOut := podStatsProcessBuildRow(pod, columnInfo)
			columnInfo.ApplyRow(&table, tblOut)
		}

		if commonFlagList.showInitContainers {
			// process init containers
			columnInfo.containerType = "I"
			for _, container := range pod.Spec.InitContainers {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				columnInfo.containerName = container.Name
				tblOut := statsProcessTableRow(container.Resources, podState[pod.Name][container.Name], columnInfo, resourceType, showRaw, commonFlagList.byteSize)
				columnInfo.ApplyRow(&table, tblOut)
			}
		} else {
			// hide the container type column as its only needed when the init containers are being shown
			if !columnInfo.treeView {
				table.HideColumn(0)
			}
		}

		// process standard containers
		columnInfo.containerType = "S"
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := statsProcessTableRow(container.Resources, podState[pod.Name][container.Name], columnInfo, resourceType, showRaw, commonFlagList.byteSize)
			columnInfo.ApplyRow(&table, tblOut)
		}

		columnInfo.containerType = "E"
		for _, container := range pod.Spec.EphemeralContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := statsProcessTableRow(container.Resources, podState[pod.Name][container.Name], columnInfo, resourceType, showRaw, commonFlagList.byteSize)
			columnInfo.ApplyRow(&table, tblOut)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	// at this point we have data on all containers

	// do we need to find the outliers, we have enough data to compute a range
	if commonFlagList.showOddities {
		row2Remove, err := table.ListOutOfRange(defaultHeaderLen+0, table.GetRows()) //1 = used column
		if err != nil {
			return err
		}
		table.HideRows(row2Remove)
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil
}

func podStatsProcessBuildRow(pod v1.Pod, info containerInfomation) []Cell {

	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.podName)), //name
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}
}

func statsProcessTableRow(res v1.ResourceRequirements, metrics v1.ResourceList, info containerInfomation, resource string, showRaw bool, bytesAs string) []Cell {
	var cellList []Cell
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
				displayValue = fmt.Sprintf("%dm", metrics.Cpu().MilliValue())
				floatfmt = "%.2f"
			}

			// limit = res.Limits.Cpu().Value()
			rawLimit = res.Limits.Cpu().MilliValue()
			limit = fmt.Sprintf("%dm", rawLimit)
			// request = res.Requests.Cpu().String()
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
			// limit = fmt.Sprintf("%d", rawLimit)
			request = res.Requests.Memory().String()
			rawRequest = res.Requests.Memory().Value()
			// request = fmt.Sprintf("%d", rawRequest)

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

	if info.treeView {
		cellList = buildTreeCell(info, cellList)
	}

	cellList = append(cellList,
		NewCellInt(displayValue, rawValue),
		NewCellInt(request, rawRequest),
		NewCellInt(limit, rawLimit),
		NewCellFloat(percentRequest, rawPercentRequest),
		NewCellFloat(percentLimit, rawPercentLimit),
	)

	return cellList
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
