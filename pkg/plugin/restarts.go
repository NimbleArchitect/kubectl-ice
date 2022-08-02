package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
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

	log := logger{location: "Restarts"}
	log.Debug("Start")

	loopinfo := restarts{}
	builder := RowBuilder{}
	builder.LoopStatus = true
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

	// do we need to find the outliers, we have enough data to compute a range
	if commonFlagList.showOddities {
		row2Remove, err := table.ListOutOfRange(4) //3 = restarts column
		if err != nil {
			return err
		}
		table.HideRows(row2Remove)
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type restarts struct{}

func (s restarts) Headers() []string {
	return []string{
		"RESTARTS",
	}
}

func (s restarts) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	out[0] = s.restartsBuildRow(info, container.RestartCount)
	return out, nil
}

func (s restarts) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	out[0] = s.restartsBuildRow(info, container.RestartCount)
	return out, nil
}

func (s restarts) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s restarts) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	rowOut := make([]Cell, 1)

	switch info.TypeName {
	case "Pod":
		for _, r := range rows {
			rowOut[0].number += r[0].number //ready
		}
		rowOut[0].text = fmt.Sprintf("%d", rowOut[0].number)
	}

	return rowOut, nil
}

func (s restarts) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	return out, nil
}

func (s restarts) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	return out, nil
}

func (s restarts) restartsBuildRow(info BuilderInformation, restartCount int32) []Cell {
	var cellList []Cell

	cellList = append(cellList,
		NewCellInt(fmt.Sprintf("%d", restartCount), int64(restartCount)),
	)

	return cellList
}
