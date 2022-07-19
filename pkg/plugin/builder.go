package plugin

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

const (
	TREE      = iota
	CONTAINER = iota
	PARENT    = iota
)

type Looper interface {
	BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error)
	BuildBranch(info BuilderInformation, podList []v1.Pod) ([][]Cell, error)
	BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error)
	BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error)
	BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error)
	Headers() []string
	HideColumns(info BuilderInformation) []int
}

type RowBuilder struct {
	Connection         *Connector
	Table              *Table
	CommonFlags        commonFlags
	PodName            []string //list of pod names to retrieve
	LoopStatus         bool     //do we need to loop over v1.Pod.Status.ContainerStatus
	LoopSpec           bool     //should we loop over v1.Pod.Spec.Containers
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
	FilterList         map[string]matchValue // used to filter out rows from the table during Print function
	info               BuilderInformation
	DefaultHeaderLen   int

	hTreeViewRow []Cell
}

type BuilderInformation struct {
	Pod           *v1.Pod
	ContainerName string
	ContainerType string
	Namespace     string
	NodeName      string
	PodName       string
	TreeView      bool
	TypeName      string
	BranchType    int
	Owner         string
	OwnerType     string
}

// SetFlagsFrom sets the common flags to match the values retrieved from the passed object
func (b *RowBuilder) SetFlagsFrom(commonFlagList commonFlags) {

	log := logger{location: "RowBuilder:SetFlagsFrom"}
	log.Debug("Start")

	b.CommonFlags = commonFlagList

	b.ShowTreeView = commonFlagList.showTreeView
	b.LabelNodeName = commonFlagList.labelNodeName
	b.LabelPodName = commonFlagList.labelPodName
	b.AnnotationPodName = commonFlagList.annotationPodName
	b.FilterList = b.CommonFlags.filterList

	// we always show the pod name by default
	b.ShowPodName = true

	// if a single pod is selected we dont need to show its name
	if len(b.PodName) == 1 {
		if len(b.PodName[0]) >= 1 {
			log.Debug("builder.ShowPodName = false")
			b.ShowPodName = false
		}
	}

	if b.ShowTreeView {
		log.Debug("b.info.TreeView = true")
		b.info.TreeView = true
	}

	if !b.ShowContainerType {
		b.ShowContainerType = b.CommonFlags.showContainerType
	}

}

// Build
func (b *RowBuilder) Build(loop Looper) error {
	log := logger{location: "RowBuilder:Build"}
	log.Debug("Start")

	err := b.LoadHeaders(loop)
	if err != nil {
		return err
	}

	podList, err := b.Connection.GetPods(b.PodName)
	if err != nil {
		return err
	}

	if b.ShowTreeView {
		return b.BuildContainerTree(loop, podList)
	} else {
		return b.BuildContainerTable(loop, podList)
	}
}

