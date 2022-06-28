package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var imageShort = "List the image name and pull status for each container"

var imageDescription = ` Print the the image used for running containers in a pod including the pull policy, single pods
and containers can be selected by name. If no name is specified the image details of all pods in
the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`

var imageExample = `  # List containers image info from pods
  %[1]s image

  # List container image info from pods output in JSON format
  %[1]s image -o json

  # List container image info from a single pod
  %[1]s image my-pod-4jh36

  # List image info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s image -c web-container

  # List image info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s image -c web-container --sort '!CONTAINER'

  # List image info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s image -c web-container --sort PODNAME

  # List container image info from all pods where label app matches web
  %[1]s image -l app=web

  # List container image info from all pods where the pod label app is either web or mail
  %[1]s image -l "app in (web,mail)"`

func Image(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
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

	tblHead = append(tblHead, "PULL", "IMAGE")
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
			tblOut := podImageBuildRow(pod, columnInfo)
			columnInfo.ApplyRow(&table, tblOut)
		}

		columnInfo.containerType = "S"
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := imageBuildRow(columnInfo, container.Image, string(container.ImagePullPolicy))
			columnInfo.ApplyRow(&table, tblOut)
			// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
			// table.AddRow(tblFullRow...)
		}

		columnInfo.containerType = "I"
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := imageBuildRow(columnInfo, container.Image, string(container.ImagePullPolicy))
			columnInfo.ApplyRow(&table, tblOut)
			// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
			// table.AddRow(tblFullRow...)
		}

		columnInfo.containerType = "E"
		for _, container := range pod.Spec.EphemeralContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			columnInfo.containerName = container.Name
			tblOut := imageBuildRow(columnInfo, container.Image, string(container.ImagePullPolicy))
			columnInfo.ApplyRow(&table, tblOut)
			// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
			// table.AddRow(tblFullRow...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func podImageBuildRow(pod v1.Pod, info containerInfomation) []Cell {
	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.podName)), //name
		NewCellText(""),
		NewCellText(""),
	}
}

func imageBuildRow(info containerInfomation, imageName string, pullPolicy string) []Cell {
	var cellList []Cell
	if info.treeView {
		cellList = buildTreeCell(info, cellList)
	}

	cellList = append(cellList,
		NewCellText(pullPolicy),
		NewCellText(imageName),
	)

	return cellList
}
