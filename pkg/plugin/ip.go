package plugin

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var ipShort = "List ip addresses of all pods in the namespace listed"

var ipDescription = ` Prints the known IP addresses of the specified pod(s). if no pod is specified the IP address of
all pods in the current namespace are shown.`

var ipExample = `  # List IP address of pods
  %[1]s ip

  # List IP address of pods output in JSON format
  %[1]s ip -o json

  # List IP address a single pod
  %[1]s ip my-pod-4jh36

  # List IP address of all pods where label app matches web
  %[1]s ip -l app=web

  # List IP address of all pods where the pod label app is either web or mail
  %[1]s ip -l "app in (web,mail)"`

func IP(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
	}
	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList

	podList, err := connect.GetPods(podname)
	if err != nil {
		return err
	}
	table := Table{}
	table.SetHeader(
		"NAME", "IP",
	)

	for _, pod := range podList {

		table.AddRow(
			NewCellText(pod.Name),
			NewCellText(pod.Status.PodIP),
		)
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}