// BuildRows
func (b *RowBuilder) BuildContainerTree(loop Looper, podList []v1.Pod) error {
	var err error
	var nodeLabels map[string]map[string]string
	var podLabels map[string]map[string]string
	var podAnnotations map[string]map[string]string

	log := logger{location: "RowBuilder:BuildContainerTree"}
	log.Debug("Start")

	if b.LabelNodeName != "" {
		log.Debug("b.LabelNodeName", b.LabelNodeName)
		nodeLabels, err = b.Connection.GetNodeLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.LabelPodName != "" {
		log.Debug("b.LabelPodName", b.LabelPodName)
		podLabels, err = b.Connection.GetPodLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.AnnotationPodName != "" {
		log.Debug("b.AnnotationPodName", b.AnnotationPodName)
		podAnnotations, err = b.Connection.GetPodAnnotations(podList)
		if err != nil {
			return err
		}
	}

	ownerList, ownerTypeList := b.Connection.GetOwnersList()

	for owner, podList := range ownerList {
		// ########################################
		// somehow need to count columns and need a list of ownerNames and ownerTypes from ownerList
		b.info.ContainerType = "R"

		b.info.BranchType = PARENT
		b.info.Owner = owner
		b.info.OwnerType = ownerTypeList[owner]

		allRows, err := loop.BuildBranch(b.info, podList)
		if err != nil {
			return err
		}
		// tblOut := []Cell{NewCellText(owner)}
		parentType := make([]Cell, 1)
		parentType[0] = NewCellText(fmt.Sprint(b.info.OwnerType, "/", b.info.Owner))
		for _, row := range allRows {
			tblOut := append(parentType, row...)
			// true/false in this function signals we are on a container or pod line for tree view, need to create a replicaset line now though
			//  maybe an int would be better???
			b.Table.AddRow(b.makeRow(TREE, 0, b.info, tblOut)...)

			// rowsOut := b.makeRow(CONTAINER, b.info, row)
			// total += len(rowsOut)
			// b.printHeadIfNeeded()
			// b.Table.AddRow(rowsOut...)
		}
		// ########################################

		b.info.TypeName = "Pod"
		for _, pod := range podList {

			log.Debug("pod.Name =", pod.Name)
			b.info.Pod = &pod
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
			if b.ShowTreeView {
				b.info.ContainerType = "P"
				tblOut, err := loop.BuildPod(pod, b.info)
				if err != nil {
					return err
				}

				//this is a tree view, so we have a name column to deal with 'Pod/foo-6f67dcc579-znb55'
				parentType := make([]Cell, 1)
				// parentType[0] = NewCellText(fmt.Sprint(b.info.TypeName, "/", b.info.PodName))
				parentType[0] = NewCellText(indentText(1, b.info.TypeName, "/", b.info.PodName))
				tblOut = append(parentType, tblOut...)

				log.Debug("rowsOut =", b.hTreeViewRow)
				// we save the row rather than add it to the table so we can control the output later on
				b.hTreeViewRow = b.makeRow(TREE, 1, b.info, tblOut)
			}

			_, err = b.podLoop(2, loop, pod)
			if err != nil {
				return err
			}
		}
	}

	// sorting by column breaks the tree view so we sliently skip
	if !b.ShowTreeView {
		if err := b.Table.SortByNames(b.CommonFlags.sortList...); err != nil {
			return err
		}
	}

	return nil
}

// Build normal table
func (b *RowBuilder) BuildContainerTable(loop Looper, podList []v1.Pod) error {
	var err error
	var nodeLabels map[string]map[string]string
	var podLabels map[string]map[string]string
	var podAnnotations map[string]map[string]string

	log := logger{location: "RowBuilder:BuildContainerTable"}
	log.Debug("Start")

	if b.LabelNodeName != "" {
		log.Debug("b.LabelNodeName", b.LabelNodeName)
		nodeLabels, err = b.Connection.GetNodeLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.LabelPodName != "" {
		log.Debug("b.LabelPodName", b.LabelPodName)
		podLabels, err = b.Connection.GetPodLabels(podList)
		if err != nil {
			return err
		}
	}

	if b.AnnotationPodName != "" {
		log.Debug("b.AnnotationPodName", b.AnnotationPodName)
		podAnnotations, err = b.Connection.GetPodAnnotations(podList)
		if err != nil {
			return err
		}
	}

	b.info.TypeName = "Pod"
	for _, pod := range podList {
		log.Debug("pod.Name =", pod.Name)
		b.info.Pod = &pod
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
		if b.ShowTreeView {
			b.info.ContainerType = "P"
			tblOut, err := loop.BuildPod(pod, b.info)
			if err != nil {
				return err
			}

			//this is a tree view, so we have a name column to deal with 'Pod/foo-6f67dcc579-znb55'
			parentType := make([]Cell, 1)
			parentType[0] = NewCellText(fmt.Sprint(b.info.TypeName, "/", b.info.PodName))
			tblOut = append(parentType, tblOut...)

			log.Debug("rowsOut =", b.hTreeViewRow)
			// we save the row rather than add it to the table so we can control the output later on
			b.hTreeViewRow = b.makeRow(TREE, 0, b.info, tblOut)
		}

		_, err = b.podLoop(0, loop, pod)
		if err != nil {
			return err
		}

	}

	// sorting by column breaks the tree view so we sliently skip
	if !b.ShowTreeView {
		if err := b.Table.SortByNames(b.CommonFlags.sortList...); err != nil {
			return err
		}
	}

	return nil
}

// LoadHeaders sets the default column headers hiding as needed
func (b *RowBuilder) LoadHeaders(loop Looper) error {
	var tblHead []string
	var hideColumns []int

	log := logger{location: "RowBuilder:LoadHeaders"}
	log.Debug("Start")

	tblHead = b.getDefaultHead()

	// save the default lengh now as we need to use it in other functions
	defaultHeaderLen := len(tblHead)
	log.Debug("len(defaultHeaderLen) =", defaultHeaderLen)
	b.DefaultHeaderLen = defaultHeaderLen

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

	b.setVisibleColumns()

	log.Debug("len(hideColumns) =", len(hideColumns))
	for _, id := range hideColumns {
		b.Table.HideColumn(defaultHeaderLen + id)
	}

	return nil
}

// SetVisibleColumns hides default columns based on various flags
func (b *RowBuilder) setVisibleColumns() {
	log := logger{location: "RowBuilder:SetVisibleColumns"}
	log.Debug("Start")

	if !b.ShowContainerType {
		b.Table.HideColumn(0)
	}

	if b.info.TreeView {
		//only hide the nodename in tree view
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

//printHeadIfNeeded prints the tree view pod line once printed b.hTreeViewRow is cleared
//  as we only want to print it once per pod
func (b *RowBuilder) printHeadIfNeeded() {
	if len(b.hTreeViewRow) > 0 {
		b.Table.AddRow(b.hTreeViewRow...)
		b.hTreeViewRow = []Cell{}
	}
}

// PodLoop given a pod we loop over all containers adding to the table as we go
//  returns total count of rows added and nil on success
func (b *RowBuilder) podLoop(indentLevel int, loop Looper, pod v1.Pod) (int, error) {
	var total int

	log := logger{location: "RowBuilder:PodLoop"}
	log.Debug("Start")

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
					return 0, err
				}
				for _, row := range allRows {
					rowsOut := b.makeRow(CONTAINER, indentLevel, b.info, row)
					total += len(rowsOut)
					b.printHeadIfNeeded()
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
					return 0, err
				}
				for _, row := range allRows {
					rowsOut := b.makeRow(CONTAINER, indentLevel, b.info, row)
					total += len(rowsOut)
					b.printHeadIfNeeded()
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
				return 0, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(CONTAINER, indentLevel, b.info, row)
				total += len(rowsOut)
				b.printHeadIfNeeded()
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
				return 0, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(CONTAINER, indentLevel, b.info, row)
				total += len(rowsOut)
				b.printHeadIfNeeded()
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
				return 0, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(CONTAINER, indentLevel, b.info, row)
				total += len(rowsOut)
				b.printHeadIfNeeded()
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
				return 0, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(CONTAINER, indentLevel, b.info, row)
				total += len(rowsOut)
				b.printHeadIfNeeded()
				b.Table.AddRow(rowsOut...)
			}
		}
	}
	// }

	return total, nil
}

// MakeRow adds the listed columns to the default columns, outputs
//  the complete row as a list of columns
func (b *RowBuilder) makeRow(containerLine int, indentLevel int, info BuilderInformation, columns ...[]Cell) []Cell {
	log := logger{location: "RowBuilder:MakeRow"}
	log.Debug("Start")

	rowList := b.getDefaultCells()

	if b.LabelNodeName != "" {
		rowList = append(rowList, NewCellText(b.labelNodeValue))
	}

	if b.LabelPodName != "" {
		rowList = append(rowList, NewCellText(b.labelPodValue))
	}

	if b.AnnotationPodName != "" {
		rowList = append(rowList, NewCellText(b.annotationPodValue))
	}

	if containerLine == CONTAINER && b.info.TreeView {
		rowList = b.info.BuildTreeCell(indentLevel, rowList)
	}

	for _, c := range columns {
		rowList = append(rowList, c...)
	}

	log.Debug("len(rowList) =", len(rowList))
	return rowList
}

// GetDefaultHead: returns the common headers in order
func (b *RowBuilder) getDefaultHead() []string {
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

	if b.AnnotationPodName != "" {
		log.Debug("AnnotationPodName =", b.AnnotationPodName)
		headList = append(headList, b.AnnotationPodName)
	}

	if b.info.TreeView {
		headList = append(headList, "NAME")
	}

	log.Debug("headList =", headList)
	return headList
}

// GetDefaultCells: returns an array of cells prepopulated with the common information
func (b *RowBuilder) getDefaultCells() []Cell {
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

func (b *BuilderInformation) BuildTreeCell(indentLevel int, cellList []Cell) []Cell {
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
		// NewCellText(fmt.Sprint("  └─", namePrefix, b.ContainerName)),
		NewCellText(indentText(indentLevel, namePrefix, b.ContainerName)),
	)
	return cellList
}
