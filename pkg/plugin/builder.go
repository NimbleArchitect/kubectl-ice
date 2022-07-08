package plugin

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type Looper interface {
	Apply([]Cell)
	// Build(container v1.ContainerStatus, columnInfo container) ([]Cell, error)
	BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error)
	BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([]Cell, error)
	BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([]Cell, error)
}

type RowBuilder struct {
	Connection *Connector
	Table      *Table
	// ColumnInfo         *containerInfomation
	CommonFlags        *commonFlags
	LabelNodeName      string
	labelNodeValue     string
	LabelPodName       string
	labelPodValue      string
	AnnotationPodName  string
	annotationPodValue string
}

type BuilderInformation struct {
	ContainerName string
	ContainerType string
	Namespace     string
	NodeName      string
	PodName       string
	TreeView      bool
}

func (b *RowBuilder) PodLoop(loop Looper) error {
	var nodeLabels map[string]map[string]string
	var podLabels map[string]map[string]string
	var podAnnotations map[string]map[string]string
	var info BuilderInformation

	podList, err := b.Connection.GetPods([]string{})
	if err != nil {
		return err
	}

	if b.LabelNodeName != "" {
		// columnInfo.labelNodeName = cmd.Flag("node-label").Value.String()
		nodeLabels, err = b.Connection.GetNodeLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.LabelPodName != "" {
		// columnInfo.labelPodName = cmd.Flag("pod-label").Value.String()
		podLabels, err = b.Connection.GetPodLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.AnnotationPodName != "" {
		// columnInfo.annotationPodName = cmd.Flag("pod-annotation").Value.String()
		podAnnotations, err = b.Connection.GetPodAnnotations(podList)
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

		info.PodName = pod.Name
		info.Namespace = pod.Namespace
		info.NodeName = pod.Spec.NodeName

		//check if we have any labels that need to be shown as columns
		if b.LabelNodeName != "" {
			b.labelNodeValue = nodeLabels[pod.Spec.NodeName][b.LabelNodeName]
		}
		if b.LabelPodName != "" {
			b.labelPodValue = podLabels[pod.Name][b.LabelPodName]
		}
		if b.AnnotationPodName != "" {
			b.annotationPodValue = podAnnotations[pod.Name][b.AnnotationPodName]
		}

		//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
		if info.TreeView {
			tblOut, err := loop.BuildPod(pod, info)
			if err != nil {

			}
			rowsOut := b.MakeRow(info, tblOut)
			loop.Apply(rowsOut)
		}

		//now show the container line
		info.ContainerType = "S"
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(*b.CommonFlags, container.Name) {
				continue
			}
			info.ContainerName = container.Name
			tblOut, err := loop.BuildContainerStatus(container, info)
			if err != nil {

			}
			rowsOut := b.MakeRow(info, tblOut)
			loop.Apply(rowsOut)
			// tblOut := statusBuildRow(container, columnInfo, commonFlagList)
			// columnInfo.ApplyRow(&table, tblOut)
		}

		info.ContainerType = "I"
		for _, container := range pod.Status.InitContainerStatuses {
			// should the container be processed
			if skipContainerName(*b.CommonFlags, container.Name) {
				continue
			}
			info.ContainerName = container.Name
			tblOut, err := loop.BuildContainerStatus(container, info)
			if err != nil {

			}
			rowsOut := b.MakeRow(info, tblOut)
			loop.Apply(rowsOut)
			// tblOut := statusBuildRow(container, columnInfo)
			// columnInfo.ApplyRow(&table, tblOut)
		}

		info.ContainerType = "E"
		for _, container := range pod.Status.EphemeralContainerStatuses {
			// should the container be processed
			if skipContainerName(*b.CommonFlags, container.Name) {
				continue
			}
			info.ContainerName = container.Name
			tblOut, err := loop.BuildEphemeralContainerStatus(container, info)
			if err != nil {

			}
			rowsOut := b.MakeRow(info, tblOut)
			loop.Apply(rowsOut)
			// tblOut := statusBuildRow(container, columnInfo)
			// columnInfo.ApplyRow(&table, tblOut)
		}
	}

	return nil
}

// MakeRow adds the listed columns to the default columns, outputs
//  the complete row as a list of columns
func (b *RowBuilder) MakeRow(info BuilderInformation, columns ...[]Cell) []Cell {
	rowList := info.GetDefaultCells()

	if b.LabelNodeName != "" {
		rowList = append(rowList, NewCellText(b.labelNodeValue))
	}

	if b.LabelPodName != "" {
		rowList = append(rowList, NewCellText(b.labelPodValue))
	}

	for _, c := range columns {
		rowList = append(rowList, c...)
	}

	return rowList
}

// GetDefaultCells: returns an array of cells prepopulated with the common information
func (b *BuilderInformation) GetDefaultCells() []Cell {
	if b.TreeView {
		return []Cell{
			NewCellText(b.Namespace),
			NewCellText(b.NodeName),
		}
	} else {
		return []Cell{
			NewCellText(b.ContainerType),
			NewCellText(b.Namespace),
			NewCellText(b.NodeName),
			NewCellText(b.PodName),
			NewCellText(b.ContainerName),
		}
	}
}

func (b *BuilderInformation) BuildTreeCell(cellList []Cell) []Cell {
	var namePrefix string

	if b.ContainerType == "S" {
		namePrefix = "Container/"
	}
	if b.ContainerType == "I" {
		namePrefix = "InitContainer/"
	}
	if b.ContainerType == "E" {
		namePrefix = "EphemeralContainer/"
	}

	cellList = append(cellList,
		NewCellText(fmt.Sprint("└─", namePrefix, b.ContainerName)),
	)
	return cellList
}
