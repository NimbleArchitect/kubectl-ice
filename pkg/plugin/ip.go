package plugin

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func IP(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname string
	var showPodName bool
	var idx int

	clientset, err := loadConfig(kubeFlags, cmd)
	if err != nil {
		return err
	}

	// fmt.Println(args)
	//TODO: allow multipule pods to be specified on cmdline
	showPodName = false
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
	table[0] = []string{"NAME", "IP"}

	if showPodName == true {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		idx++
		table[idx] = []string{
			pod.Name,
			pod.Status.PodIP,
		}
	}
	showTable(table)
	return nil

}