package plugin

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type Looper interface {
	// Build(container v1.ContainerStatus, columnInfo container) ([]Cell, error)
	BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error)
	BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error)
	BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error)
	BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error)
	Headers() []string
	HideColumns(info BuilderInformation) []int
}

type RowBuilder struct {
	Connection *Connector
	Table      *Table
	// ColumnInfo         *containerInfomation
	CommonFlags        commonFlags
	LoopStatus         bool
	LoopSpec           bool
	LabelNodeName      string
	labelNodeValue     string
	LabelPodName       string
	labelPodValue      string
	AnnotationPodName  string
	annotationPodValue string
	ShowTreeView       bool
	ShowPodName        bool
	ShowInitContainers bool
	ShowContainerType  bool
	// ShowDetail         bool
	// ShowNamespaceName  bool
	// ShowNodeName       bool
	FilterList       map[string]matchValue // used to filter out rows form the table during Print function
	info             BuilderInformation
	DefaultHeaderLen int
}

type BuilderInformation struct {
	ContainerName string
	ContainerType string
	Namespace     string
	NodeName      string
	PodName       string
	TreeView      bool
}

func (b *RowBuilder) BuildRows(loop Looper) error {
	log := logger{location: "RowBuilder:BuildRows"}
	log.Debug("Start")

	if b.ShowTreeView {
		log.Debug("b.info.TreeView = true")
		b.info.TreeView = true
	}

	if !b.ShowContainerType {
		b.ShowContainerType = b.CommonFlags.showContainerType
	}

	err := b.LoadHeaders(loop)
	if err != nil {
		return err
	}

	err = b.PodLoop(loop)
	if err != nil {
		return err
	}

	// sorting by column breaks the tree view also previous is not valid so we sliently skip those actions
	if !b.ShowTreeView {
		if err := b.Table.SortByNames(b.CommonFlags.sortList...); err != nil {
			return err
		}
	}

	return nil
}

func (b *RowBuilder) LoadHeaders(loop Looper) error {
	var tblHead []string
	var hideColumns []int

	log := logger{location: "RowBuilder:LoadHeaders"}
	log.Debug("Start")

	tblHead = b.GetDefaultHead()
	defaultHeaderLen := len(tblHead)
	log.Debug("len(defaultHeaderLen) =", defaultHeaderLen)

	b.DefaultHeaderLen = defaultHeaderLen
	log.Debug("b.info.TreeView =", b.info.TreeView)
	if b.info.TreeView {
		tblHead = append(tblHead, "NAME")
	}
	hideColumns = loop.HideColumns(b.info)

	tblHead = append(tblHead, loop.Headers()...)
	log.Debug("len(tblHead) =", len(tblHead))
	log.Debug("tblHead =", tblHead)
	b.Table.SetHeader(tblHead...)

	log.Debug("len(b.FilterList) =", len(b.FilterList))
	if len(b.FilterList) >= 1 {
		err := b.Table.SetFilter(b.FilterList)
		if err != nil {
			return err
		}
	}

	b.SetVisibleColumns()

	log.Debug("len(hideColumns) =", len(hideColumns))
	for _, id := range hideColumns {
		b.Table.HideColumn(defaultHeaderLen + id)
	}

	return nil
}

func (b *RowBuilder) SetVisibleColumns() {
	log := logger{location: "RowBuilder:SetVisibleColumns"}
	log.Debug("Start")

	if !b.ShowContainerType {
		b.Table.HideColumn(0)
	}

	if b.info.TreeView {
		//only hide the nodename and namespace, podname is always show in tree view
		if !b.CommonFlags.showNamespaceName {
			b.Table.HideColumn(1)
		}

		if !b.CommonFlags.showNodeName {
			b.Table.HideColumn(2)
		}
		return
	}

	if !b.CommonFlags.showNamespaceName {
		b.Table.HideColumn(1)
	}

	if !b.CommonFlags.showNodeName {
		b.Table.HideColumn(2)
	}

	if !b.ShowPodName {
		// we need to hide the pod name in the table
		b.Table.HideColumn(3)
	}

}

