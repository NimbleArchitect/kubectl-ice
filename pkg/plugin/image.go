package plugin

import (
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var imageShort = "list the image name and pull status for each container"

var imageDescription = `.`

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
  %[1]s image -c web-container --sort 'PODNAME"

  # List container image info from all pods where label app matches web
  %[1]s image -l app=web

  # List container image info from all pods where the pod label app is either web or mail
  %[1]s image -l "app in (web,mail)"`

func Image(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
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
		"T", "PODNAME", "CONTAINER", "PULL", "IMAGE",
	)

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
			tblOut := imageBuildRow(container, pod.Name, "S")
			table.AddRow(tblOut...)
		}
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := imageBuildRow(container, pod.Name, "I")
			table.AddRow(tblOut...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func imageBuildRow(container v1.Container, podName string, containerType string) []string {
	return []string{
		containerType,
		podName,
		container.Name,
		string(container.ImagePullPolicy),
		container.Image,
	}
}
