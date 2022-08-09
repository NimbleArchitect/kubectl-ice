package plugin

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var volumesShort = "Display container volumes and mount points"

var volumesDescription = ` Prints configured volume information at the container level, volume type, backing information,
read-write state and mount point are all avaliable, volume size is only available if found in
the pod configuration. If no name is specified the volume information for all pods in the
current namespace are shown.`

var volumesExample = `  # List volumes from containers inside pods from current namespace
  %[1]s volumes

  # List volumes from conttainers output in JSON format
  %[1]s volumes -o json

  # List all container volumes from a single pod
  %[1]s volumes my-pod-4jh36

  # List volumes from all containers named web-container searching all 
  # pods in the current namespace
  %[1]s volumes -c web-container

  # List volumes from container web-container searching all pods in current
  # namespace sorted by volume name in descending order (notice the ! charator)
  %[1]s volumes -c web-container --sort '!VOLUME'

  # List volumes from container web-container searching all pods in current
  # namespace sorted by volume name in ascending order
  %[1]s volumes -c web-container --sort MOUNT-POINT

  # List container volume info from all pods where label app equals web
  %[1]s volumes -l app=web

  # List volumes from all containers where the pod label app is web or mail
  %[1]s volumes -l "app in (web,mail)"`

func Volumes(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {

	log := logger{location: "Volumes"}
	log.Debug("Start")

	loopinfo := volumes{}
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

	if cmd.Flag("device").Value.String() == "true" {
		loopinfo.ShowVolumeDevice = true
	}

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

type volumes struct {
	ShowVolumeDevice bool
}

func (s *volumes) Headers() []string {
	if !s.ShowVolumeDevice {
		return []string{
			"VOLUME",
			"TYPE",
			"BACKING",
			"SIZE",
			"RO",
			"MOUNT-POINT",
		}
	} else {
		return []string{
			"PVC_NAME",
			"DEVICE_PATH",
		}
	}
}

func (s *volumes) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *volumes) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *volumes) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s *volumes) BuildBranch(info BuilderInformation, rows [][]Cell) ([]Cell, error) {
	var out []Cell

	if !s.ShowVolumeDevice {
		out = []Cell{
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
		}
	} else {
		out = []Cell{
			NewCellText(""),
			NewCellText(""),
		}
	}

	return out, nil
}

func (s *volumes) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	Pod := info.Data.pod
	if !s.ShowVolumeDevice {
		podVolumes := s.createVolumeMap(Pod.Spec.Volumes)
		for _, mount := range container.VolumeMounts {
			out = append(out, s.volumesBuildRow(info, podVolumes, mount))
		}
	} else {
		for _, mount := range container.VolumeDevices {
			out = append(out, s.mountsBuildRow(mount))
		}
	}
	return out, nil
}

func (s *volumes) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	if !s.ShowVolumeDevice {
		podVolumes := s.createVolumeMap(info.Data.pod.Spec.Volumes)
		for _, mount := range container.VolumeMounts {
			out = append(out, s.volumesBuildRow(info, podVolumes, mount))
		}
	} else {
		for _, mount := range container.VolumeDevices {
			out = append(out, s.mountsBuildRow(mount))
		}
	}
	return out, nil
}

func (s *volumes) createVolumeMap(volumes []v1.Volume) map[string]map[string]Cell {
	podMap := make(map[string]map[string]Cell)
	// podVolumes := map[string]map[string]string{}
	for _, vol := range volumes {
		v := reflect.ValueOf(vol.VolumeSource)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				name := fmt.Sprintf("%v", typeOfS.Field(i).Name)
				podMap[vol.Name] = s.decodeVolumeType(name, vol.VolumeSource)
			}
		}
	}

	return podMap
}

