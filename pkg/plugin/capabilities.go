package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var capabilitiesShort = "Shows details of configured containers POSIX capabilities"

var capabilitiesDescription = ` View POSIX Capabilities that have been applied to the running containers.
`

var capabilitiesExample = `  # List container capabilities from pods
  %[1]s capabilities

  # List container capabilities info from pods output in JSON format
  %[1]s capabilities -o json

  # List container capabilities info from a single pod
  %[1]s capabilities my-pod-4jh36

  # List capabilities info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s capabilities -c web-container

  # List capabilities info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s capabilities -c web-container --sort '!CONTAINER'

  # List capabilities info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s capabilities -c web-container --sort PODNAME

  # List container capabilities info from all pods where label app matches web
  %[1]s capabilities -l app=web

  # List container capabilities info from all pods where the pod label app is either web or mail
  %[1]s capabilities -l "app in (web,mail)"`

//list details of configured liveness readiness and startup capabilities
func Capabilities(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var columnInfo containerInfomation
	var tblHead []string
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
	tblHead = append(columnInfo.GetDefaultHead(), "ADD", "DROP")
	table.SetHeader(tblHead...)

	if len(commonFlagList.filterList) >= 1 {
		err = table.SetFilter(commonFlagList.filterList)
		if err != nil {
			return err
		}
	}

	commonFlagList.showPodName = showPodName
	columnInfo.SetVisibleColumns(table, commonFlagList)

	for _, pod := range podList {
		columnInfo.LoadFromPod(pod)

		columnInfo.containerType = "S"
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := capabilitiesBuildRow(container.SecurityContext, columnInfo)
			tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
			table.AddRow(tblFullRow...)
		}

		columnInfo.containerType = "I"
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := capabilitiesBuildRow(container.SecurityContext, columnInfo)
			tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
			table.AddRow(tblFullRow...)
		}

		columnInfo.containerType = "E"
		for _, container := range pod.Spec.EphemeralContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := capabilitiesBuildRow(container.SecurityContext, columnInfo)
			tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
			table.AddRow(tblFullRow...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func capabilitiesBuildRow(securityContext *v1.SecurityContext, info containerInfomation) []Cell {
	capAdd := ""
	capDrop := ""

	if securityContext != nil {
		if securityContext.Capabilities != nil {
			for i, v := range securityContext.Capabilities.Add {
				sep := ","
				if i == 0 {
					sep = ""
				}
				capAdd += sep + fmt.Sprint(v)
			}

			for i, v := range securityContext.Capabilities.Drop {
				sep := ","
				if i == 0 {
					sep = ""
				}
				capDrop += sep + fmt.Sprint(v)
			}
		}
	}

	// capDrop := container.SecurityContext.Capabilities.Drop

	return []Cell{
		NewCellText(capAdd),
		NewCellText(capDrop),
	}
}
