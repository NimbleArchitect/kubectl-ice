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
	// var columnInfo containerInfomation
	// var podname []string

	log := logger{location: "Capabilities"}
	log.Debug("Start")

	loopinfo := capabilities{}
	builder := RowBuilder{}
	builder.LoopSpec = true
	// builder.ShowPodName = true
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
	builder.Connection = &connect
	builder.SetFlagsFrom(commonFlagList)

	table := Table{}
	builder.Table = &table
	// columnInfo.table = &table
	builder.ShowTreeView = commonFlagList.showTreeView

	builder.BuildRows(loopinfo)

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
