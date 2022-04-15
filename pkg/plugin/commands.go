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

func Commands(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true

	clientset, err := loadConfig(kubeFlags)
	if err != nil {
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

	podList, err := getPods(clientset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	table := Table{}
	table.SetHeader(
		"T", "PODNAME", "CONTAINER", "COMMAND", "ARGUMENTS",
	)

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
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := commandsBuildRow(container, pod.Name, "S")
			table.AddRow(tblOut...)
		}
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := commandsBuildRow(container, pod.Name, "I")
			table.AddRow(tblOut...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func commandsBuildRow(container v1.Container, podName string, containerType string) []Cell {

	return []Cell{
		NewCellText(containerType),
		NewCellText(podName),
		NewCellText(container.Name),
		NewCellText(strings.Join(container.Command, " ")),
		NewCellText(strings.Join(container.Args, " ")),
	}
}
