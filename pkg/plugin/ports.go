package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var portsExample = ``

func Ports(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
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
		"T", "PODNAME", "CONTAINER", "PORTNAME", "PORT", "PROTO", "HOSTPORT",
	)

	if !showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(1)
	}

	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				tblOut := portsBuildRow(container, pod.Name, port, "S")
				table.AddRow(tblOut...)
			}
		}
		for _, container := range pod.Spec.InitContainers {
			for _, port := range container.Ports {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				tblOut := portsBuildRow(container, pod.Name, port, "I")
				table.AddRow(tblOut...)
			}
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func portsBuildRow(container v1.Container, podName string, port v1.ContainerPort, containerType string) []string {
	hostPort := ""

	if port.HostPort > 0 {
		hostPort = fmt.Sprintf("%d", port.HostPort)
	}
	return []string{
		containerType,
		podName,
		container.Name,
		port.Name,
		fmt.Sprintf("%d", port.ContainerPort),
		string(port.Protocol),
		hostPort,
	}
}
