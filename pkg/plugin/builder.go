package plugin

import (
	v1 "k8s.io/api/core/v1"
)

type Looper interface {
	BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error)
	BuildBranch(info BuilderInformation) ([]Cell, error)
	BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error)
	BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error)
	BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error)
	Headers() []string
	HideColumns(info BuilderInformation) []int
	Sum(rows [][]Cell) []Cell
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
	ShowTreeView       bool //show the standard tree view with the resource sets as the root
	ShowPodName        bool
	ShowInitContainers bool
	ShowContainerType  bool
	ShowNodeTree       bool                  //show the tree view with the nodes at the root level rather than just the resource sets at root
	FilterList         map[string]matchValue // used to filter out rows from the table during Print function
	info               BuilderInformation
	DefaultHeaderLen   int

	annotationLabel map[string]map[string]map[string]map[string]string
}

type BuilderInformation struct {
	Pod             *v1.Pod
	ContainerName   string
	ContainerType   string
	ContainerPrefix string
	Namespace       string
	NodeName        string
	Name            string //objects name
	TreeView        bool
	TypeName        string //k8s kind
	// BranchType      int
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

	b.ShowNodeTree = true

	if b.ShowTreeView {
		err := b.populateAnnotationsLabels(podList)
		if err != nil {
			return err
		}
		ol := b.Connection.BuildOwnersList()

		for _, value := range ol {
			var rowid int

			if b.ShowNodeTree {
				rowid = b.Table.AddPlaceHolderRow()
			}

			totals, err := b.walkTreeCreateRow(loop, *value)
			if err != nil {
				return err
			}

			if b.ShowNodeTree {
				b.info.Namespace = value.namespace
				b.info.Name = value.name
				b.info.ContainerType = "N"
				b.info.TypeName = value.kind
				// partOut, _ := loop.BuildBranch(b.info)
				partOut := loop.Sum(totals)
				tblOut := b.makeRow(value.indent, b.info, partOut)
				if len(tblOut) > 0 {
					b.Table.UpdatePlaceHolderRow(rowid, tblOut)
				}
			}
		}

	} else {
		return b.BuildContainerTable(loop, podList)
	}
	return nil
}

func (b *RowBuilder) walkTreeCreateRow(loop Looper, parent node) ([][]Cell, error) {
	var parentTotals [][]Cell

	for _, value := range parent.child {
		var totals [][]Cell
		var tblOut []Cell

		rowid := b.Table.AddPlaceHolderRow()

		b.info.Namespace = value.namespace
		b.info.Name = value.name
		b.info.TypeName = value.kind
		b.info.ContainerType = value.kindIndicator

		if value.kind == "Pod" {
			pod := value.data.pod
			partOut, err := b.buildPodRow(loop, pod, value.indent, value.kind)
			if err != nil {
				return [][]Cell{}, err
			}
			totals = append(totals, partOut...)
		} else {
			//make the row for the table heade line
			partOut, _ := loop.BuildBranch(b.info)
			resourceTotals, err := b.walkTreeCreateRow(loop, *value)
			if err == nil {
				totals = append(totals, resourceTotals...)
			}
			totals = append(totals, partOut)
		}

		//we have to reset these as they are changed during code run
		b.info.Namespace = value.namespace
		b.info.Name = value.name
		b.info.TypeName = value.kind
		b.info.ContainerType = value.kindIndicator

		partOut := loop.Sum(totals)
		tblOut = b.makeRow(value.indent, b.info, partOut)
		if len(tblOut) > 0 {
			b.Table.UpdatePlaceHolderRow(rowid, tblOut)
		}

		b.labelNodeValue = ""
		b.labelPodValue = ""
		b.annotationPodValue = ""
		parentTotals = append(parentTotals, partOut)
	}

	return parentTotals, nil
}

