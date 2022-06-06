package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var restartsShort = "Show restart counts for each container in a named pod"

var restartsDescription = ` Prints container name and restart count for individual containers. If no name is specified the
container restart counts of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard, I for init and E for Ephemerial containers`

var restartsExample = `  # List individual container restart count from pods
  %[1]s restarts

  # List conttainers restart count from pods output in JSON format
  %[1]s restarts -o json

  # List restart count from all containers in a single pod
  %[1]s restarts my-pod-4jh36

  # List restart count of all containers named web-container searching all 
  # pods in the current namespace
  %[1]s restarts -c web-container

  # List restart count of containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s restarts -c web-container --sort '!CONTAINER'

  # List restart count of containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s restarts -c web-container --sort PODNAME

  # List container restart count from all pods where label app equals web
  %[1]s restarts -l app=web

  # List restart count from all containers where the pod label app is either web or mail
  %[1]s restarts -l "app in (web,mail)"`

func Restarts(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true

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

	table := Table{}
	table.SetHeader(
		"T", "PODNAME", "CONTAINER", "RESTARTS",
	)
	table.SetColumnTypeInt(3)

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

	for _, pod := range podList {
		info := containerInfomation{
			podName: pod.Name,
		}

		info.containerType = "S"
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			info.containerName = container.Name
			tblOut := restartsBuildRow(info, container.RestartCount)
			table.AddRow(tblOut...)
		}

		info.containerType = "I"
		for _, container := range pod.Status.InitContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			info.containerName = container.Name
			tblOut := restartsBuildRow(info, container.RestartCount)
			table.AddRow(tblOut...)
		}

		info.containerType = "E"
		for _, container := range pod.Status.EphemeralContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			info.containerName = container.Name
			tblOut := restartsBuildRow(info, container.RestartCount)
			table.AddRow(tblOut...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	// do we need to find the outliers, we have enough data to compute a range
	if commonFlagList.showOddities {
		row2Remove, err := table.ListOutOfRange(3, table.GetRows()) //3 = restarts column
		if err != nil {
			return err
		}
		table.HideRows(row2Remove)
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func restartsBuildRow(info containerInfomation, restartCount int32) []Cell {
	// if container.RestartCount == 0
	// restarts := fmt.Sprintf("%d", container.RestartCount)

	return []Cell{
		NewCellText(info.containerType),
		NewCellText(info.podName),
		NewCellText(info.containerName),
		NewCellInt(fmt.Sprintf("%d", restartCount), int64(restartCount)),
	}
}
