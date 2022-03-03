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
	var idx int
	var allNamespaces bool

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

	if cmd.Flag("all-namespaces").Value.String() == "true" {
		allNamespaces = true
	}

	podList, err := getPods(clientset, kubeFlags, podname, allNamespaces)
	if err != nil {
		return err
	}

	table := make(map[int][]string)
	table[0] = []string{"T", "CONTAINER", "COMMAND", "ARGUMENTS"}

	if showPodName {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {
			idx++
			table[idx] = commandsBuildRow(container, "S")
			if showPodName {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
		for _, container := range pod.Spec.InitContainers {
			idx++
			table[idx] = commandsBuildRow(container, "I")
			if showPodName {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
	}
	showTable(table)
	return nil

}

func commandsBuildRow(container v1.Container, containerType string) []string {

	return []string{
		containerType,
		container.Name,
		strings.Join(container.Command, " "),
		strings.Join(container.Args, " "),
	}
}
