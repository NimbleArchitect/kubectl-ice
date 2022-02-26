package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Restarts(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname string
	var showPodName bool = true
	var idx int

	clientset, err := loadConfig(kubeFlags, cmd)
	if err != nil {
		return err
	}

	// fmt.Println(args)
	//TODO: allow multipule pods to be specified on cmdline
	if len(args) >= 1 {
		podname = args[0]
		if len(podname) >= 1 {
			showPodName = false
		}
	}

	podList, err := getPods(clientset, kubeFlags, podname)
	if err != nil {
		return err
	}

	table := make(map[int][]string)
	table[0] = []string{"T", "NAME", "RESTARTS"}

	if showPodName == true {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		for _, container := range pod.Status.ContainerStatuses {
			idx++
			table[idx] = restartsBuildRow(container, "S")
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
		for _, container := range pod.Status.InitContainerStatuses {
			idx++
			table[idx] = restartsBuildRow(container, "I")
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
	}
	showTable(table)
	return nil

}

func restartsBuildRow(container v1.ContainerStatus, containerType string) []string {
	// if container.RestartCount == 0
	// restarts := fmt.Sprintf("%d", container.RestartCount)

	return []string{
		containerType,
		container.Name,
		fmt.Sprintf("%d", container.RestartCount),
	}
}
