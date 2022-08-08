package plugin

import (
	"strings"

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

	if cmd.Flag("id").Value.String() == "true" {
		log.Debug("loopinfo.ShowID = true")
		loopinfo.ShowID = true
	}

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
	ShowID bool
}

func (s *image) Headers() []string {
	return []string{
		"PULL", "IMAGEID", "CONTAINERID", "IMAGE", "TAG",
	}
}

func (s *image) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *image) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *image) HideColumns(info BuilderInformation) []int {
	var hideColumns []int

	if !s.ShowID {
		hideColumns = append(hideColumns, 1, 2)
	}

	return hideColumns
}

func (s *image) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	out := make([]Cell, len(s.Headers()))

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
	var imageID string
	var containerID string
	var cellList []Cell

	name := imageName
	tag := ""

	if strings.Contains(imageName, "/") {
		arrPath := strings.Split(imageName, "/")
		if c := len(arrPath); c > 0 {
			tmp := strings.Split(arrPath[c-1], ":")
			if len(tmp) > 0 {
				tag = strings.Join(tmp[1:], ":")
				//calculate the uri length
				namelen := len(imageName) - len(tag)
				if len(tag) > 0 {
					// check a tag was supplied so we dont cut off the last char of the image name
					namelen--
				}
				name = imageName[0:namelen]
			}
		}
	} else {
		arrImage := strings.Split(imageName, ":")
		if c := len(arrImage); c > 0 {
			tag = arrImage[c-1]
			name = strings.Join(arrImage[:c-1], ":")
		}
	}

	for _, status := range info.Data.pod.Status.InitContainerStatuses {
		if status.Image == imageName {
			imageID = status.ImageID
			containerID = status.ContainerID
		}
	}
	for _, status := range info.Data.pod.Status.ContainerStatuses {
		if status.Image == imageName {
			imageID = status.ImageID
			containerID = status.ContainerID
		}
	}
	for _, status := range info.Data.pod.Status.EphemeralContainerStatuses {
		if status.Image == imageName {
			imageID = status.ImageID
			containerID = status.ContainerID
		}
	}

	if val := strings.Split(imageID, "@"); len(val) == 2 {
		imageID = val[1]
	}

	cellList = append(cellList,
		NewCellText(pullPolicy),
		NewCellText(imageID),
		NewCellText(containerID),
		NewCellText(name),
		NewCellText(tag),
	)

	return cellList
}
