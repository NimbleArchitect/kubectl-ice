package plugin

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Volumes(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true
	var idx int

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

	commonFlagList := processCommonFlags(cmd)

	podList, err := getPods(clientset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	table := make(map[int][]string)
	table[0] = []string{"CONTAINER", "VOLUME", "TYPE", "BACKING", "SIZE", "RO", "MOUNT-POINT"}

	if showPodName {
		// we need to add the pod name to the table
		table[0] = append([]string{"PODNAME"}, table[0]...)
	}

	for _, pod := range podList {
		podVolumes := createVolumeMap(pod.Spec.Volumes)

		containerList := append(pod.Spec.InitContainers, pod.Spec.Containers...)
		for _, container := range containerList {
			for _, mount := range container.VolumeMounts {
				idx++
				table[idx] = volumesBuildRow(container, podVolumes, mount)
				if showPodName {
					table[idx] = append([]string{pod.Name}, table[idx]...)
				}
			}
		}
	}
	showTable(table)
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

func volumesBuildRow(container v1.Container, podVolumes map[string]map[string]string, mount v1.VolumeMount) []string {
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

	return []string{
		container.Name,
		mount.Name,
		volumeType,
		backing,
		size,
		fmt.Sprintf("%t", mount.ReadOnly),
		mount.MountPath,
	}

}
