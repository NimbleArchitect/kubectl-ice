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

// list details of configured liveness readiness and startup probes
func Probes(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {

	log := logger{location: "Probes"}
	log.Debug("Start")

	loopinfo := probes{}
	builder := RowBuilder{}
	builder.LoopSpec = true
	builder.ShowInitContainers = true
	builder.PodName = args

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList
	builder.Connection = &connect

	builder.SetFlagsFrom(commonFlagList)

	table := Table{}
	builder.Table = &table
	builder.ShowTreeView = commonFlagList.showTreeView

	if err := builder.Build(&loopinfo); err != nil {
		return err
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type probes struct {
}

func (s *probes) Headers() []string {
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

func (s *probes) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *probes) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *probes) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s *probes) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	out := []Cell{
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
		NewCellText(""),
	}
	return out, nil
}

func (s *probes) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	probeList := s.buildProbeList(container)
	for _, probe := range probeList {
		for _, action := range probe {
			out = append(out, s.probesBuildRow(info, action))
		}
	}
	return out, nil
}

func (s *probes) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	return out, nil
}

func (s *probes) probesBuildRow(info BuilderInformation, action probeAction) []Cell {
	var cellList []Cell

	// if info.TreeView {
	// 	cellList = info.BuildTreeCell(cellList)
	// }

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

// check each type of probe and return a list
func (s *probes) buildProbeList(container v1.Container) map[string][]probeAction {
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

// given a probe return an array of probeAction with the action translated to a string
func (s *probes) buildProbeAction(name string, probe *v1.Probe) []probeAction {
	probeList := []probeAction{}
	item := probeAction{
		probeName: name,
		probe:     probe,
	}

	// translate Exec action
	if probe.Exec != nil {
		item.actionName = "Exec"
		item.action = strings.Join(probe.Exec.Command, " ")
		probeList = append(probeList, item)
	}

	// translate HTTP action
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

	// translate GRPC action
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

	// translate TCPSocket action
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
