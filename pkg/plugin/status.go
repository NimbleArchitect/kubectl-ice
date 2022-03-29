package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var statusShort = "List status of each container in a pod"

var statusDescription = ` Prints container status information from pods, current and previous exit code, reason and signal
are shown slong with current ready and running state. Pods and containers can also be selected
by name. If no name is specified the container state of all pods in the current namespace is
shown.

The T column in the table output denotes S for Standard and I for init containers`

var statusExample = `  # List individual container status from pods
  %[1]s status

  # List conttainers status from pods output in JSON format
  %[1]s status -o json

  # List status from all container in a single pod
  %[1]s status my-pod-4jh36

  # List previous container status from a single pod
  %[1]s status -p my-pod-4jh36

  # List status of all containers named web-container searching all 
  # pods in the current namespace
  %[1]s status -c web-container

  # List status of containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s status -c web-container --sort '!CONTAINER'

  # List status of containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s status -c web-container --sort 'PODNAME"

  # List container status from all pods where label app equals web
  %[1]s status -l app=web

  # List status from all containers where the pods label app is either web or mail
  %[1]s status -l "app in (web,mail)"`

func Status(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true
	var showPrevious bool

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

	if cmd.Flag("previous").Value.String() == "true" {
		showPrevious = true
	}

	table := Table{}
	if !showPrevious {
		table.SetHeader(
			"T", "PODNAME", "CONTAINER", "READY", "STARTED", "RESTARTS", "STATE", "REASON", "EXIT-CODE", "SIGNAL", "TIMESTAMP", "MESSAGE",
		)
		table.SetColumnTypeInt(5, 8, 9)
	} else {
		table.SetHeader(
			"T", "PODNAME", "CONTAINER", "STATE", "REASON", "EXIT-CODE", "SIGNAL", "TIMESTAMP", "MESSAGE",
		)
		table.SetColumnTypeInt(5, 6)
	}

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
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := statusBuildRow(container, pod.Name, "S", showPrevious)
			table.AddRow(tblOut...)
		}
		for _, container := range pod.Status.InitContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := statusBuildRow(container, pod.Name, "I", showPrevious)
			table.AddRow(tblOut...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func statusBuildRow(container v1.ContainerStatus, podName string, containerType string, showPrevious bool) []Cell {
	var reason string
	var exitCode string
	var signal string
	var message string
	var startedAt string
	var started string
	var strState string
	var state v1.ContainerState
	var rawExitCode, rawSignal, rawRestarts int64

	// fmt.Println("F:statusBuildRow:Name=", container.Name)

	if showPrevious {
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
		rawExitCode = int64(state.Terminated.ExitCode)
		signal = fmt.Sprintf("%d", state.Terminated.Signal)
		rawSignal = int64(state.Terminated.Signal)
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
	rawRestarts = int64(container.RestartCount)

	if showPrevious {
		return []Cell{
			NewCellText(containerType),
			NewCellText(podName),
			NewCellText(container.Name),
			NewCellText(strState),
			NewCellText(reason),
			NewCellInt(exitCode, rawExitCode),
			NewCellInt(signal, rawSignal),
			NewCellText(startedAt),
			NewCellText(message),
		}
	} else {
		return []Cell{
			NewCellText(containerType),
			NewCellText(podName),
			NewCellText(container.Name),
			NewCellText(ready),
			NewCellText(started),
			NewCellInt(restarts, rawRestarts),
			NewCellText(strState),
			NewCellText(reason),
			NewCellInt(exitCode, rawExitCode),
			NewCellInt(signal, rawSignal),
			NewCellText(startedAt),
			NewCellText(message),
		}
	}

}
