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

	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}

	podList, err := getPods(clientset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	table := Table{}
	table.SetHeader(
		"PODNAME", "CONTAINER", "VOLUME", "TYPE", "BACKING", "SIZE", "RO", "MOUNT-POINT",
	)
	table.SetColumnTypeInt(5)

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
		podVolumes := createVolumeMap(pod.Spec.Volumes)

		containerList := append(pod.Spec.InitContainers, pod.Spec.Containers...)
		for _, container := range containerList {
			for _, mount := range container.VolumeMounts {
				// should the container be processed
				if skipContainerName(commonFlagList, container.Name) {
					continue
				}
				tblOut := volumesBuildRow(container, pod.Name, podVolumes, mount)
				table.AddRow(tblOut...)
			}
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func createVolumeMap(volumes []v1.Volume) map[string]map[string]string {
	podMap := make(map[string]map[string]string)
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

func decodeVolumeType(volType string, volume v1.VolumeSource) map[string]string {
	outMap := make(map[string]string)

	if volType == "" {
		return nil
	}

	outMap["type"] = volType
	outMap["size"] = ""
	outMap["backing"] = ""

	switch volType {
	case "AWSElasticBlockStore":
		outMap["backing"] = volume.AWSElasticBlockStore.VolumeID

	case "AzureDisk":
		outMap["backing"] = volume.AzureDisk.DataDiskURI

	case "AzureFile":
		outMap["backing"] = volume.AzureFile.ShareName

	case "Cinder":
		outMap["backing"] = volume.Cinder.VolumeID

	case "ConfigMap":
		outMap["backing"] = volume.ConfigMap.Name

	case "EmptyDir":
		if volume.EmptyDir.SizeLimit != nil {
			outMap["size"] = volume.EmptyDir.SizeLimit.String()
		}
		outMap["backing"] = string(volume.EmptyDir.Medium)

	case "Ephemeral":
		outMap["backing"] = volume.Ephemeral.VolumeClaimTemplate.Name

	case "FC":
		outMap["backing"] = volume.FC.TargetWWNs[0]

	case "Flocker":
		outMap["backing"] = volume.Flocker.DatasetUUID

	case "GCEPersistentDisk":
		outMap["backing"] = volume.GCEPersistentDisk.PDName

	case "HostPath":
		outMap["backing"] = volume.HostPath.Path

	case "ISCSI":
		outMap["backing"] = volume.ISCSI.IQN

	case "NFS":
		outMap["backing"] = volume.NFS.Server + "/" + volume.NFS.Path

	case "PersistentVolumeClaim":
		outMap["backing"] = volume.PersistentVolumeClaim.ClaimName

	case "PhotonPersistentDisk":
		outMap["backing"] = volume.PhotonPersistentDisk.PdID

	case "PortworxVolume":
		outMap["backing"] = volume.PortworxVolume.VolumeID

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
		outMap["backing"] = tmp

	case "Quobyte":
		outMap["backing"] = volume.Quobyte.Tenant

	case "RBD":
		outMap["backing"] = volume.RBD.RBDImage

	case "Secret":
		outMap["backing"] = volume.Secret.SecretName

	case "StorageOS":
		outMap["backing"] = volume.StorageOS.VolumeNamespace + "/" + volume.StorageOS.VolumeName

	case "VsphereVolume":
		outMap["backing"] = volume.VsphereVolume.VolumePath

	default:
		fmt.Println("ERROR: unknown volume type", volType)
		return nil
	}
	return outMap
}

func volumesBuildRow(container v1.Container, podName string, podVolumes map[string]map[string]string, mount v1.VolumeMount) []Cell {
	var volumeType string
	var size string
	var backing string

	// fmt.Println(volume["name"])
	if podVolumes[mount.Name] != nil {
		volume := podVolumes[mount.Name]
		volumeType = volume["type"]
		size = volume["size"]
		backing = volume["backing"]
	}

	return []Cell{
		NewCellText(podName),
		NewCellText(container.Name),
		NewCellText(mount.Name),
		NewCellText(volumeType),
		NewCellText(backing),
		NewCellText(size),
		NewCellText(fmt.Sprintf("%t", mount.ReadOnly)),
		NewCellText(mount.MountPath),
	}

}
