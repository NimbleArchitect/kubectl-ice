package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func Stats(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname string
	var floatfmt string
	var showPodName bool = true
	var showRaw bool
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

	metricset, err := loadMetricConfig(kubeFlags)
	if err != nil {
		return err
	}
	podStateList, err := getMetricPods(metricset, kubeFlags, podname)
	if err != nil {
		return err
	}

	if cmd.Flag("raw").Value.String() == "true" {
		showRaw = true
	}

	table := make(map[int][]string)
	table[0] = []string{"NAME", "USED-CPU", "CPU-%-REQ", "CPU-%-LIMIT", "USED-MEM", "MEM-%-REQ", "MEM-%-LIMIT"}
	if showPodName == true {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	podState := podMetrics2Hashtable(podStateList)
	for _, pod := range podList {
		// fmt.Println("podname:", pod.Name)
		// table[pidx] = make([]string)
		containerList := append(pod.Spec.InitContainers, pod.Spec.Containers...)
		for _, container := range containerList {
			displayCpu := ""
			displayMemory := ""
			pCpuLimit := 0.0
			pCpuRequest := 0.0
			pMemoryLimit := 0.0
			pMemoryRequest := 0.0
			idx++
			metrics := podState[pod.Name][container.Name]
			if metrics.Cpu() != nil {
				pCpuLimit = validateFloat64(metrics.Cpu().AsApproximateFloat64()/container.Resources.Limits.Cpu().AsApproximateFloat64()) * 100
				pCpuRequest = validateFloat64(metrics.Cpu().AsApproximateFloat64()/container.Resources.Requests.Cpu().AsApproximateFloat64()) * 100
				pMemoryLimit = validateFloat64(metrics.Memory().AsApproximateFloat64()/container.Resources.Limits.Memory().AsApproximateFloat64()) * 100
				pMemoryRequest = validateFloat64(metrics.Memory().AsApproximateFloat64()/container.Resources.Requests.Memory().AsApproximateFloat64()) * 100
			}

			if showRaw == true {
				displayCpu = metrics.Cpu().String()
				displayMemory = fmt.Sprintf("%d", metrics.Memory().Value())
				floatfmt = "%f"
			} else {
				displayCpu = fmt.Sprintf("%d", metrics.Cpu().MilliValue())
				//TODO: make values human readable
				displayMemory = memoryHumanReadable(metrics.Memory().Value())
				// displayMemory = metrics.Memory().String()
				floatfmt = "%.2f"
			}

			table[idx] = []string{
				container.Name,
				displayCpu,
				fmt.Sprintf(floatfmt, pCpuRequest),
				fmt.Sprintf(floatfmt, pCpuLimit),
				displayMemory,
				fmt.Sprintf(floatfmt, pMemoryRequest),
				fmt.Sprintf(floatfmt, pMemoryLimit),
			}
			if showPodName == true {
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}

			// fmt.Println(pod.Name, container.Name, met.Cpu().Value())

		}
	}

	showTable(table)
	return nil
}

func podMetrics2Hashtable(stateList []v1beta1.PodMetrics) map[string]map[string]v1.ResourceList {
	podState := make(map[string]map[string]v1.ResourceList)

	for _, pod := range stateList {
		podState[pod.Name] = make(map[string]v1.ResourceList)
		for _, container := range pod.Containers {
			podState[pod.Name][container.Name] = container.Usage
		}
	}
	return podState
}
