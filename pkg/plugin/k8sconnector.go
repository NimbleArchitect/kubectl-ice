package plugin

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Connector struct {
	clientSet      kubernetes.Clientset
	metricSet      metricsclientset.Clientset
	Flags          commonFlags
	configFlags    *genericclioptions.ConfigFlags
	metricFlags    *genericclioptions.ConfigFlags
	configMapArray map[string]map[string]string
	setNameSpace   string
}

//load config for the k8s endpoint
func (c *Connector) LoadConfig(configFlags *genericclioptions.ConfigFlags) error {
	c.clientSet = kubernetes.Clientset{}
	c.configFlags = configFlags
	config, err := configFlags.ToRESTConfig()

	if err != nil {
		return fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}
	c.clientSet = *clientset
	return nil
}

//load config for the metrics endpoint
func (c *Connector) LoadMetricConfig(configFlags *genericclioptions.ConfigFlags) error {
	c.metricSet = metricsclientset.Clientset{}
	c.metricFlags = configFlags
	config, err := configFlags.ToRESTConfig()

	if err != nil {
		return fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	metricset, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset for metrics: %w", err)
	}

	c.metricSet = *metricset
	return nil
}

// returns a list of pods or a list with one pod when given a pod name
func (c *Connector) GetPods(podNameList []string) ([]v1.Pod, error) {
	podList := []v1.Pod{}
	selector := metav1.ListOptions{}

	namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(podNameList) > 0 {
		if len(c.Flags.labels) > 0 {
			return []v1.Pod{}, fmt.Errorf("error: you cannot specify a pod name and a selector together")
		}

		// single pod
		for _, podname := range podNameList {
			pod, err := c.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
			if err == nil {
				podList = append(podList, []v1.Pod{*pod}...)
			} else {
				return []v1.Pod{}, fmt.Errorf("failed to retrieve pod from server: %w", err)
			}
		}

		return podList, nil
	}

	// multi pods
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	pods, err := c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), selector)
	if err == nil {
		if len(pods.Items) == 0 {
			return []v1.Pod{}, errors.New("no pods found in default namespace")
		} else {
			if len(c.Flags.matchSpecList) > 0 {
				return c.SelectMatchinghPodSpec(pods.Items)
			} else {
				return pods.Items, nil
			}
		}
	} else {
		return []v1.Pod{}, fmt.Errorf("failed to retrieve pod list from server: %w", err)
	}

}

// SelectMatchingPodSpec select pods to inclue or eclude based on the field in v1.Pods.Spec an operator (!=, ==, =) and a string value to match with
func (c *Connector) SelectMatchinghPodSpec(pods []v1.Pod) ([]v1.Pod, error) {
	var newPodList []v1.Pod

	//grab and compare the field name to the user suppilied string as the user may have typed all in caps
	includeList := make(map[string]matchValue)

	fields := reflect.VisibleFields(reflect.TypeOf(v1.Pod{}.Spec))
	for _, field := range fields {
		isValid := false

		name := strings.ToUpper(field.Name)
		//restrict to basic types (string, int, bool)
		switch field.Type.String() {
		case "string", "*string":
			fallthrough
		case "int", "*int":
			fallthrough
		case "int32", "*int32":
			fallthrough
		case "int64", "*int64":
			fallthrough
		case "bool", "*bool":
			isValid = true
		}

		if !isValid {
			continue
		}

		if value, ok := c.Flags.matchSpecList[name]; ok {
			includeList[field.Name] = value
		}
	}

	// now we can loop through doing a name lookup with should be faster than searching each name to find a match
	for _, i := range pods {
		fields := reflect.ValueOf(i.Spec)
		for k, v := range includeList {
			field := fields.FieldByName(k)
			fieldString := convertToString(field, field.Interface())
			switch v.operator {
			case "=":
				fallthrough
			case "==":
				if fieldString == v.value {
					newPodList = append(newPodList, i)
				}
			case "!=":
				if fieldString != v.value {
					newPodList = append(newPodList, i)
				}
			default:
				return []v1.Pod{}, errors.New("invalid operator found")
			}
		}

	}

	return newPodList, nil
}

