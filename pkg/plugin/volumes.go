package plugin

import (
	"fmt"
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
	var podname []string
	var showPodName bool = true
	var showVolumeDevice bool

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

	if cmd.Flag("device").Value.String() == "true" {
		showVolumeDevice = true
	}

	table := Table{}
	if !showVolumeDevice {
		table.SetHeader(
			"PODNAME", "CONTAINER", "VOLUME", "TYPE", "BACKING", "SIZE", "RO", "MOUNT-POINT",
		)
	} else {
		table.SetHeader(
			"T", "PODNAME", "CONTAINER", "PVC_NAME", "DEVICE_PATH",
		)
	}

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
		info := containerInfomation{
			podName: pod.Name,
		}

		if !showPodName {
			podVolumes := createVolumeMap(pod.Spec.Volumes)

			containerList := append(pod.Spec.InitContainers, pod.Spec.Containers...)
			for _, container := range containerList {
				info.containerName = container.Name
				for _, mount := range container.VolumeMounts {
					// should the container be processed
					if skipContainerName(commonFlagList, container.Name) {
						continue
					}
					tblOut := volumesBuildRow(info, podVolumes, mount)
					table.AddRow(tblOut...)
				}
			}
		} else {
			info.containerType = "S"
			for _, container := range pod.Spec.Containers {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				info.containerName = container.Name
				for _, mount := range container.VolumeDevices {
					tblOut := mountsBuildRow(info, mount)
					table.AddRow(tblOut...)
				}
			}

			info.containerType = "I"
			for _, container := range pod.Spec.InitContainers {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				info.containerName = container.Name
				for _, mount := range container.VolumeDevices {
					tblOut := mountsBuildRow(info, mount)
					table.AddRow(tblOut...)
				}
			}

			info.containerType = "E"
			for _, container := range pod.Spec.EphemeralContainers {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				info.containerName = container.Name
				for _, mount := range container.VolumeDevices {
					tblOut := mountsBuildRow(info, mount)
					table.AddRow(tblOut...)
				}
			}
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func createVolumeMap(volumes []v1.Volume) map[string]map[string]Cell {
	podMap := make(map[string]map[string]Cell)
	// podVolumes := map[string]map[string]string{}
	for _, vol := range volumes {
		v := reflect.ValueOf(vol.VolumeSource)
		typeOfS := v.Type()
		// fmt.Println("===", v.Kind())

		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				name := fmt.Sprintf("%v", typeOfS.Field(i).Name)
				podMap[vol.Name] = decodeVolumeType(name, vol.VolumeSource)
			}
		}
	}

	return podMap
}

func decodeVolumeType(volType string, volume v1.VolumeSource) map[string]Cell {
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
		//TODO: needs reworking it looks fuggly
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
		fmt.Println("ERROR: unknown volume type", volType)
		return nil
	}
	return outMap
}

func volumesBuildRow(info containerInfomation, podVolumes map[string]map[string]Cell, mount v1.VolumeMount) []Cell {
	var volumeType Cell
	var size Cell
	var backing Cell

	// fmt.Println(volume["name"])
	if podVolumes[mount.Name] != nil {
		volume := podVolumes[mount.Name]
		volumeType = volume["type"]
		size = volume["size"]
		backing = volume["backing"]
	}

	return []Cell{
		NewCellText(info.podName),
		NewCellText(info.containerName),
		NewCellText(mount.Name),
		volumeType,
		backing,
		size,
		NewCellText(fmt.Sprintf("%t", mount.ReadOnly)),
		NewCellText(mount.MountPath),
	}

}

func mountsBuildRow(info containerInfomation, mountInfo v1.VolumeDevice) []Cell {

	return []Cell{
		NewCellText(info.containerType),
		NewCellText(info.podName),
		NewCellText(info.containerName),
		NewCellText(mountInfo.Name),
		NewCellText(mountInfo.DevicePath),
	}
}
