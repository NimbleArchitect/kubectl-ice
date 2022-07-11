package plugin

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var probesShort = "Shows details of configured startup, readiness and liveness probes of each container"

var probesDescription = ` Prints details of the currently configured startup, liveness and rediness probes for each 
container. Details like the delay timeout and action are printed along with the configured probe
type. If no name is specified the container probe details of all pods in the current namespace
are shown.`

var probesExample = `  # List containers probe info from pods
  %[1]s probes

  # List container probe info from pods output in JSON format
  %[1]s probes -o json

  # List container probe info from a single pod
  %[1]s probes my-pod-4jh36

  # List probe info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s probes -c web-container

  # List probe info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s probes -c web-container --sort '!CONTAINER'

  # List probe info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s probes -c web-container --sort PODNAME

  # List container probe info from all pods where label app matches web
  %[1]s probes -l app=web

  # List container probe info from all pods where the pod label app is either web or mail
  %[1]s probes -l "app in (web,mail)"`

type probeAction struct {
	probeName  string
	action     string
	actionName string
	probe      *v1.Probe
}

//list details of configured liveness readiness and startup probes
func Probes(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var columnInfo containerInfomation
	// var tblHead []string
	var podname []string
	// var showPodName bool = true
	// var nodeLabels map[string]map[string]string
	// var podLabels map[string]map[string]string

	log := logger{location: "Probes"}
	log.Debug("Start")

	loopinfo := probes{}
	builder := RowBuilder{}
	builder.LoopSpec = true
	builder.ShowPodName = true
	builder.ShowInitContainers = true

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			log.Debug("builder.ShowPodName = false")
			builder.ShowPodName = false
		}
	}
	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList
	builder.CommonFlags = commonFlagList

	builder.Connection = &connect

	// podList, err := connect.GetPods(podname)
	// if err != nil {
	// 	return err
	// }

	if cmd.Flag("node-label").Value.String() != "" {
		label := cmd.Flag("node-label").Value.String()
		log.Debug("builder.LabelNodeName =", label)
		builder.LabelNodeName = label
	}

	if cmd.Flag("pod-label").Value.String() != "" {
		label := cmd.Flag("pod-label").Value.String()
		log.Debug("builder.LabelPodName =", label)
		builder.LabelPodName = label
	}

	table := Table{}
	builder.Table = &table
	columnInfo.table = &table
	builder.ShowTreeView = commonFlagList.showTreeView

	// tblHead = columnInfo.GetDefaultHead()
	// if commonFlagList.showTreeView {
	// 	// we have to control the name when displaying a tree view as the table
	// 	//  object dosent have the extra info to be able to process it
	// 	tblHead = append(tblHead, "NAME")
	// }

	// tblHead = append(tblHead, "PROBE", "DELAY", "PERIOD", "TIMEOUT", "SUCCESS", "FAILURE", "CHECK", "ACTION")
	// table.SetHeader(tblHead...)

	// if len(commonFlagList.filterList) >= 1 {
	// 	err = table.SetFilter(commonFlagList.filterList)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// commonFlagList.showPodName = showPodName
	// columnInfo.SetVisibleColumns(table, commonFlagList)

	builder.BuildRows(loopinfo)

	// for _, pod := range podList {
	// 	columnInfo.LoadFromPod(pod)

	// 	if columnInfo.labelNodeName != "" {
	// 		columnInfo.labelNodeValue = nodeLabels[pod.Spec.NodeName][columnInfo.labelNodeName]
	// 	}
	// 	if columnInfo.labelPodName != "" {
	// 		columnInfo.labelPodValue = podLabels[pod.Name][columnInfo.labelPodName]
	// 	}

	// 	//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
	// 	if columnInfo.treeView {
	// 		tblOut := podProbesBuildRow(pod, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 	} else {
	// 		// we force hide the container type column as probes can only be set in standard conatiners
	// 		table.HideColumn(0)
	// 	}

	// 	columnInfo.containerType = "S"
	// 	for _, container := range pod.Spec.Containers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		// add the probes to our map (if defined) so we can loop through each
	// 		probeList := buildProbeList(container)
	// 		// loop over all probes build the output table and add the podname if multipule pods will be output
	// 		for _, probe := range probeList {
	// 			for _, action := range probe {
	// 				tblOut := probesBuildRow(columnInfo, action)
	// 				columnInfo.ApplyRow(&table, tblOut)
	// 				// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 				// table.AddRow(tblFullRow...)
	// 			}
	// 		}
	// 	}
	// }

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type probes struct {
}

func (s probes) Headers() []string {
	return []string{
		"PROBE",
		"DELAY",
		"PERIOD",
		"TIMEOUT",
		"SUCCESS",
		"FAILURE",
		"CHECK",
		"ACTION",
	}
}

func (s probes) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s probes) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s probes) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s probes) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.PodName)), //name
		NewCellText(""),
	}, nil
}

func (s probes) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	probeList := s.buildProbeList(container)
	for _, probe := range probeList {
		for _, action := range probe {
			out = append(out, s.probesBuildRow(info, action))
		}
	}
	return out, nil
}

func (s probes) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	return out, nil
}

func (s probes) podProbesBuildRow(pod v1.Pod, info containerInfomation) []Cell {

	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.podName)), //name
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		// NewCellText(""),
	}
}

func (s probes) probesBuildRow(info BuilderInformation, action probeAction) []Cell {
	var cellList []Cell

	if info.TreeView {
		cellList = info.BuildTreeCell(cellList)
	}

	cellList = append(cellList,
		NewCellText(action.probeName),
		NewCellInt(fmt.Sprintf("%d", action.probe.InitialDelaySeconds), int64(action.probe.InitialDelaySeconds)),
		NewCellInt(fmt.Sprintf("%d", action.probe.PeriodSeconds), int64(action.probe.PeriodSeconds)),
		NewCellInt(fmt.Sprintf("%d", action.probe.TimeoutSeconds), int64(action.probe.TimeoutSeconds)),
		NewCellInt(fmt.Sprintf("%d", action.probe.SuccessThreshold), int64(action.probe.SuccessThreshold)),
		NewCellInt(fmt.Sprintf("%d", action.probe.FailureThreshold), int64(action.probe.FailureThreshold)),
		NewCellText(action.actionName),
		NewCellText(action.action),
	)

	return cellList
}

//check each type of probe and return a list
func (s probes) buildProbeList(container v1.Container) map[string][]probeAction {
	probes := make(map[string][]probeAction)
	if container.LivenessProbe != nil {
		probes["liveness"] = s.buildProbeAction("liveness", container.LivenessProbe)
	}
	if container.ReadinessProbe != nil {
		probes["readiness"] = s.buildProbeAction("readiness", container.ReadinessProbe)
	}
	if container.StartupProbe != nil {
		probes["startup"] = s.buildProbeAction("liveness", container.StartupProbe)
	}

	return probes
}

//given a probe return an array of probeAction with the action translated to a string
func (s probes) buildProbeAction(name string, probe *v1.Probe) []probeAction {
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
