package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Status(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname string
	var showPodName bool = true
	var showPrevious bool
	var idx int

	// onfigTest(cmd, kubeFlags, args)
	// if true {
	// 	return nil
	// }

	// kubeFlags.AddFlags(flagList)
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

	if cmd.Flag("previous").Value.String() == "true" {
		showPrevious = true
	}

	table := make(map[int][]string)
	if showPrevious == false {
		table[0] = []string{"T", "NAME", "READY", "STARTED", "RESTARTS", "STATE", "REASON", "EXIT-CODE", "SIGNAL", "TIMESTAMP", "MESSAGE"}
	} else {
		table[0] = []string{"T", "NAME", "STATE", "REASON", "EXIT-CODE", "SIGNAL", "TIMESTAMP", "MESSAGE"}
	}

	if showPodName == true {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		for _, container := range pod.Status.ContainerStatuses {
			idx++
			table[idx] = statusBuildRow(container, "S", showPrevious)
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
		for _, container := range pod.Status.InitContainerStatuses {
			idx++
			table[idx] = statusBuildRow(container, "I", showPrevious)
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
	}
	showTable(table)
	return nil

}

func statusBuildRow(container v1.ContainerStatus, containerType string, showPrevious bool) []string {
	var reason string
	var exitCode string
	var signal string
	var message string
	var startedAt string
	var started string
	var strState string
	var state v1.ContainerState

	// fmt.Println("F:statusBuildRow:Name=", container.Name)

	if showPrevious == true {
		state = container.LastTerminationState
	} else {
		state = container.State
	}

	if state.Waiting != nil {
		strState = "Waiting"
		reason = state.Waiting.Reason
		message = state.Waiting.Message
	}

	if state.Terminated != nil {
		strState = "Terminated"
		exitCode = fmt.Sprintf("%d", state.Terminated.ExitCode)
		signal = fmt.Sprintf("%d", state.Terminated.Signal)
		startedAt = state.Terminated.StartedAt.String()
		reason = state.Terminated.Reason
		message = state.Terminated.Message
	}

	if state.Running != nil {
		strState = "Running"
		startedAt = state.Running.StartedAt.String()
	}

	if container.Started != nil {
		started = fmt.Sprintf("%t", *container.Started)
	}
	ready := fmt.Sprintf("%t", container.Ready)
	restarts := fmt.Sprintf("%d", container.RestartCount)

	if showPrevious == true {
		return []string{
			containerType,
			container.Name,
			strState,
			reason,
			exitCode,
			signal,
			startedAt,
			message,
		}
	} else {
		return []string{
			containerType,
			container.Name,
			ready,
			started,
			restarts,
			strState,
			reason,
			exitCode,
			signal,
			startedAt,
			message,
		}
	}

}