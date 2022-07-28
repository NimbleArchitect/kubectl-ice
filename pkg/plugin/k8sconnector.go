package plugin

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	a1 "k8s.io/api/apps/v1"
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
	podList        []v1.Pod         //List of Pods
	replicaList    []a1.ReplicaSet  //list of ReplicaSets
	daemonList     []a1.DaemonSet   //list of DaemonSets
	statefulList   []a1.StatefulSet //list of StatefulSet
	deploymentList []a1.Deployment  //list of Deployments
}

type parentData struct {
	name       string
	kind       string
	namespace  string
	deployment a1.Deployment
	replica    a1.ReplicaSet
	stateful   a1.StatefulSet
	daemon     a1.DaemonSet
	pod        v1.Pod
}

type node struct {
	child     map[string]*node
	name      string
	kind      string
	namespace string
	indent    int
	data      parentData
}

func (n *node) getChild(name string) *node {

	for k, v := range n.child {
		if k == name {
			//return matching child if we have it
			return v
		}
	}

	//if we got here we dont have a match so we create a new entry
	child := node{
		name:  name,
		child: make(map[string]*node),
	}

	// and add it as a child
	n.child[name] = &child

	return n.child[name]
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
	if len(c.podList) == 0 {
		err := c.LoadPods(podNameList)
		return c.podList, err
	}

	if len(podNameList) > 0 {
		err := c.LoadPods(podNameList)
		return c.podList, err
	}

	return c.podList, nil

}

// // returns a list of pods or a list with one pod when given a pod name
// func (c *Connector) GetReplica(replicaNameList []string) ([]a1.ReplicaSet, error) {
// 	if len(c.replicaList) == 0 {
// 		err := c.LoadPods(replicaNameList)
// 		return c.replicaList, err
// 	}

// 	if len(replicaNameList) > 0 {
// 		err := c.LoadReplicaSet(replicaNameList)
// 		return c.replicaList, err
// 	}

// 	return c.replicaList, nil

// }

func (c *Connector) GetPodAnnotations(podList []v1.Pod) (map[string]map[string]string, error) {
	//
	annotationsMap := make(map[string]map[string]string)

	for _, pod := range c.podList {
		podName := pod.Name
		annotations := pod.Annotations
		annotationsMap[podName] = annotations
	}

	return annotationsMap, nil
}

func (c *Connector) GetPodLabels(podList []v1.Pod) (map[string]map[string]string, error) {
	//
	labelMap := make(map[string]map[string]string)

	for _, pod := range c.podList {
		podName := pod.Name
		labels := pod.Labels
		labelMap[podName] = labels
	}

	return labelMap, nil
}

func (c *Connector) GetNodeLabels(podList []v1.Pod) (map[string]map[string]string, error) {
	//
	var nameList []string

	labelMap := make(map[string]map[string]string)
	nodeNames := make(map[string]int)

	for _, pod := range c.podList {
		nodeName := pod.Spec.NodeName
		if _, ok := nodeNames[nodeName]; !ok {
			nodeNames[nodeName] = 1
			nameList = append(nameList, nodeName)
		}
	}

	nodeList, err := c.GetNodes(nameList)
	if err != nil {
		return map[string]map[string]string{}, err
	}

	for _, node := range nodeList {
		name := node.Name
		labels := node.Labels
		labelMap[name] = labels
	}

	return labelMap, nil
}

// returns a list of nodes
func (c *Connector) GetNodes(nodeNameList []string) ([]v1.Node, error) {
	nodeList := []v1.Node{}
	selector := metav1.ListOptions{}

	if len(nodeNameList) > 0 {
		if len(c.Flags.labels) > 0 {
			return []v1.Node{}, fmt.Errorf("error: you cannot specify a node name and a selector together")
		}

		// single node
		for _, nodename := range nodeNameList {
			node, err := c.clientSet.CoreV1().Nodes().Get(context.TODO(), nodename, metav1.GetOptions{})
			if err == nil {
				nodeList = append(nodeList, []v1.Node{*node}...)
			} else {
				return []v1.Node{}, fmt.Errorf("failed to retrieve node from server: %w", err)
			}
		}

		return nodeList, nil
	}

	// multi nodes
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	nodes, err := c.clientSet.CoreV1().Nodes().List(context.TODO(), selector)
	if err == nil {
		if len(nodes.Items) == 0 {
			return []v1.Node{}, errors.New("no nodes found in default namespace")
		}
	} else {
		return []v1.Node{}, fmt.Errorf("failed to retrieve node list from server: %w", err)
	}

	return nodes.Items, nil
}

