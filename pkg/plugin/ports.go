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
	var columnInfo containerInfomation
	var tblHead []string
	var podname []string
	var showPodName bool = true
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
	if commonFlagList.showTreeView {
		// we have to control the name when displaying a tree view as the table
		//  object dosent have the extra info to be able to process it
		tblHead = append(tblHead, "NAME")
	}

	tblHead = append(tblHead, "PORTNAME", "PORT", "PROTO", "HOSTPORT")
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

		if columnInfo.labelNodeName != "" {
			columnInfo.labelNodeValue = nodeLabels[pod.Spec.NodeName][columnInfo.labelNodeName]
		}
		if columnInfo.labelPodName != "" {
			columnInfo.labelPodValue = podLabels[pod.Name][columnInfo.labelPodName]
		}

		//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
		if columnInfo.treeView {
			tblOut := podPortsBuildRow(pod, columnInfo)
			columnInfo.ApplyRow(&table, tblOut)
		}

		columnInfo.containerType = "S"
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				columnInfo.containerName = container.Name
				tblOut := portsBuildRow(columnInfo, port)
				columnInfo.ApplyRow(&table, tblOut)
				// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
				// table.AddRow(tblFullRow...)
			}
		}

		columnInfo.containerType = "I"
		for _, container := range pod.Spec.InitContainers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				columnInfo.containerName = container.Name
				tblOut := portsBuildRow(columnInfo, port)
				columnInfo.ApplyRow(&table, tblOut)
				// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
				// table.AddRow(tblFullRow...)
			}
		}

		columnInfo.containerType = "E"
		for _, container := range pod.Spec.EphemeralContainers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				columnInfo.containerName = container.Name
				tblOut := portsBuildRow(columnInfo, port)
				columnInfo.ApplyRow(&table, tblOut)
				// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
				// table.AddRow(tblFullRow...)
			}
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func podPortsBuildRow(pod v1.Pod, info containerInfomation) []Cell {

	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.podName)), //name
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}
}

func portsBuildRow(info containerInfomation, port v1.ContainerPort) []Cell {
	var cellList []Cell

	hostPort := Cell{}

	if port.HostPort > 0 {
		hostPort = NewCellInt(fmt.Sprintf("%d", port.HostPort), int64(port.HostPort))
	} else {
		hostPort = NewCellText("")
	}

	if info.treeView {
		cellList = buildTreeCell(info, cellList)
	}

	cellList = append(cellList,
		NewCellText(port.Name),
		NewCellInt(fmt.Sprintf("%d", port.ContainerPort), int64(port.ContainerPort)),
		NewCellText(string(port.Protocol)),
		hostPort,
	)
	return cellList
}
