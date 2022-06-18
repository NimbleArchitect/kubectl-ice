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
	var tblHead []string
	var podname []string
	var showPodName bool = true
	var showSELinuxOptions bool

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

	if cmd.Flag("selinux").Value.String() == "true" {
		showSELinuxOptions = true
		tblHead = append(infoTableHead(), "USER", "ROLE", "TYPE", "LEVEL")
	} else {
		tblHead = append(infoTableHead(), "ALLOW_PRIVILEGE_ESCALATION", "PRIVILEGED", "RO_ROOT_FS", "RUN_AS_NON_ROOT", "RUN_AS_USER", "RUN_AS_GROUP")
	}
	table.SetHeader(tblHead...)

	if len(commonFlagList.filterList) >= 1 {
		err = table.SetFilter(commonFlagList.filterList)
		if err != nil {
			return err
		}
	}

	if !showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(2)
	}

	if !commonFlagList.showNamespaceName {
		table.HideColumn(1)
	}

	for _, pod := range podList {
		info := containerInfomation{
			podName:   pod.Name,
			namespace: pod.Namespace,
		}

		info.containerType = "S"
		for _, container := range pod.Spec.Containers {
			var tblOut []Cell
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			info.containerName = container.Name
			if showSELinuxOptions {
				tblOut = seLinuxBuildRow(info, container.SecurityContext, pod.Spec.SecurityContext)
			} else {
				tblOut = securityBuildRow(info, container.SecurityContext, pod.Spec.SecurityContext)
			}
			tblFullRow := append(infoTable(info), tblOut...)
			table.AddRow(tblFullRow...)
		}

		info.containerType = "I"
		for _, container := range pod.Spec.InitContainers {
			var tblOut []Cell
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			info.containerName = container.Name
			if showSELinuxOptions {
				tblOut = seLinuxBuildRow(info, container.SecurityContext, pod.Spec.SecurityContext)
			} else {
				tblOut = securityBuildRow(info, container.SecurityContext, pod.Spec.SecurityContext)
			}
			tblFullRow := append(infoTable(info), tblOut...)
			table.AddRow(tblFullRow...)
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func securityBuildRow(info containerInfomation, csc *v1.SecurityContext, psc *v1.PodSecurityContext) []Cell {

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
		ape,
		p,
		rorfs,
		ranr,
		rau,
		rag,
	}
}

func seLinuxBuildRow(info containerInfomation, csc *v1.SecurityContext, psc *v1.PodSecurityContext) []Cell {

	seLevel := Cell{}
	seRole := Cell{}
	seType := Cell{}
	seUser := Cell{}

	if psc != nil {
		if psc.SELinuxOptions != nil {
			pselinux := psc.SELinuxOptions
			if len(pselinux.Level) > 0 {
				seLevel = NewCellText(pselinux.Level)
			}

			if len(pselinux.Role) > 0 {
				seRole = NewCellText(pselinux.Role)
			}

			if len(pselinux.Type) > 0 {
				seType = NewCellText(pselinux.Type)
			}

			if len(pselinux.User) > 0 {
				seUser = NewCellText(pselinux.User)
			}
		}
	}

	if csc != nil {
		if csc.SELinuxOptions != nil {
			cselinux := psc.SELinuxOptions
			if len(cselinux.Level) > 0 {
				seLevel = NewCellText(cselinux.Level)
			}

			if len(cselinux.Role) > 0 {
				seRole = NewCellText(cselinux.Role)
			}

			if len(cselinux.Type) > 0 {
				seType = NewCellText(cselinux.Type)
			}

			if len(cselinux.User) > 0 {
				seUser = NewCellText(cselinux.User)
			}
		}
	}

	return []Cell{
		seUser,
		seRole,
		seType,
		seLevel,
	}
}