func (s *volumes) decodeVolumeType(volType string, volume v1.VolumeSource) map[string]Cell {
	outMap := make(map[string]Cell)

	if volType == "" {
		return nil
	}

	outMap["type"] = NewCellText(volType)
	outMap["size"] = Cell{}
	outMap["backing"] = Cell{}

	switch volType {
	case "AWSElasticBlockStore":
		outMap["backing"] = NewCellText(volume.AWSElasticBlockStore.VolumeID)

	case "AzureDisk":
		outMap["backing"] = NewCellText(volume.AzureDisk.DataDiskURI)

	case "AzureFile":
		outMap["backing"] = NewCellText(volume.AzureFile.ShareName)

	case "Cinder":
		outMap["backing"] = NewCellText(volume.Cinder.VolumeID)

	case "ConfigMap":
		outMap["backing"] = NewCellText(volume.ConfigMap.Name)

	case "DownwardAPI":
		str := ""
		sep := ""
		for i, value := range volume.DownwardAPI.Items {
			str += sep + value.Path
			if i == 0 {
				sep = ","
			}
		}
		outMap["backing"] = NewCellText(str)

	case "EmptyDir":
		if volume.EmptyDir.SizeLimit != nil {
			outMap["size"] = NewCellInt(volume.EmptyDir.SizeLimit.String(), volume.EmptyDir.SizeLimit.Value())
		}
		outMap["backing"] = NewCellText(string(volume.EmptyDir.Medium))

	case "Ephemeral":
		outMap["backing"] = NewCellText(volume.Ephemeral.VolumeClaimTemplate.Name)

	case "FC":
		outMap["backing"] = NewCellText(volume.FC.TargetWWNs[0])

	case "Flocker":
		outMap["backing"] = NewCellText(volume.Flocker.DatasetUUID)

	case "GCEPersistentDisk":
		outMap["backing"] = NewCellText(volume.GCEPersistentDisk.PDName)

	case "HostPath":
		outMap["backing"] = NewCellText(volume.HostPath.Path)

	case "ISCSI":
		outMap["backing"] = NewCellText(volume.ISCSI.IQN)

	case "NFS":
		outMap["backing"] = NewCellText(volume.NFS.Server + "/" + volume.NFS.Path)

	case "PersistentVolumeClaim":
		outMap["backing"] = NewCellText(volume.PersistentVolumeClaim.ClaimName)

	case "PhotonPersistentDisk":
		outMap["backing"] = NewCellText(volume.PhotonPersistentDisk.PdID)

	case "PortworxVolume":
		outMap["backing"] = NewCellText(volume.PortworxVolume.VolumeID)

	case "Projected":
		tmp := ""
		// TODO: needs reworking it looks fuggly
		for _, val := range volume.Projected.Sources {
			if val.ConfigMap != nil {
				tmp += val.ConfigMap.Name + ","
			}
		}
		if len(tmp) > 0 {
			tmp = tmp[:len(tmp)-1]
		}
		outMap["backing"] = NewCellText(tmp)

	case "Quobyte":
		outMap["backing"] = NewCellText(volume.Quobyte.Tenant)

	case "RBD":
		outMap["backing"] = NewCellText(volume.RBD.RBDImage)

	case "Secret":
		outMap["backing"] = NewCellText(volume.Secret.SecretName)

	case "StorageOS":
		outMap["backing"] = NewCellText(volume.StorageOS.VolumeNamespace + "/" + volume.StorageOS.VolumeName)

	case "VsphereVolume":
		outMap["backing"] = NewCellText(volume.VsphereVolume.VolumePath)

	default:
		fmt.Fprintln(os.Stderr, "ERROR: unknown volume type", volType)
		return nil
	}
	return outMap
}

func (s *volumes) volumesBuildRow(info BuilderInformation, podVolumes map[string]map[string]Cell, mount v1.VolumeMount) []Cell {
	var cellList []Cell
	var volumeType Cell
	var size Cell
	var backing Cell

	if podVolumes[mount.Name] != nil {
		volume := podVolumes[mount.Name]
		volumeType = volume["type"]
		size = volume["size"]
		backing = volume["backing"]
	}

	cellList = append(cellList,
		NewCellText(mount.Name),
		volumeType,
		backing,
		size,
		NewCellText(fmt.Sprintf("%t", mount.ReadOnly)),
		NewCellText(mount.MountPath))

	return cellList
}

func (s *volumes) mountsBuildRow(mountInfo v1.VolumeDevice) []Cell {
	var cellList []Cell

	cellList = append(cellList,
		NewCellText(mountInfo.Name),
		NewCellText(mountInfo.DevicePath),
	)

	return cellList
}