// SelectMatchingPodSpec select pods to inclue or exclude based on the field in v1.Pods.Spec an operator (!=, ==, =) and a string value to match with
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

	return c.configMapArray[configMap][key]
}

// GetNamespace retrieves the namespace that is currently set as default
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

// SetNamespace sets the namespace to use when searching for pods
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

func (c *Connector) LoadPods(podNameList []string) error {
	podList := []v1.Pod{}
	selector := metav1.ListOptions{}

	namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(podNameList) > 0 {
		if len(c.Flags.labels) > 0 {
			c.podList = []v1.Pod{}
			return fmt.Errorf("error: you cannot specify a pod name and a selector together")
		}

		// single pod
		for _, podname := range podNameList {
			pod, err := c.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
			if err == nil {
				podList = append(podList, []v1.Pod{*pod}...)
			} else {
				c.podList = []v1.Pod{}
				return fmt.Errorf("failed to retrieve pod from server: %w", err)
			}
		}

		c.podList = podList
		return nil
	}

	// multi pods
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	pods, err := c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), selector)
	if err == nil {
		if len(pods.Items) == 0 {
			c.podList = []v1.Pod{}
			return errors.New("no pods found in default namespace")
		} else {
			if len(c.Flags.matchSpecList) > 0 {
				c.podList, err = c.SelectMatchinghPodSpec(pods.Items)
				return err
			} else {
				c.podList = pods.Items
				return nil
			}
		}
	} else {
		c.podList = []v1.Pod{}
		return fmt.Errorf("failed to retrieve pod list from server: %w", err)
	}
}

// GetOwnersList calls GetOwnerReference for each pod and returns a unique list of owner types as the key with an array of pods as the value
func (c *Connector) GetOwnersList() (map[string][]v1.Pod, map[string]string) {
	parentList := map[string][]v1.Pod{}
	typeList := map[string]string{}

	for _, pod := range c.podList {
		ownerRef := pod.GetOwnerReferences()
		if len(ownerRef) == 0 {
			continue
		}

		for _, a := range ownerRef {
			parentList[a.Name] = append(parentList[a.Name], pod)
			typeList[a.Name] = a.Kind
		}
	}

	return parentList, typeList
}

func (c *Connector) LoadReplicaSet(replicaNameList []string, namespace string) error {
	replicaList := []a1.ReplicaSet{}
	selector := metav1.ListOptions{}

	// namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(replicaNameList) > 0 {
		// single pod
		for _, replicaName := range replicaNameList {
			rs, err := c.clientSet.AppsV1().ReplicaSets(namespace).Get(context.TODO(), replicaName, metav1.GetOptions{})
			if err == nil {
				replicaList = append(replicaList, []a1.ReplicaSet{*rs}...)
			} else {
				c.replicaList = []a1.ReplicaSet{}
				return fmt.Errorf("failed to retrieve ReplicaSet from server: %w", err)
			}
		}

		c.replicaList = replicaList
		return nil
	}

	// multi pods
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	rs, err := c.clientSet.AppsV1().ReplicaSets(namespace).List(context.TODO(), selector)

	if err == nil {
		if len(rs.Items) == 0 {
			c.replicaList = []a1.ReplicaSet{}
			return errors.New("no ReplicaSet found in default namespace")
		} else {
			if len(c.Flags.matchSpecList) > 0 {
				// c.replicaList, err = c.SelectMatchinghPodSpec(rs.Items)
				return err
			} else {
				// c.replicaList = pods.Items
				return nil
			}
		}
	} else {
		c.replicaList = []a1.ReplicaSet{}
		return fmt.Errorf("failed to retrieve ReplicaSet list from server: %w", err)
	}
}

