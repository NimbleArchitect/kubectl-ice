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

// list details of configured liveness readiness and startup security
func Security(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {

	log := logger{location: "Security"}
	log.Debug("Start")

	loopinfo := security{}
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

	if cmd.Flag("selinux").Value.String() == "true" {
		log.Debug("loopinfo.ShowSELinuxOptions = true")
		loopinfo.ShowSELinuxOptions = true
	}

	if err := builder.Build(&loopinfo); err != nil {
		return err
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type security struct {
	ShowSELinuxOptions bool
}

func (s *security) Headers() []string {
	if s.ShowSELinuxOptions {
		return []string{
			"USER",
			"ROLE",
			"TYPE",
			"LEVEL",
		}
	} else {
		return []string{
			"ALLOW_PRIVILEGE_ESCALATION",
			"PRIVILEGED",
			"RO_ROOT_FS",
			"RUN_AS_NON_ROOT",
			"RUN_AS_USER",
			"RUN_AS_GROUP",
		}
	}
}

func (s *security) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *security) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *security) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s *security) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	var rowOut []Cell

	if s.ShowSELinuxOptions {
		rowOut = make([]Cell, 4)
	} else {
		rowOut = make([]Cell, 6)
	}
	return rowOut, nil
}

func (s *security) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	if s.ShowSELinuxOptions {
		out[0] = s.seLinuxBuildRow(info, container.SecurityContext, info.Data.pod.Spec.SecurityContext)
	} else {
		out[0] = s.securityBuildRow(info, container.SecurityContext, info.Data.pod.Spec.SecurityContext)
	}
	return out, nil
}

func (s *security) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	if s.ShowSELinuxOptions {
		out[0] = s.seLinuxBuildRow(info, container.SecurityContext, info.Data.pod.Spec.SecurityContext)
	} else {
		out[0] = s.securityBuildRow(info, container.SecurityContext, info.Data.pod.Spec.SecurityContext)
	}
	return out, nil
}

func (s *security) securityBuildRow(info BuilderInformation, csc *v1.SecurityContext, psc *v1.PodSecurityContext) []Cell {
	var cellList []Cell
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

	// if info.TreeView {
	// 	cellList = info.BuildTreeCell(cellList)
	// }

	cellList = append(cellList,
		ape,
		p,
		rorfs,
		ranr,
		rau,
		rag,
	)

	return cellList

}

func (s *security) seLinuxBuildRow(info BuilderInformation, csc *v1.SecurityContext, psc *v1.PodSecurityContext) []Cell {
	var cellList []Cell
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

	// if info.TreeView {
	// 	cellList = info.BuildTreeCell(cellList)
	// }

	cellList = append(cellList,
		seUser,
		seRole,
		seType,
		seLevel,
	)

	return cellList
}
