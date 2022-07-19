package plugin

import (
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var environmentShort = "List the env name and value for each container"

var environmentDescription = ` Print the the environment variables used in running containers in a pod, single pods
and containers can be selected by name. If no name is specified the environment details of all pods in
the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`

var environmentExample = `  # List containers env info from pods
  %[1]s env

  # List container env info from pods output in JSON format
  %[1]s env -o json

  # List container env info from a single pod
  %[1]s env my-pod-4jh36

  # List env info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s env -c web-container

  # List env info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s env -c web-container --sort '!CONTAINER'

  # List env info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s env -c web-container --sort PODNAME

  # List container env info from all pods where label app matches web
  %[1]s env -l app=web

  # List container env info from all pods where the pod label app is either web or mail
  %[1]s env -l "app in (web,mail)"`

func Environment(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	log := logger{location: "Environment"}
	log.Debug("Start")

	loopinfo := environment{}
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

	//we need the connection details so we can translate the environment variables
	loopinfo.Connection = &connect

	if cmd.Flag("translate").Value.String() == "true" {
		loopinfo.TranslateConfigMap = true
	}

	table := Table{}
	builder.Table = &table
	builder.ShowTreeView = commonFlagList.showTreeView
	builder.Build(loopinfo)

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type environment struct {
	Connection         *Connector
	TranslateConfigMap bool
}

func (s environment) Headers() []string {
	return []string{
		"NAME", "VALUE",
	}
}

func (s environment) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s environment) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s environment) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s environment) BuildBranch(info BuilderInformation, podList []v1.Pod) ([][]Cell, error) {
	out := []Cell{
		NewCellText(""),
		NewCellText(""),
	}
	return [][]Cell{out}, nil
}

// func podStatsProcessBuildRow(pod v1.Pod, info containerInfomation) []Cell {
func (s environment) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(""),
		NewCellText(""),
	}, nil
}

func (s environment) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	allRows := s.buildEnvFromContainer(container)
	for _, envRow := range allRows {
		out = append(out, s.envBuildRow(info, envRow, s.Connection, s.TranslateConfigMap))
	}
	return out, nil
}

func (s environment) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	allRows := s.buildEnvFromEphemeral(container)
	for _, envRow := range allRows {
		out = append(out, s.envBuildRow(info, envRow, s.Connection, s.TranslateConfigMap))
	}
	return out, nil
}

func (s environment) envBuildRow(info BuilderInformation, env v1.EnvVar, connect *Connector, translate bool) []Cell {
	var envKey, envValue string
	var configName string
	var key string

	envKey = env.Name
	if len(env.Value) == 0 {
		if env.ValueFrom.ConfigMapKeyRef != nil {
			configName = env.ValueFrom.ConfigMapKeyRef.LocalObjectReference.Name
			key = env.ValueFrom.ConfigMapKeyRef.Key
			envValue = "CONFIGMAP:" + configName + " KEY:" + key
		}

		if env.ValueFrom.SecretKeyRef != nil {
			configName = env.ValueFrom.SecretKeyRef.LocalObjectReference.Name
			key = env.ValueFrom.SecretKeyRef.Key
			envValue = "SECRETMAP:" + configName + " KEY:" + key
			translate = false //never translate secrets
		}

		if env.ValueFrom.FieldRef != nil {
			configName = env.ValueFrom.FieldRef.FieldPath
			envValue = "FIELDREF:" + configName
			translate = false //we cant translate FieldRef at the minute
		}

		if env.ValueFrom.ResourceFieldRef != nil {
			configName = env.ValueFrom.ResourceFieldRef.Resource
			envValue = "RESOURCE:" + configName
			translate = false //we cant translate resourceFieldRef at the moment
		}

		if translate {
			envValue = connect.GetConfigMapValue(configName, key)
		}

	} else {
		envValue = env.Value
	}

	return []Cell{
		NewCellText(envKey),
		NewCellText(envValue),
	}
}

func (s environment) buildEnvFromContainer(container v1.Container) []v1.EnvVar {
	if len(container.Env) == 0 {
		return []v1.EnvVar{}
	}
	return container.Env
}

func (s environment) buildEnvFromEphemeral(container v1.EphemeralContainer) []v1.EnvVar {
	if len(container.Env) == 0 {
		return []v1.EnvVar{}
	}
	return container.Env
}
