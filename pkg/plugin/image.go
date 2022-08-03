package plugin

import (
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
	log := logger{location: "Image"}
	log.Debug("Start")

	loopinfo := image{}
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
	builder.SetFlagsFrom(commonFlagList)

	table := Table{}
	builder.Table = &table
	builder.CommonFlags = commonFlagList
	builder.Connection = &connect

	builder.Build(&loopinfo)

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type image struct {
}

func (s *image) Headers() []string {
	return []string{
		"PULL", "IMAGE",
	}
}

func (s *image) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *image) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *image) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s *image) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	out := []Cell{
		NewCellText(""),
		NewCellText(""),
	}
	return out, nil
}

func (s *image) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	out[0] = s.imageBuildRow(info, container.Image, string(container.ImagePullPolicy))
	return out, nil
}

func (s *image) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	out[0] = s.imageBuildRow(info, container.Image, string(container.ImagePullPolicy))
	return out, nil
}

func (s *image) imageBuildRow(info BuilderInformation, imageName string, pullPolicy string) []Cell {
	var cellList []Cell

	cellList = append(cellList,
		NewCellText(pullPolicy),
		NewCellText(imageName),
	)

	return cellList
}
