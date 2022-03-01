package plugin

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

//load config for the k8s endpoint
func loadConfig(configFlags *genericclioptions.ConfigFlags) (kubernetes.Clientset, error) {

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return kubernetes.Clientset{}, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return kubernetes.Clientset{}, fmt.Errorf("failed to create clientset: %w", err)
	}

	return *clientset, nil
}

//load config for the metrics endpoint
func loadMetricConfig(configFlags *genericclioptions.ConfigFlags) (metricsclientset.Clientset, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return metricsclientset.Clientset{}, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	metricset, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return metricsclientset.Clientset{}, fmt.Errorf("failed to create clientset: %w", err)
	}

	return *metricset, nil
}

func getNamespace(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) string {
	namespace := ""
	ctx := ""

	if allNamespaces {
		// get/list pods will search all namespaces in the current context
		return ""
	}

	// was a namespace specified on the cmd line
	if len(*configFlags.Namespace) > 0 {
		return *configFlags.Namespace
	}

	// now try to load the current namespace for our context
	clientCfg, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	// if context was suppiled on cmd line use that
	if len(*configFlags.Context) > 0 {
		ctx = *configFlags.Context
	} else {
		ctx = clientCfg.CurrentContext
	}

	namespace = clientCfg.Contexts[ctx].Namespace
	if len(namespace) > 0 {
		return namespace
	}

	return "default"
}

// returns a list of pods or a list with one pod when given a pod name
func getPods(clientSet kubernetes.Clientset, configFlags *genericclioptions.ConfigFlags, podNameList []string, allNamespaces bool) ([]v1.Pod, error) {
	podList := []v1.Pod{}

	namespace := getNamespace(configFlags, allNamespaces)

	if len(podNameList) > 0 {
		for _, podname := range podNameList {
			// single pod
			pod, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
			if err == nil {
				podList = append(podList, []v1.Pod{*pod}...)
			} else {
				return []v1.Pod{}, fmt.Errorf("failed to retrieve pod from server: %w", err)
			}
		}

		return podList, nil
	} else {
		// multi pods
		podList, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			if len(podList.Items) == 0 {
				return []v1.Pod{}, errors.New("No pods found in default namespace.")
			} else {
				return podList.Items, nil
			}
		} else {
			return []v1.Pod{}, fmt.Errorf("failed to retrieve pod list from server: %w", err)
		}
	}

}

//get an array of pod metrics
func getMetricPods(clientSet metricsclientset.Clientset, configFlags *genericclioptions.ConfigFlags, podNameList []string, allNamespaces bool) ([]v1beta1.PodMetrics, error) {
	podList := []v1beta1.PodMetrics{}

	namespace := getNamespace(configFlags, allNamespaces)

	if len(podNameList) > 0 {
		for _, podname := range podNameList {
			// single pod
			pod, err := clientSet.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
			if err == nil {
				podList = append(podList, []v1beta1.PodMetrics{*pod}...)
			} else {
				return []v1beta1.PodMetrics{}, fmt.Errorf("failed to retrieve pod from metrics: %w", err)
			}
		}

		return podList, nil
	} else {
		podList, err := clientSet.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			if len(podList.Items) == 0 {
				return []v1beta1.PodMetrics{}, errors.New("No pods found in default namespace.")
			} else {
				return podList.Items, nil
			}
		} else {
			return []v1beta1.PodMetrics{}, fmt.Errorf("failed to retrieve pod list from metrics: %w", err)
		}
	}
}

//print the array as a table, auto adjusts column widths
func showTable(table map[int][]string) {
	colWidth := make([]int, len(table[0]))
	for _, row := range table {
		for idx, word := range row {
			if colWidth[idx] <= len(word) {
				colWidth[idx] = len(word) + 2
			}
		}
	}
	// fmt.Println("F:showTable:len(table)=", len(table))
	for row := 0; row <= len(table); row++ {
		//for _, row := range table {
		for idx, word := range table[row] {
			if len(word) == 0 {
				word = "-"
			}
			pad := strings.Repeat(" ", colWidth[idx]-len(word))
			fmt.Print(word, pad)
		}
		fmt.Println()
	}
}

//returns a list of memory sizes with their multipacation amount
func memoryGetUnitLst() map[string]int64 {
	// Ki | Mi | Gi | Ti | Pi | Ei = 1024 = 1Ki
	// m "" k | M | G | T | P | E = 1000 = 1k
	var d int64 = 1000 // decimal
	var b int64 = 1024 // binary

	return map[string]int64{
		"Ki": b, "Mi": b * b, "Gi": b * b * b, "Ti": b * b * b * b, "Pi": b * b * b * b * b, "Ei": b * b * b * b * b * b,
		"k": d, "M": d * d, "G": d * d * d, "T": d * d * d * d, "P": d * d * d * d * d, "E": d * d * d * d * d * d,
	}
}

// takes a float and converts to a nearest size with unit discriptor as a string
func memoryHumanReadable(memorySize int64) string {
	var floatfmt string
	power := 100.0
	outVal := ""

	if memorySize == 0 {
		return "0"
	}

	byteList := memoryGetUnitLst()

	for k, v := range byteList {
		if len(k) == 2 {
			size := float64(memorySize) / float64(v)
			val := math.Round(size*power) / power

			remain := int64(math.Round(size*power)) % int64(power)
			if remain == 0 {
				floatfmt = "%d%s"
			} else {
				floatfmt = "%.2f%s"
			}

			// TODO: it works but its clunky and a bit crap, needs work :(
			if val > 0.0 && val <= 900 {
				outVal = fmt.Sprintf(floatfmt, val, k)
			}
			if val > 0.9 && val <= 900 {
				outVal = fmt.Sprintf(floatfmt, val, k)
			}
		}
	}
	return outVal
}

//checks if number is NaN, always returns a valid number
func validateFloat64(number float64) float64 {
	if number != number {
		return 0.0
	}
	return number
}