//get an array of pod metrics
func (c *Connector) GetMetricPods(podNameList []string) ([]v1beta1.PodMetrics, error) {
	podList := []v1beta1.PodMetrics{}
	selector := metav1.ListOptions{}

	namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(podNameList) > 0 {
		for _, podname := range podNameList {
			if len(c.Flags.labels) > 0 {
				return []v1beta1.PodMetrics{}, fmt.Errorf("error: you cannot specify a pod name and a selector together")
			}

			// single pod
			pod, err := c.metricSet.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
			if err == nil {
				podList = append(podList, []v1beta1.PodMetrics{*pod}...)
			} else {
				return []v1beta1.PodMetrics{}, fmt.Errorf("failed to retrieve pod from metrics: %w", err)
			}
		}

		return podList, nil
	} else {
		if len(c.Flags.labels) > 0 {
			selector.LabelSelector = c.Flags.labels
		}

		podList, err := c.metricSet.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), selector)
		if err == nil {
			if len(podList.Items) == 0 {
				return []v1beta1.PodMetrics{}, errors.New("no metric info found for pods in namespace")
			} else {
				return podList.Items, nil
			}
		} else {
			return []v1beta1.PodMetrics{}, fmt.Errorf("failed to retrieve pod list from metrics: %w", err)
		}
	}
}

func (c *Connector) GetConfigMaps(configMapName string) (v1.ConfigMap, error) {

	namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(configMapName) == 0 {
		return v1.ConfigMap{}, nil
	}

	cm, err := c.clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err == nil {
		return *cm, nil
	}

	return v1.ConfigMap{}, nil
}

func (c *Connector) GetConfigMapValue(configMap string, key string) string {
	var val map[string]map[string]string

	if len(configMap) <= 0 {
		return ""
	}

	if _, ok := c.configMapArray[configMap]; !ok {
		//fmt.Println("Loadme", configMap)
		cm, err := c.GetConfigMaps(configMap)
		if err != nil {
			c.configMapArray[configMap] = make(map[string]string)
			return ""
		}

		if len(c.configMapArray) > 0 {
			val = c.configMapArray
		} else {
			val = make(map[string]map[string]string)
		}
		val[configMap] = cm.Data
		c.configMapArray = val

	}

	//fmt.Println("===", configMap, " + ", key, " - ", c.configMapArray[configMap][key], "===")
	return c.configMapArray[configMap][key]
}

func (c *Connector) GetNamespace(allNamespaces bool) string {
	namespace := ""
	ctx := ""

	if len(c.setNameSpace) >= 1 {
		return c.setNameSpace
	}

	if allNamespaces {
		// get/list pods will search all namespaces in the current context
		return ""
	}

	// was a namespace specified on the cmd line
	if len(*c.configFlags.Namespace) > 0 {
		return *c.configFlags.Namespace
	}

	// now try to load the current namespace for our context
	clientCfg, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	// if context was suppiled on cmd line use that
	if len(*c.configFlags.Context) > 0 {
		ctx = *c.configFlags.Context
	} else {
		ctx = clientCfg.CurrentContext
	}

	if clientCfg.Contexts[ctx] == nil {
		return "default"
	}

	namespace = clientCfg.Contexts[ctx].Namespace
	if len(namespace) > 0 {
		return namespace
	}

	return "default"
}

func (c *Connector) SetNamespace(namespace string) {
	if len(namespace) >= 1 {
		c.setNameSpace = namespace
	}
}

// convertToString expects a reflect value and the raw interface value and returns the value
// as a string, it also handles pointers correctly
func convertToString(field reflect.Value, value interface{}) string {

	switch value.(type) {
	case *bool:
		if !field.IsNil() {
			return fmt.Sprint(reflect.Indirect(field).Bool())
		}

	case *string:
		if !field.IsNil() {
			return fmt.Sprint(reflect.Indirect(field).String())
		}

	case *int, *int32, *int64:
		if !field.IsNil() {
			return fmt.Sprint(reflect.Indirect(field).Int())
		}
	}

	return fmt.Sprint(value)
}