func (b *RowBuilder) PodLoop(loop Looper) error {
	var nodeLabels map[string]map[string]string
	var podLabels map[string]map[string]string
	var podAnnotations map[string]map[string]string

	log := logger{location: "RowBuilder:PodLoop"}
	log.Debug("Start")

	podList, err := b.Connection.GetPods([]string{})
	if err != nil {
		return err
	}

	if b.LabelNodeName != "" {
		log.Debug("b.LabelNodeName", b.LabelNodeName)
		// columnInfo.labelNodeName = cmd.Flag("node-label").Value.String()
		nodeLabels, err = b.Connection.GetNodeLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.LabelPodName != "" {
		log.Debug("b.LabelPodName", b.LabelPodName)
		// columnInfo.labelPodName = cmd.Flag("pod-label").Value.String()
		podLabels, err = b.Connection.GetPodLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.AnnotationPodName != "" {
		log.Debug("b.AnnotationPodName", b.AnnotationPodName)
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
		log.Debug("pod.Name =", pod.Name)
		b.info.PodName = pod.Name
		b.info.Namespace = pod.Namespace
		b.info.NodeName = pod.Spec.NodeName

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
		log.Debug("b.info.TreeView", b.info.TreeView)
		if b.info.TreeView {
			b.info.ContainerType = "P"
			tblOut, err := loop.BuildPod(pod, b.info)
			if err != nil {

			}
			// log.Debug("len(tblOut)", len(tblOut))
			rowsOut := b.MakeRow(b.info, tblOut)
			log.Debug("rowsOut =", rowsOut)
			b.Table.AddRow(rowsOut...)
		}

		if b.ShowInitContainers {
			log.Debug("loop init ContainerStatuses")
			b.info.ContainerType = "I"
			if b.LoopStatus {
				for _, container := range pod.Status.InitContainerStatuses {
					// should the container be processed
					log.Debug("processing -", container.Name)
					if skipContainerName(b.CommonFlags, container.Name) {
						continue
					}
					b.info.ContainerName = container.Name
					allRows, err := loop.BuildContainerStatus(container, b.info)
					if err != nil {

					}
					for _, row := range allRows {
						rowsOut := b.MakeRow(b.info, row)
						b.Table.AddRow(rowsOut...)
					}
				}
			}

			if b.LoopSpec {
				for _, container := range pod.Spec.InitContainers {
					// should the container be processed
					log.Debug("processing -", container.Name)
					if skipContainerName(b.CommonFlags, container.Name) {
						continue
					}
					b.info.ContainerName = container.Name
					allRows, err := loop.BuildContainerSpec(container, b.info)
					if err != nil {

					}
					for _, row := range allRows {
						rowsOut := b.MakeRow(b.info, row)
						b.Table.AddRow(rowsOut...)
					}
				}
			}
		}

		//now show the container line
		log.Debug("loop standard ContainerStatuses")
		b.info.ContainerType = "S"
		if b.LoopStatus {
			for _, container := range pod.Status.ContainerStatuses {
				// should the container be processed
				if skipContainerName(b.CommonFlags, container.Name) {
					continue
				}
				log.Debug("processing -", container.Name)
				b.info.ContainerName = container.Name
				allRows, err := loop.BuildContainerStatus(container, b.info)
				if err != nil {

				}
				for _, row := range allRows {
					rowsOut := b.MakeRow(b.info, row)
					b.Table.AddRow(rowsOut...)
				}
			}
		}

		if b.LoopSpec {
			for _, container := range pod.Spec.Containers {
				// should the container be processed
				if skipContainerName(b.CommonFlags, container.Name) {
					continue
				}
				log.Debug("processing -", container.Name)
				b.info.ContainerName = container.Name
				allRows, err := loop.BuildContainerSpec(container, b.info)
				if err != nil {

				}
				for _, row := range allRows {
					rowsOut := b.MakeRow(b.info, row)
					b.Table.AddRow(rowsOut...)
				}
			}
		}

		log.Debug("loop ephemeral ContainerStatuses")
		b.info.ContainerType = "E"

		if b.LoopStatus {
			for _, container := range pod.Status.EphemeralContainerStatuses {
				// should the container be processed
				if skipContainerName(b.CommonFlags, container.Name) {
					continue
				}
				log.Debug("processing -", container.Name)
				b.info.ContainerName = container.Name
				allRows, err := loop.BuildContainerStatus(container, b.info)
				if err != nil {

				}
				for _, row := range allRows {
					rowsOut := b.MakeRow(b.info, row)
					b.Table.AddRow(rowsOut...)
				}
			}
		}

		if b.LoopSpec {
			for _, container := range pod.Spec.EphemeralContainers {
				// should the container be processed
				if skipContainerName(b.CommonFlags, container.Name) {
					continue
				}
				log.Debug("processing -", container.Name)
				b.info.ContainerName = container.Name
				allRows, err := loop.BuildEphemeralContainerSpec(container, b.info)
				if err != nil {

				}
				for _, row := range allRows {
					rowsOut := b.MakeRow(b.info, row)
					b.Table.AddRow(rowsOut...)
				}
			}
		}
	}

	return nil
}

// MakeRow adds the listed columns to the default columns, outputs
//  the complete row as a list of columns
func (b *RowBuilder) MakeRow(info BuilderInformation, columns ...[]Cell) []Cell {
	log := logger{location: "RowBuilder:MakeRow"}
	log.Debug("Start")

	rowList := b.GetDefaultCells()

	if b.LabelNodeName != "" {
		rowList = append(rowList, NewCellText(b.labelNodeValue))
	}

	if b.LabelPodName != "" {
		rowList = append(rowList, NewCellText(b.labelPodValue))
	}

	for _, c := range columns {
		rowList = append(rowList, c...)
	}

	log.Debug("len(rowList) =", len(rowList))
	return rowList
}

// GetDefaultHead: returns the common headers in order
func (b *RowBuilder) GetDefaultHead() []string {
	log := logger{location: "RowBuilder:GetDefaultHead"}
	log.Debug("Start")

	var headList []string

	log.Debug("b.info.TreeView =", b.info.TreeView)
	if b.info.TreeView {
		//in tree view we only create the namespace and nodename columns, the name colume is created outside of this
		// function so we have full control over its contents
		headList = []string{
			"T", "NAMESPACE", "NODE",
		}
	} else {
		headList = []string{
			"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER",
		}
	}

	if b.LabelNodeName != "" {
		log.Debug("LabelNodeName =", b.LabelNodeName)
		headList = append(headList, b.LabelNodeName)
	}

	if b.LabelPodName != "" {
		log.Debug("LabelPodName =", b.LabelPodName)
		headList = append(headList, b.LabelPodName)
	}

	log.Debug("headList =", headList)
	return headList
}

// GetDefaultCells: returns an array of cells prepopulated with the common information
func (b *RowBuilder) GetDefaultCells() []Cell {
	log := logger{location: "RowBuilder:GetDefaultCells"}
	log.Debug("Start")

	if b.info.TreeView {
		return []Cell{
			NewCellText(b.info.ContainerType),
			NewCellText(b.info.Namespace),
			NewCellText(b.info.NodeName),
		}
	} else {
		return []Cell{
			NewCellText(b.info.ContainerType),
			NewCellText(b.info.Namespace),
			NewCellText(b.info.NodeName),
			NewCellText(b.info.PodName),
			NewCellText(b.info.ContainerName),
		}
	}
}

func (b *BuilderInformation) BuildTreeCell(cellList []Cell) []Cell {
	var namePrefix string

	log := logger{location: "RowBuilder:BuildTreeCell"}
	log.Debug("Start")

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
