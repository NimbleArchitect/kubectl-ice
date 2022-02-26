package plugin

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

//load config for the k8s endpoint
func loadConfig(configFlags *genericclioptions.ConfigFlags, cmd *cobra.Command) (kubernetes.Clientset, error) {

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return kubernetes.Clientset{}, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	// clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)
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

// returns a list of pods or a list with one pod when given a pod name
func getPods(clientSet kubernetes.Clientset, configFlags *genericclioptions.ConfigFlags, podname string) ([]v1.Pod, error) {
	var err error

	namespace := *configFlags.Namespace
	if podname != "" {
		// single pod
		pod, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
		if err == nil {
			podList := []v1.Pod{*pod}
			return podList, nil
		}
	} else {
		podList, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			return podList.Items, nil
		}
	}
	return nil, err
}

//get an array of pod metrics
func getMetricPods(clientSet metricsclientset.Clientset, configFlags *genericclioptions.ConfigFlags, podname string) ([]v1beta1.PodMetrics, error) {
	var err error

	namespace := *configFlags.Namespace
	if podname != "" {
		// single pod
		pod, err := clientSet.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
		if err == nil {
			podList := []v1beta1.PodMetrics{*pod}
			return podList, nil
		}
	} else {
		podList, err := clientSet.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			return podList.Items, nil
		}
	}
	return nil, err
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
