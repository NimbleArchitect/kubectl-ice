package plugin

import v1 "k8s.io/api/core/v1"

type Looper interface {
	Apply([]Cell)
	// Build(container v1.ContainerStatus, columnInfo container) ([]Cell, error)
	BuildPod(pod v1.Pod, info *containerInfomation) ([]Cell, error)
	BuildContainerStatus(container v1.ContainerStatus, info *containerInfomation) ([]Cell, error)
	BuildEphemeralContainerStatus(container v1.ContainerStatus, columnInfo *containerInfomation) ([]Cell, error)
}

type RowBuilder struct {
	connection  *Connector
	table       *Table
	columnInfo  *containerInfomation
	commonFlags *commonFlags
}

func (b *RowBuilder) PodLoop(loop Looper) error {
	podList, err := b.connection.GetPods([]string{})
	if err != nil {
		return err
	}

	if b.columnInfo.labelNodeName != "" {
		// columnInfo.labelNodeName = cmd.Flag("node-label").Value.String()
		b.connection.nodeLabels, err = b.connection.GetNodeLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.columnInfo.labelPodName != "" {
		// columnInfo.labelPodName = cmd.Flag("pod-label").Value.String()
		b.connection.podLabels, err = b.connection.GetPodLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.columnInfo.annotationPodName != "" {
		// columnInfo.annotationPodName = cmd.Flag("pod-annotation").Value.String()
		b.connection.podAnnotations, err = b.connection.GetPodAnnotations(podList)
		if err != nil {
			return err
		}
	}

	for _, pod := range podList {
		// p := pod.GetOwnerReferences()
		// for i, a := range p {
		// 	fmt.Println("index:", i)
		// 	fmt.Println("** name:", a.Name)
		// 	fmt.Println("** kind:", a.Kind)
		// }

		b.columnInfo.LoadFromPod(pod)

		//check if we have any labels that need to be shown as columns
		if b.columnInfo.labelNodeName != "" {
			b.columnInfo.labelNodeValue = b.connection.nodeLabels[pod.Spec.NodeName][b.columnInfo.labelNodeName]
		}
		if b.columnInfo.labelPodName != "" {
			b.columnInfo.labelPodValue = b.connection.podLabels[pod.Name][b.columnInfo.labelPodName]
		}
		if b.columnInfo.annotationPodName != "" {
			b.columnInfo.annotationPodName = b.connection.podAnnotations[pod.Name][b.columnInfo.annotationPodName]
		}

		//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
		if b.columnInfo.treeView {
			tblOut, err := loop.BuildPod(pod, b.columnInfo)
			if err != nil {

			}
			rowsOut := b.columnInfo.AddRow(tblOut)
			loop.Apply(rowsOut)
		}

		//now show the container line
		b.columnInfo.containerType = "S"
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(*b.commonFlags, container.Name) {
				continue
			}
			b.columnInfo.containerName = container.Name
			tblOut, err := loop.BuildContainerStatus(container, b.columnInfo)
			if err != nil {

			}
			rowsOut := b.columnInfo.AddRow(tblOut)
			loop.Apply(rowsOut)
			// tblOut := statusBuildRow(container, columnInfo, commonFlagList)
			// columnInfo.ApplyRow(&table, tblOut)
		}

		b.columnInfo.containerType = "I"
		for _, container := range pod.Status.InitContainerStatuses {
			// should the container be processed
			if skipContainerName(*b.commonFlags, container.Name) {
				continue
			}
			b.columnInfo.containerName = container.Name
			tblOut, err := loop.BuildContainerStatus(container, b.columnInfo)
			if err != nil {

			}
			rowsOut := b.columnInfo.AddRow(tblOut)
			loop.Apply(rowsOut)
			// tblOut := statusBuildRow(container, columnInfo)
			// columnInfo.ApplyRow(&table, tblOut)
		}

		b.columnInfo.containerType = "E"
		for _, container := range pod.Status.EphemeralContainerStatuses {
			// should the container be processed
			if skipContainerName(*b.commonFlags, container.Name) {
				continue
			}
			b.columnInfo.containerName = container.Name
			tblOut, err := loop.BuildEphemeralContainerStatus(container, b.columnInfo)
			if err != nil {

			}
			rowsOut := b.columnInfo.AddRow(tblOut)
			loop.Apply(rowsOut)
			// tblOut := statusBuildRow(container, columnInfo)
			// columnInfo.ApplyRow(&table, tblOut)
		}
	}

	return nil
}
