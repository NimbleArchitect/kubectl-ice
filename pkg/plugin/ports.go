package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Ports(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool
	var idx int
	var allNamespaces bool

	clientset, err := loadConfig(kubeFlags)
	if err != nil {
		return err
	}

	// fmt.Println(args)
	//TODO: allow multipule pods to be specified on cmdline
	showPodName = false
	if len(args) >= 1 {
		podname = args
		if len(podname) >= 1 {
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
	table[0] = []string{"T", "CONTAINER", "PORTNAME", "PORT", "PROTO", "HOSTPORT"}

	if showPodName {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				idx++
				table[idx] = portsBuildRow(container, port, "S")
				if showPodName {
					table[idx] = append([]string{pod.Name}, table[idx]...)
				}
			}
		}
		for _, container := range pod.Spec.InitContainers {
			for _, port := range container.Ports {
				idx++
				table[idx] = portsBuildRow(container, port, "I")
				if showPodName {
					table[idx] = append([]string{pod.Name}, table[idx]...)
				}
			}
		}
	}
	showTable(table)
	return nil

}

func portsBuildRow(container v1.Container, port v1.ContainerPort, containerType string) []string {

	return []string{
		containerType,
		container.Name,
		port.Name,
		fmt.Sprintf("%d", port.ContainerPort),
		string(port.Protocol),
	}
}
