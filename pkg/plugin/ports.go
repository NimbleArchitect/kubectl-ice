package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var portsShort = "Shows ports exposed by the containers in a pod"

var portsDescription = ` Print a details of service ports exposed by containers in a pod. Details include the container 
name, port number and protocol type. Port name and host port are only show if avaliable. If no
name is specified the container port details of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`

var portsExample = `  # List containers port info from pods
  %[1]s ports

  # List container port info from pods output in JSON format
  %[1]s ports -o json

  # List container port info from a single pod
  %[1]s ports my-pod-4jh36

  # List port info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s ports -c web-container

  # List port info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s ports -c web-container --sort '!CONTAINER'

  # List port info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s ports -c web-container --sort PODNAME

  # List container port info from all pods where label app matches web
  %[1]s ports -l app=web

  # List container port info from all pods where the pod label app is either web or mail
  %[1]s ports -l "app in (web,mail)"`

func Ports(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
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
		"T", "PODNAME", "CONTAINER", "PORTNAME", "PORT", "PROTO", "HOSTPORT",
	)
	table.SetColumnTypeInt(4, 6)

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
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				info.containerName = container.Name
				tblOut := portsBuildRow(info, port)
				table.AddRow(tblOut...)
			}
		}

		info.containerType = "I"
		for _, container := range pod.Spec.InitContainers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				info.containerName = container.Name
				tblOut := portsBuildRow(info, port)
				table.AddRow(tblOut...)
			}
		}

		info.containerType = "E"
		for _, container := range pod.Spec.EphemeralContainers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				info.containerName = container.Name
				tblOut := portsBuildRow(info, port)
				table.AddRow(tblOut...)
			}
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func portsBuildRow(info containerInfomation, port v1.ContainerPort) []Cell {
	hostPort := ""

	if port.HostPort > 0 {
		hostPort = fmt.Sprintf("%d", port.HostPort)
	}
	return []Cell{
		NewCellText(info.containerType),
		NewCellText(info.podName),
		NewCellText(info.containerName),
		NewCellText(port.Name),
		NewCellText(fmt.Sprintf("%d", port.ContainerPort)),
		NewCellText(string(port.Protocol)),
		NewCellText(hostPort),
	}
}
