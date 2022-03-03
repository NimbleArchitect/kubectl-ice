package plugin

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type probeAction struct {
	probeName  string
	action     string
	actionName string
	probe      *v1.Probe
}

//list details of configured liveness readiness and startup probes
func Probes(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true
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

	table := make(map[int][]string)
	table[0] = []string{"CONTAINER", "PROBE", "DELAY", "PERIOD", "TIMEOUT", "SUCCESS", "FAILURE", "CHECK", "ACTION"}

	if showPodName {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {

			// add the probes to our map (if defined) so we can loop through each
			probeList := buildProbeList(container)
			// loop over all probes build the output table and add the podname if multipule pods will be output
			for _, probe := range probeList {
				for _, action := range probe {
					idx++
					table[idx] = probesBuildRow(container, action)
					if showPodName {
						table[idx] = append([]string{pod.Name}, table[idx]...)
					}
				}
			}
		}
	}
	showTable(table)
	return nil

}

func probesBuildRow(container v1.Container, action probeAction) []string {

	return []string{
		container.Name,
		action.probeName,
		fmt.Sprintf("%d", action.probe.InitialDelaySeconds),
		fmt.Sprintf("%d", action.probe.PeriodSeconds),
		fmt.Sprintf("%d", action.probe.TimeoutSeconds),
		fmt.Sprintf("%d", action.probe.SuccessThreshold),
		fmt.Sprintf("%d", action.probe.FailureThreshold),
		action.actionName,
		action.action,
	}
}

//check each type of probe and return a list
func buildProbeList(container v1.Container) map[string][]probeAction {
	probes := make(map[string][]probeAction)
	if container.LivenessProbe != nil {
		probes["liveness"] = buildProbeAction("liveness", container.LivenessProbe)
	}
	if container.ReadinessProbe != nil {
		probes["readiness"] = buildProbeAction("liveness", container.ReadinessProbe)
	}
	if container.StartupProbe != nil {
		probes["startup"] = buildProbeAction("liveness", container.StartupProbe)
	}

	return probes
}

//given a probe return an array of probeAction with the action translated to a string
func buildProbeAction(name string, probe *v1.Probe) []probeAction {
	probeList := []probeAction{}
	item := probeAction{
		probeName: name,
		probe:     probe,
	}

	//translate Exec action
	if probe.Exec != nil {
		item.actionName = "Exec"
		item.action = strings.Join(probe.Exec.Command, " ")
		probeList = append(probeList, item)
	}

	//translate HTTP action
	if probe.HTTPGet != nil {
		item.actionName = "HTTPGet"
		actionStr := ""
		p := probe.HTTPGet
		if len(p.Scheme) > 0 {
			actionStr = strings.ToLower(string(p.Scheme)) + "://"
		}

		if len(p.Host) > 0 {
			actionStr += p.Host
		}

		actionStr += portAsString(p.Port)

		if len(p.Path) > 0 {
			actionStr += p.Path
		}
		item.action = actionStr
		probeList = append(probeList, item)
	}

	//translate GRPC action
	if probe.GRPC != nil {
		item.actionName = "GRPC"
		if probe.GRPC.Service == nil {
			item.action = *probe.GRPC.Service
		}
		if probe.GRPC.Port > 0 {
			item.action += fmt.Sprintf(":%d", probe.GRPC.Port)
		}
		probeList = append(probeList, item)
	}

	//translate TCPSocket action
	if probe.TCPSocket != nil {
		item.actionName = "TCPSocket"
		actionStr := ""
		item.action = probe.TCPSocket.String()
		if len(probe.TCPSocket.Host) > 0 {
			actionStr += probe.TCPSocket.Host
		}
		actionStr += portAsString(probe.TCPSocket.Port)
		item.action = actionStr
		probeList = append(probeList, item)
	}

	return probeList
}

// takes a port object and returns either the number or the name as a string with a proceeding :
// returns empty string if port is empty
func portAsString(port intstr.IntOrString) string {
	//port number provided
	if port.Type == 0 {
		if port.IntVal > 0 {
			return fmt.Sprintf(":%d", port.IntVal)
		} else {
			return ""
		}
	}

	//port name provided
	if port.Type == 1 {
		if len(port.StrVal) > 0 {
			return ":" + port.StrVal
		} else {
			return ""
		}
	}

	return ""
}