func (c *Connector) LoadDeployment(deploymentNameList []string, namespace string) error {
	deploymentList := []a1.Deployment{}
	selector := metav1.ListOptions{}

	// namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(deploymentNameList) > 0 {
		// single pod
		for _, replicaName := range deploymentNameList {
			d, err := c.clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), replicaName, metav1.GetOptions{})
			if err == nil {
				deploymentList = append(deploymentList, []a1.Deployment{*d}...)
			} else {
				c.deploymentList = []a1.Deployment{}
				return fmt.Errorf("failed to retrieve Deployment from server: %w", err)
			}
		}

		c.deploymentList = deploymentList
		return nil
	}

	// multi pods
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	d, err := c.clientSet.AppsV1().Deployments(namespace).List(context.TODO(), selector)

	if err == nil {
		if len(d.Items) == 0 {
			c.deploymentList = []a1.Deployment{}
			return errors.New("no Deployment found in default namespace")
		} else {
			if len(c.Flags.matchSpecList) > 0 {
				// c.deploymentList, err = c.SelectMatchinghPodSpec(rs.Items)
				return err
			} else {
				// c.deploymentList = pods.Items
				return nil
			}
		}
	} else {
		c.deploymentList = []a1.Deployment{}
		return fmt.Errorf("failed to retrieve Deployment list from server: %w", err)
	}
}

func (c *Connector) LoadDaemonSet(daemonNameList []string, namespace string) error {
	daemonList := []a1.DaemonSet{}
	selector := metav1.ListOptions{}

	// namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(daemonNameList) > 0 {
		// single pod
		for _, replicaName := range daemonNameList {
			d, err := c.clientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), replicaName, metav1.GetOptions{})
			if err == nil {
				daemonList = append(daemonList, []a1.DaemonSet{*d}...)
			} else {
				c.daemonList = []a1.DaemonSet{}
				return fmt.Errorf("failed to retrieve DaemonSet from server: %w", err)
			}
		}

		c.daemonList = daemonList
		return nil
	}

	// multi pods
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	d, err := c.clientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), selector)

	if err == nil {
		if len(d.Items) == 0 {
			c.daemonList = []a1.DaemonSet{}
			return errors.New("no DaemonSet found in default namespace")
		} else {
			if len(c.Flags.matchSpecList) > 0 {
				// c.deploymentList, err = c.SelectMatchinghPodSpec(rs.Items)
				return err
			} else {
				// c.deploymentList = pods.Items
				return nil
			}
		}
	} else {
		c.daemonList = []a1.DaemonSet{}
		return fmt.Errorf("failed to retrieve DaemonSet list from server: %w", err)
	}
}

func (c *Connector) LoadStatefulSet(statefulNameList []string, namespace string) error {
	statefulList := []a1.StatefulSet{}
	selector := metav1.ListOptions{}

	// namespace := c.GetNamespace(c.Flags.allNamespaces)

	if len(statefulNameList) > 0 {
		// single pod
		for _, replicaName := range statefulNameList {
			d, err := c.clientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), replicaName, metav1.GetOptions{})
			if err == nil {
				statefulList = append(statefulList, []a1.StatefulSet{*d}...)
			} else {
				c.statefulList = []a1.StatefulSet{}
				return fmt.Errorf("failed to retrieve StatefulSet from server: %w", err)
			}
		}

		c.statefulList = statefulList
		return nil
	}

	// multi pods
	if len(c.Flags.labels) > 0 {
		selector.LabelSelector = c.Flags.labels
	}

	d, err := c.clientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), selector)

	if err == nil {
		if len(d.Items) == 0 {
			c.statefulList = []a1.StatefulSet{}
			return errors.New("no StatefulSet found in default namespace")
		} else {
			if len(c.Flags.matchSpecList) > 0 {
				// c.deploymentList, err = c.SelectMatchinghPodSpec(rs.Items)
				return err
			} else {
				// c.deploymentList = pods.Items
				return nil
			}
		}
	} else {
		c.statefulList = []a1.StatefulSet{}
		return fmt.Errorf("failed to retrieve StatefulSet list from server: %w", err)
	}
}

