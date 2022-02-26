package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Resources(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string, resourceType string) error {
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
	table[0] = []string{"T", "NAME", "REQUEST", "LIMIT"}

	if showPodName == true {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {
			idx++
			table[idx] = resourcesBuildRow(container, "S", resourceType)
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
		for _, container := range pod.Spec.InitContainers {
			idx++
			table[idx] = resourcesBuildRow(container, "I", resourceType)
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
	}
	showTable(table)
	return nil

}

func resourcesBuildRow(container v1.Container, containerType string, resourceType string) []string {
	var request string
	var limit string

	if resourceType == "cpu" {
		limit = container.Resources.Limits.Cpu().String()
		request = container.Resources.Requests.Cpu().String()
	} else if resourceType == "memory" {
		limit = container.Resources.Limits.Memory().String()
		request = container.Resources.Requests.Memory().String()
	} else {
		fmt.Println("EROR: invalid resource")
	}

	return []string{
		containerType,
		container.Name,
		request,
		limit,
	}

}