func (b *RowBuilder) buildPodRow(loop Looper, pod v1.Pod, indent int, kind string) ([][]Cell, error) {
	log := logger{location: "RowBuilder:buildPodRow"}
	log.Debug("Start")

	log.Debug("pod.Name =", pod.Name)
	// b.info.TypeName = b.info.kind
	b.info.Pod = &pod
	b.info.Name = pod.Name
	b.info.Namespace = pod.Namespace
	b.info.NodeName = pod.Spec.NodeName

	//check if we have any labels that need to be shown as columns
	b.setValuesAnnotationLabel(pod)

	tblOut, err := b.podLoop(indent+1, loop, pod)
	if err != nil {
		return [][]Cell{}, err
	}

	// TODO: if returning a proper filtered and sorted list from ownerTypeList dosent work then this is a backup process to follow
	// if rowscommited <= 0 {
	// 	b.Table.HidePlaceHolderRow(rowid)
	// }

	// b.info.BranchType = POD
	//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
	b.info.ContainerType = "P"
	b.info.TypeName = kind
	tblBranch, err := loop.BuildBranch(b.info)
	if err != nil {
		return [][]Cell{}, err
	}
	tblOut = append(tblOut, tblBranch)
	return tblOut, nil
}

//check if any labels or annotations are needed and set their values
func (b *RowBuilder) setValuesAnnotationLabel(pod v1.Pod) {
	if b.LabelNodeName != "" {
		b.labelNodeValue = b.annotationLabel["label"]["node"][pod.Spec.NodeName][b.LabelNodeName]
	}
	if b.LabelPodName != "" {
		b.labelPodValue = b.annotationLabel["label"]["pod"][pod.Name][b.LabelPodName]
	}
	if b.AnnotationPodName != "" {
		b.annotationPodValue = b.annotationLabel["annotation"]["pod"][pod.Name][b.AnnotationPodName]
	}

}

func (b *RowBuilder) populateAnnotationsLabels(podList []v1.Pod) error {
	log := logger{location: "RowBuilder:BuildContainerTable"}
	log.Debug("Start")
	//                          type       kind       pod        label  value
	b.annotationLabel = make(map[string]map[string]map[string]map[string]string)
	b.annotationLabel["label"] = make(map[string]map[string]map[string]string)
	b.annotationLabel["annotation"] = make(map[string]map[string]map[string]string)

	if b.LabelNodeName != "" {
		log.Debug("b.LabelNodeName", b.LabelNodeName)
		nodeLabels, err := b.Connection.GetNodeLabels(podList)
		if err != nil {
			return err
		}
		b.annotationLabel["label"]["node"] = nodeLabels
	}

	if b.LabelPodName != "" {
		log.Debug("b.LabelPodName", b.LabelPodName)
		podLabels, err := b.Connection.GetPodLabels(podList)
		if err != nil {
			return err
		}
		b.annotationLabel["label"]["pod"] = podLabels
	}

	if b.AnnotationPodName != "" {
		log.Debug("b.AnnotationPodName", b.AnnotationPodName)
		podAnnotations, err := b.Connection.GetPodAnnotations(podList)
		if err != nil {
			return err
		}
		b.annotationLabel["annotation"]["pod"] = podAnnotations
	}

	return nil
}

