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
	// var tblHead []string
	var podname []string
	// var showPodName bool = true
	// var nodeLabels map[string]map[string]string
	// var podLabels map[string]map[string]string

	log := logger{location: "Capabilities"}
	log.Debug("Start")

	loopinfo := capabilities{}
	builder := RowBuilder{}
	builder.LoopSpec = true
	builder.ShowPodName = true
	builder.ShowInitContainers = true

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			log.Debug("builder.ShowPodName = false")
			builder.ShowPodName = false
		}
	}
	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList
	builder.CommonFlags = commonFlagList
	builder.Connection = &connect

	// podList, err := connect.GetPods(podname)
	// if err != nil {
	// 	return err
	// }

	if cmd.Flag("node-label").Value.String() != "" {
		label := cmd.Flag("node-label").Value.String()
		log.Debug("builder.LabelNodeName =", label)
		builder.LabelNodeName = label
	}

	if cmd.Flag("pod-label").Value.String() != "" {
		label := cmd.Flag("pod-label").Value.String()
		log.Debug("builder.LabelPodName =", label)
		builder.LabelPodName = label
	}

	table := Table{}
	builder.Table = &table
	columnInfo.table = &table
	builder.ShowTreeView = commonFlagList.showTreeView
	// columnInfo.treeView = commonFlagList.showTreeView

	// tblHead = columnInfo.GetDefaultHead()
	// if commonFlagList.showTreeView {
	// 	// we have to control the name when displaying a tree view as the table
	// 	//  object dosent have the extra info to be able to process it
	// 	tblHead = append(tblHead, "NAME")
	// }

	// tblHead = append(tblHead, "ADD", "DROP")
	// table.SetHeader(tblHead...)

	// if len(commonFlagList.filterList) >= 1 {
	// 	err = table.SetFilter(commonFlagList.filterList)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// commonFlagList.showPodName = showPodName
	// columnInfo.SetVisibleColumns(table, commonFlagList)

	builder.BuildRows(loopinfo)

	// for _, pod := range podList {
	// 	columnInfo.LoadFromPod(pod)

	// 	if columnInfo.labelNodeName != "" {
	// 		columnInfo.labelNodeValue = nodeLabels[pod.Spec.NodeName][columnInfo.labelNodeName]
	// 	}
	// 	if columnInfo.labelPodName != "" {
	// 		columnInfo.labelPodValue = podLabels[pod.Name][columnInfo.labelPodName]
	// 	}

	// 	//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
	// 	if columnInfo.treeView {
	// 		tblOut := podCapabilitiesBuildRow(pod, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 	}

	// 	columnInfo.containerType = "S"
	// 	for _, container := range pod.Spec.Containers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		tblOut := capabilitiesBuildRow(container.SecurityContext, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}

	// 	columnInfo.containerType = "I"
	// 	for _, container := range pod.Spec.InitContainers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		tblOut := capabilitiesBuildRow(container.SecurityContext, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}

	// 	columnInfo.containerType = "E"
	// 	for _, container := range pod.Spec.EphemeralContainers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		tblOut := capabilitiesBuildRow(container.SecurityContext, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}
	// }

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type capabilities struct {
}

func (s capabilities) Headers() []string {
	return []string{
		"ADD", "DROP",
	}
}

func (s capabilities) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s capabilities) HideColumns(info BuilderInformation) []int {
	return []int{}
}

// func podStatsProcessBuildRow(pod v1.Pod, info containerInfomation) []Cell {
func (s capabilities) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.PodName)), //name
		NewCellText(""),
		NewCellText(""),
	}, nil
}

func (s capabilities) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	out[0] = s.capabilitiesBuildRow(container.SecurityContext, info)
	return out, nil
}

func (s capabilities) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	out[0] = s.capabilitiesBuildRow(container.SecurityContext, info)
	return out, nil
}

func (s capabilities) capabilitiesBuildRow(securityContext *v1.SecurityContext, info BuilderInformation) []Cell {
	// func capabilitiesBuildRow(securityContext *v1.SecurityContext, info containerInfomation) []Cell {
	var cellList []Cell

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

	// if info.treeView {
	// 	cellList = buildTreeCell(info, cellList)
	// }

	if info.TreeView {
		cellList = info.BuildTreeCell(cellList)
	}

	cellList = append(cellList,
		NewCellText(capAdd),
		NewCellText(capDrop),
	)

	return cellList
}
