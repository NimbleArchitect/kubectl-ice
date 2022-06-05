package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var securityShort = "Shows details of configured container security settings"

var securityDescription = ` View SecurityContext configuration that has been applied to the containers. Shows 
runAsUser and runAsGroup fields among others.
`

var securityExample = `  # List container security info from pods
  %[1]s security

  # List container security info from pods output in JSON format
  %[1]s security -o json

  # List container security info from a single pod
  %[1]s security my-pod-4jh36

  # List security info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s security -c web-container

  # List security info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s security -c web-container --sort '!CONTAINER'

  # List security info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s security -c web-container --sort PODNAME

  # List container security info from all pods where label app matches web
  %[1]s security -l app=web

  # List container security info from all pods where the pod label app is either web or mail
  %[1]s security -l "app in (web,mail)"`

//list details of configured liveness readiness and startup security
func Security(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			showPodName = false
		}
	}
	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList

	podList, err := connect.GetPods(podname)
	if err != nil {
		return err
	}

	table := Table{}
	table.SetHeader(
		"T", "PODNAME", "CONTAINER", "ALLOW_PRIVILEGE_ESCALATION", "PRIVILEGED", "RO_ROOT_FS", "RUN_AS_NON_ROOT", "RUN_AS_USER", "RUN_AS_GROUP",
	)

	if len(commonFlagList.filterList) >= 1 {
		err = table.SetFilter(commonFlagList.filterList)
		if err != nil {
			return err
		}
	}

	if !showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(0)
	}

	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := securityBuildRow(container, pod.Name, "S", pod.Spec.SecurityContext)
			table.AddRow(tblOut...)
		}
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := securityBuildRow(container, pod.Name, "I", pod.Spec.SecurityContext)
			table.AddRow(tblOut...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func securityBuildRow(container v1.Container, podName string, containerType string, psc *v1.PodSecurityContext) []Cell {

	ape := Cell{}
	p := Cell{}
	rorfs := Cell{}
	ranr := Cell{}
	rau := Cell{}
	rag := Cell{}

	if psc != nil {
		if psc.RunAsNonRoot != nil {
			ranr = NewCellText(fmt.Sprintf("%t", *psc.RunAsNonRoot))
		}

		if psc.RunAsUser != nil {
			rau = NewCellInt(fmt.Sprintf("%d", *psc.RunAsUser), *psc.RunAsUser)
		}

		if psc.RunAsGroup != nil {
			rag = NewCellInt(fmt.Sprintf("%d", *psc.RunAsGroup), *psc.RunAsGroup)
		}
	}

	csc := container.SecurityContext
	if csc != nil {
		if csc.AllowPrivilegeEscalation != nil {
			ape = NewCellText(fmt.Sprintf("%t", *csc.AllowPrivilegeEscalation))
		}

		if csc.Privileged != nil {
			p = NewCellText(fmt.Sprintf("%t", *csc.Privileged))
		}

		if csc.ReadOnlyRootFilesystem != nil {
			rorfs = NewCellText(fmt.Sprintf("%t", *csc.ReadOnlyRootFilesystem))
		}

		if csc.RunAsNonRoot != nil {
			ranr = NewCellText(fmt.Sprintf("%t", *csc.RunAsNonRoot))
		}

		if csc.RunAsUser != nil {
			rau = NewCellInt(fmt.Sprintf("%d", *csc.RunAsUser), *csc.RunAsUser)
		}

		if csc.RunAsGroup != nil {
			rag = NewCellInt(fmt.Sprintf("%d", *psc.RunAsGroup), *csc.RunAsGroup)
		}
	}

	return []Cell{
		NewCellText(containerType),
		NewCellText(podName),
		NewCellText(container.Name),
		ape,
		p,
		rorfs,
		ranr,
		rau,
		rag,
	}
}
