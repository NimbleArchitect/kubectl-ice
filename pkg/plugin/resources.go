package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func Resources(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string, resourceType string) error {
	var podname []string
	var showPodName bool = true
	var showRaw bool
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

	metricset, err := loadMetricConfig(kubeFlags)
	if err != nil {
		return err
	}
	podStateList, err := getMetricPods(metricset, kubeFlags, podname, allNamespaces)
	if err != nil {
		return err
	}

	if cmd.Flag("raw").Value.String() == "true" {
		showRaw = true
	}

	table := make(map[int][]string)
	table[0] = []string{"T", "CONTAINER", "USED", "REQUEST", "LIMIT", "%REQ", "%LIMIT"}

	if showPodName {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	podState := podMetrics2Hashtable(podStateList)
	for _, pod := range podList {
		// process init containers
		for _, container := range pod.Spec.InitContainers {
			idx++
			table[idx] = statsProcessTableRow(container, podState[pod.Name][container.Name], "I", resourceType, showRaw)
			if showPodName {
				// add podname to the beginning of the row
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}

		// process standard containers
		for _, container := range pod.Spec.Containers {
			idx++
			table[idx] = statsProcessTableRow(container, podState[pod.Name][container.Name], "S", resourceType, showRaw)
			if showPodName {
				// add podname to the beginning of the row
				table[idx] = append([]string{pod.Name}, table[idx]...)
			}
		}
	}

	showTable(table)
	return nil
}

func statsProcessTableRow(container v1.Container, metrics v1.ResourceList, containerType string, resource string, showRaw bool) []string {
	floatfmt := "%f"
	displayValue := ""
	request := ""
	limit := ""
	percentLimit := ""
	percentRequest := ""

	if resource == "cpu" {
		if metrics.Cpu() != nil {
			if showRaw {
				displayValue = metrics.Cpu().String()
			} else {
				displayValue = fmt.Sprintf("%d", metrics.Cpu().MilliValue())
				floatfmt = "%.2f"
			}

			limit = container.Resources.Limits.Cpu().String()
			request = container.Resources.Requests.Cpu().String()
			if cpuVal := metrics.Cpu().AsApproximateFloat64(); cpuVal > 0 {
				// check cpu limits has a value
				if container.Resources.Limits.Cpu().AsApproximateFloat64() == 0 {
					percentLimit = "-"
				} else {
					val := validateFloat64(cpuVal / container.Resources.Limits.Cpu().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
				}
				// check cpu requests has a value
				if container.Resources.Requests.Cpu().AsApproximateFloat64() == 0 {
					percentRequest = "-"
				} else {
					val := validateFloat64(cpuVal / container.Resources.Requests.Cpu().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
				}
			}
		}

	}

	if resource == "memory" {
		if metrics.Memory() != nil {
			if showRaw {
				displayValue = fmt.Sprintf("%d", metrics.Memory().Value())
			} else {
				displayValue = memoryHumanReadable(metrics.Memory().Value())
				floatfmt = "%.2f"
			}

			limit = container.Resources.Limits.Memory().String()
			request = container.Resources.Requests.Memory().String()
			if memVal := metrics.Memory().AsApproximateFloat64(); memVal > 0 {
				// check memory limits has a value
				if container.Resources.Limits.Memory().AsApproximateFloat64() == 0 {
					percentLimit = "-"
				} else {
					val := validateFloat64(memVal / container.Resources.Limits.Memory().AsApproximateFloat64() * 100)
					percentLimit = fmt.Sprintf(floatfmt, val)
				}
				// check memory requests has a value
				if container.Resources.Requests.Memory().AsApproximateFloat64() == 0 {
					percentRequest = "-"
				} else {
					val := validateFloat64(memVal / container.Resources.Requests.Memory().AsApproximateFloat64() * 100)
					percentRequest = fmt.Sprintf(floatfmt, val)
				}
			}
		}
	}

	return []string{
		containerType,
		container.Name,
		displayValue,
		request,
		limit,
		percentRequest,
		percentLimit,
	}
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
