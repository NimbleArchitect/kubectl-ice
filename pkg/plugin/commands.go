package plugin

import (
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

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

	commonFlagList := processCommonFlags(cmd)

	podList, err := getPods(clientset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	table := Table{}
	table.SetHeader(
		"T", "PODNAME", "CONTAINER", "COMMAND", "ARGUMENTS",
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
	table.Print()
	return nil

}

func commandsBuildRow(container v1.Container, podName string, containerType string) []string {

	return []string{
		containerType,
		podName,
		container.Name,
		strings.Join(container.Command, " "),
		strings.Join(container.Args, " "),
	}
}
