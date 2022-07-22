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

	log := logger{location: "Ports"}
	log.Debug("Start")

	loopinfo := ports{}
	builder := RowBuilder{}
	builder.LoopSpec = true
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
	builder.ShowTreeView = commonFlagList.showTreeView

	builder.Build(&loopinfo)

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type ports struct {
}

func (s *ports) Headers() []string {
	return []string{
		"PORTNAME", "PORT", "PROTO", "HOSTPORT",
	}
}

func (s *ports) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *ports) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *ports) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s *ports) BuildBranch(info BuilderInformation, podList []v1.Pod) ([]Cell, error) {
	out := []Cell{
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}
	return out, nil
}

func (s *ports) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}, nil
}

func (s *ports) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	for _, port := range container.Ports {
		out = append(out, s.portsBuildRow(info, port))
	}
	return out, nil
}

func (s *ports) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	for _, port := range container.Ports {
		out = append(out, s.portsBuildRow(info, port))
	}
	return out, nil
}

func (s *ports) Sum(rows [][]Cell) []Cell {
	rowOut := make([]Cell, 4)
	return rowOut
}

func (s *ports) portsBuildRow(info BuilderInformation, port v1.ContainerPort) []Cell {
	var cellList []Cell

	hostPort := Cell{}

	if port.HostPort > 0 {
		hostPort = NewCellInt(fmt.Sprintf("%d", port.HostPort), int64(port.HostPort))
	} else {
		hostPort = NewCellText("")
	}

	// if info.TreeView {
	// 	cellList = info.BuildTreeCell(cellList)
	// }

	cellList = append(cellList,
		NewCellText(port.Name),
		NewCellInt(fmt.Sprintf("%d", port.ContainerPort), int64(port.ContainerPort)),
		NewCellText(string(port.Protocol)),
		hostPort,
	)
	return cellList
}