// Build normal table
func (b *RowBuilder) BuildContainerTable(loop Looper, podList []v1.Pod) error {
	log := logger{location: "RowBuilder:BuildContainerTable"}
	log.Debug("Start")

	err := b.populateAnnotationsLabels(podList)
	if err != nil {
		return err
	}

	for _, pod := range podList {
		log.Debug("pod.Name =", pod.Name)
		b.info.Pod = &pod
		b.info.Name = pod.Name
		b.info.Namespace = pod.Namespace
		b.info.NodeName = pod.Spec.NodeName
		b.info.ContainerType = "P"
		b.info.TypeName = "Pod"
		// b.info.BranchType = FLAT

		//check if we have any labels that need to be shown as columns
		b.setValuesAnnotationLabel(pod)

		_, err := b.podLoop(0, loop, pod)
		if err != nil {
			return err
		}

	}

	if err := b.Table.SortByNames(b.CommonFlags.sortList...); err != nil {
		return err
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

// PodLoop given a pod we loop over all containers adding to the table as we go
//  returns a copy of rows added and nil on success
func (b *RowBuilder) podLoop(indentLevel int, loop Looper, pod v1.Pod) ([][]Cell, error) {
	var total int
	var podRowsOut [][]Cell

	log := logger{location: "RowBuilder:PodLoop"}
	log.Debug("Start")

	// b.info.BranchType = CONTAINER

	if b.ShowInitContainers {
		log.Debug("loop init ContainerStatuses")
		b.info.ContainerType = "I"
		b.info.TypeName = "InitContainer"
		if b.LoopStatus {
			for _, container := range pod.Status.InitContainerStatuses {
				// should the container be processed
				log.Debug("processing -", container.Name)
				if skipContainerName(b.CommonFlags, container.Name) {
					continue
				}
				b.info.ContainerName = container.Name
				b.info.Name = container.Name
				allRows, err := loop.BuildContainerStatus(container, b.info)
				if err != nil {
					return [][]Cell{}, err
				}
				for _, row := range allRows {
					rowsOut := b.makeRow(indentLevel, b.info, row)
					total += len(rowsOut)
					// b.printHeadIfNeeded()
					b.Table.AddRow(rowsOut...)
				}
				podRowsOut = append(podRowsOut, allRows...)
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
				b.info.Name = container.Name
				allRows, err := loop.BuildContainerSpec(container, b.info)
				if err != nil {
					return [][]Cell{}, err
				}
				for _, row := range allRows {
					rowsOut := b.makeRow(indentLevel, b.info, row)
					total += len(rowsOut)
					// b.printHeadIfNeeded()
					b.Table.AddRow(rowsOut...)
				}
				podRowsOut = append(podRowsOut, allRows...)
			}
		}
	}

	//now show the container line
	log.Debug("loop standard ContainerStatuses")
	b.info.ContainerType = "C"
	b.info.TypeName = "Container"
	if b.LoopStatus {
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(b.CommonFlags, container.Name) {
				continue
			}
			log.Debug("processing -", container.Name)
			b.info.ContainerName = container.Name
			b.info.Name = container.Name
			allRows, err := loop.BuildContainerStatus(container, b.info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(indentLevel, b.info, row)
				total += len(rowsOut)
				// b.printHeadIfNeeded()
				b.Table.AddRow(rowsOut...)
			}
			podRowsOut = append(podRowsOut, allRows...)
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
			b.info.Name = container.Name
			allRows, err := loop.BuildContainerSpec(container, b.info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(indentLevel, b.info, row)
				total += len(rowsOut)
				// b.printHeadIfNeeded()
				b.Table.AddRow(rowsOut...)
			}
			podRowsOut = append(podRowsOut, allRows...)
		}
	}

	log.Debug("loop ephemeral ContainerStatuses")
	b.info.ContainerType = "E"
	b.info.TypeName = "EphemeralContainer"
	if b.LoopStatus {
		for _, container := range pod.Status.EphemeralContainerStatuses {
			// should the container be processed
			if skipContainerName(b.CommonFlags, container.Name) {
				continue
			}
			log.Debug("processing -", container.Name)
			b.info.ContainerName = container.Name
			b.info.Name = container.Name
			allRows, err := loop.BuildContainerStatus(container, b.info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(indentLevel, b.info, row)
				total += len(rowsOut)
				// b.printHeadIfNeeded()
				b.Table.AddRow(rowsOut...)
			}
			podRowsOut = append(podRowsOut, allRows...)
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
			b.info.Name = container.Name
			allRows, err := loop.BuildEphemeralContainerSpec(container, b.info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeRow(indentLevel, b.info, row)
				total += len(rowsOut)
				// b.printHeadIfNeeded()
				b.Table.AddRow(rowsOut...)
			}
			podRowsOut = append(podRowsOut, allRows...)
		}
	}
	// }

	return podRowsOut, nil
}

// MakeRow adds the listed columns to the default columns, outputs
//  the complete row as a list of columns
func (b *RowBuilder) makeRow(indentLevel int, info BuilderInformation, columns ...[]Cell) []Cell {
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

	if b.info.TreeView {
		name := ""
		//default cells dont have name column, need to add it in tree view
		if len(b.info.TypeName) == 0 {
			name = b.info.Name
		} else {
			name = b.info.TypeName + "/" + b.info.Name
		}
		if !b.ShowNodeTree {
			rowList = append(rowList, NewCellText(indentText(indentLevel-1, name)))
		} else {
			rowList = append(rowList, NewCellText(indentText(indentLevel, name)))

		}

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
			NewCellText(b.info.Name),
			NewCellText(b.info.ContainerName),
		}
	}
}