func (c *Connector) BuildOwnersList() map[string]*node {

	children := make(map[string]*node)
	rootnode := node{child: children}

	for _, pod := range c.podList {
		nodename := pod.Spec.NodeName
		//first create a list with the pod as the first entry
		parentList := []parentData{{
			name:      pod.Name,
			namespace: pod.Namespace,
			kind:      "Pod",
			pod:       pod,
		}}
		oref := pod.GetOwnerReferences()

		//then append each owner to the begining of the list, this way we end up with a list that runs from Node to Pod
		parentList = c.appendParents(parentList, oref, nodename, pod.Namespace)

		// finally we can loop through the above list adding children to the tree where they are needed and using child nodes if they already exist
		current := rootnode
		for i, v := range parentList {
			child := current.getChild(v.name)
			child.kind = v.kind
			child.namespace = v.namespace
			child.indent = i //- len(parentList)
			child.data = v
			current = *child
		}

	}

	return rootnode.child

}

func (c *Connector) appendParents(current []parentData, oref []metav1.OwnerReference, nodename string, namespace string) []parentData {
	//check if parent exists based on kind
	if len(oref) == 0 {
		current = append([]parentData{{
			name: nodename,
			kind: "Node",
		}}, current...)
	}
	for _, v := range oref {
		if v.Kind == "Node" {
			current = append([]parentData{{
				name: v.Name,
				kind: v.Kind,
			}}, current...)
		}
		if v.Kind == "Deployment" {
			err := c.LoadDeployment([]string{v.Name}, namespace)
			if err != nil {
				panic(err)
			}

			n := v.Name
			for _, deployment := range c.deploymentList {
				if n == deployment.Name {
					current = append([]parentData{{
						name:       v.Name,
						kind:       v.Kind,
						namespace:  deployment.Namespace,
						deployment: deployment,
					}}, current...)

					return c.appendParents(current, deployment.GetOwnerReferences(), nodename, namespace)
				}
			}
		}
		if v.Kind == "ReplicaSet" {
			err := c.LoadReplicaSet([]string{v.Name}, namespace)
			if err != nil {
				panic(err)
			}

			n := v.Name
			for _, replica := range c.replicaList {
				if n == replica.Name {
					current = append([]parentData{{
						name:      v.Name,
						kind:      v.Kind,
						namespace: replica.Namespace,
						replica:   replica,
					}}, current...)

					return c.appendParents(current, replica.GetOwnerReferences(), nodename, namespace)
				}
			}
		}
		if v.Kind == "DaemonSet" {
			err := c.LoadDaemonSet([]string{v.Name}, namespace)
			if err != nil {
				panic(err)
			}

			n := v.Name
			for _, daemon := range c.daemonList {
				if n == daemon.Name {
					current = append([]parentData{{
						name:      v.Name,
						kind:      v.Kind,
						namespace: daemon.Namespace,
						daemon:    daemon,
					}}, current...)

					return c.appendParents(current, daemon.GetOwnerReferences(), nodename, namespace)
				}
			}
		}
		if v.Kind == "StatefulSet" {
			err := c.LoadStatefulSet([]string{v.Name}, namespace)
			if err != nil {
				panic(err)
			}

			n := v.Name
			for _, stateful := range c.statefulList {
				if n == stateful.Name {
					current = append([]parentData{{
						name:      v.Name,
						kind:      v.Kind,
						namespace: stateful.Namespace,
						stateful:  stateful,
					}}, current...)

					return c.appendParents(current, stateful.GetOwnerReferences(), nodename, namespace)
				}
			}
		}
	}

	return current
}
