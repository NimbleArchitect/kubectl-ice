package plugin

import (
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var commandsShort = "Retrieves the command line and any arguments specified at the container level"

var commandsDescription = ` Prints command and arguments used to start each container (if specified), single pods and 
containers can be selected by name.  If no name is specified the container commands of all pods
in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`

var commandsExample = `  # List containers command info from pods
  %[1]s command

  # List container command info from pods output in JSON format
  %[1]s command -o json

  # List container command info from a single pod
  %[1]s command my-pod-4jh36

  # List command info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s command -c web-container

  # List command info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s command -c web-container --sort '!CONTAINER'

  # List command info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s command -c web-container --sort PODNAME

  # List container command info from all pods where label app matches web
  %[1]s command -l app=web

  # List container command info from all pods where the pod label app is either web or mail
  %[1]s command -l "app in (web,mail)"`

type commandLine struct {
	cmd  []string
	args []string
}

func Commands(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {

	log := logger{location: "Commands"}
	log.Debug("Start")

	loopinfo := commands{}
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

	builder.Build(loopinfo)

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type commands struct {
}

func (s commands) Headers() []string {
	return []string{
		"COMMAND", "ARGUMENTS",
	}
}

func (s commands) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s commands) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s commands) BuildBranch(info BuilderInformation, podList []v1.Pod) ([][]Cell, error) {
	out := []Cell{
		NewCellText(""),
		NewCellText("")}
	return [][]Cell{out}, nil
}

// func podStatsProcessBuildRow(pod v1.Pod, info containerInfomation) []Cell {
func (s commands) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(""),
		NewCellText(""),
	}, nil
}

func (s commands) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	cmdLine := commandLine{
		cmd:  container.Command,
		args: container.Args,
	}
	out := make([][]Cell, 1)
	out[0] = s.commandsBuildRow(cmdLine, info)
	return out, nil
}

func (s commands) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	cmdLine := commandLine{
		cmd:  container.Command,
		args: container.Args,
	}
	out := make([][]Cell, 1)
	out[0] = s.commandsBuildRow(cmdLine, info)
	return out, nil
}

func (s commands) commandsBuildRow(cmdLine commandLine, info BuilderInformation) []Cell {
	var cellList []Cell

	// if info.TreeView {
	// 	cellList = info.BuildTreeCell(cellList)
	// }

	cellList = append(cellList,
		NewCellText(strings.Join(cmdLine.cmd, " ")),
		NewCellText(strings.Join(cmdLine.args, " ")),
	)

	return cellList
}
