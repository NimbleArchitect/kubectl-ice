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
	DefaultHeaderLen   int

	annotationLabel map[string]map[string]map[string]map[string]string
}

type BuilderInformation struct {
	Pod     *v1.Pod
	PodName string
	// ContainerName string // container name
	ContainerType string // single letter type id
	Namespace     string
	NodeName      string
	Name          string // objects name
	TreeView      bool
	TypeName      string // k8s kind
}

// SetFlagsFrom sets the common flags to match the values retrieved from the passed object
func (b *RowBuilder) SetFlagsFrom(commonFlagList commonFlags) {

	log := logger{location: "RowBuilder:SetFlagsFrom"}
	log.Debug("Start")

	b.CommonFlags = commonFlagList

	b.ShowTreeView = commonFlagList.showTreeView
	b.ShowNodeTree = commonFlagList.showNodeTree
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

	if !b.ShowContainerType {
		b.ShowContainerType = b.CommonFlags.showContainerType
	}

}

// Build
func (b *RowBuilder) Build(loop Looper) error {

	log := logger{location: "RowBuilder:Build"}
	log.Debug("Start")

	info := BuilderInformation{TreeView: b.ShowTreeView}

	err := b.LoadHeaders(loop, &info)
	if err != nil {
		return err
	}

	podList, err := b.Connection.GetPods(b.PodName)
	if err != nil {
		return err
	}

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

			totals, err := b.walkTreeCreateRow(loop, &info, *value)
			if err != nil {
				return err
			}

			if b.ShowNodeTree {
				info.Namespace = value.namespace
				info.Name = value.name
				info.ContainerType = "N"
				info.TypeName = value.kind
				log.Debug("call loop.Sum for", info.PodName, info.Name)
				partOut := loop.Sum(totals)
				tblOut := b.makeFullRow(&info, value.indent, partOut)
				if len(tblOut) > 0 {
					b.Table.UpdatePlaceHolderRow(rowid, tblOut)
				}
			}
		}

	} else {
		return b.BuildContainerTable(loop, &info, podList)
	}
	return nil
}

// walkTreeCreateRow - recursive function to loop over each child item along with all sub children, buildPodTree
//  is called on each child with the results passed to Sum so we can calculate parent values from the children
func (b *RowBuilder) walkTreeCreateRow(loop Looper, info *BuilderInformation, parent node) ([][]Cell, error) {
	var parentTotals [][]Cell

	log := logger{location: "RowBuilder:walkTreeCreateRow"}
	log.Debug("Start")

	for _, value := range parent.child {
		var totals [][]Cell
		var tblOut []Cell

		rowid := b.Table.AddPlaceHolderRow()
		info.Namespace = value.namespace
		info.Name = value.name
		info.TypeName = value.kind
		info.ContainerType = value.kindIndicator

		if value.kind == "Pod" {
			infoPod := *info
			partOut, err := b.buildPodTree(loop, &infoPod, value.data.pod, value.indent, value.kind)
			if err != nil {
				return [][]Cell{}, err
			}
			totals = append(totals, partOut...)
		} else {
			//make the row for the table heade line
			infoSet := *info
			partOut, _ := loop.BuildBranch(infoSet)
			resourceTotals, err := b.walkTreeCreateRow(loop, &infoSet, *value)
			if err == nil {
				totals = append(totals, resourceTotals...)
			}
			totals = append(totals, partOut)
		}

		log.Debug("call loop.Sum for", info.PodName, info.Name)
		partOut := loop.Sum(totals)
		tblOut = b.makeFullRow(info, value.indent, partOut)
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

// buildPodTree - sets info properties ready to call podLoop and then buildBranch
func (b *RowBuilder) buildPodTree(loop Looper, info *BuilderInformation, pod v1.Pod, indent int, kind string) ([][]Cell, error) {
	log := logger{location: "RowBuilder:buildPodTree"}
	log.Debug("Start")

	log.Debug("pod.Name =", pod.Name)
	info.Pod = &pod
	info.PodName = pod.Name
	info.Namespace = pod.Namespace
	info.NodeName = pod.Spec.NodeName
	info.ContainerType = "P"
	info.TypeName = kind

	//check if we have any labels that need to be shown as columns
	b.setValuesAnnotationLabel(pod)
	infoPod := *info
	tblOut, err := b.podLoop(loop, infoPod, pod, indent+1)
	if err != nil {
		return [][]Cell{}, err
	}

	//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
	tblBranch, err := loop.BuildBranch(*info)
	if err != nil {
		return [][]Cell{}, err
	}
	tblOut = append(tblOut, tblBranch)
	return tblOut, nil
}

// check if any labels or annotations are needed and set their values
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
func (b *RowBuilder) BuildContainerTable(loop Looper, info *BuilderInformation, podList []v1.Pod) error {
	log := logger{location: "RowBuilder:BuildContainerTable"}
	log.Debug("Start")

	err := b.populateAnnotationsLabels(podList)
	if err != nil {
		return err
	}

	for _, pod := range podList {
		log.Debug("pod.Name =", pod.Name)
		infoPod := *info
		infoPod.Pod = &pod
		infoPod.PodName = pod.Name
		infoPod.Namespace = pod.Namespace
		infoPod.NodeName = pod.Spec.NodeName
		infoPod.ContainerType = "P"
		infoPod.TypeName = "Pod"

		//check if we have any labels that need to be shown as columns
		b.setValuesAnnotationLabel(pod)

		_, err := b.podLoop(loop, infoPod, pod, 0)
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
func (b *RowBuilder) LoadHeaders(loop Looper, info *BuilderInformation) error {
	var tblHead []string
	var hideColumns []int

	log := logger{location: "RowBuilder:LoadHeaders"}
	log.Debug("Start")

	tblHead = b.getDefaultHead(info)

	// save the default lengh now as we need to use it in other functions
	defaultHeaderLen := len(tblHead)
	log.Debug("len(defaultHeaderLen) =", defaultHeaderLen)
	b.DefaultHeaderLen = defaultHeaderLen

	hideColumns = loop.HideColumns(*info)

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

	b.setVisibleColumns(info)

	log.Debug("len(hideColumns) =", len(hideColumns))
	for _, id := range hideColumns {
		b.Table.HideColumn(defaultHeaderLen + id)
	}

	return nil
}

// SetVisibleColumns hides default columns based on various flags
func (b *RowBuilder) setVisibleColumns(info *BuilderInformation) {
	log := logger{location: "RowBuilder:SetVisibleColumns"}
	log.Debug("Start")

	if !b.ShowContainerType {
		b.Table.HideColumn(0)
	}

	if info.TreeView {
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
func (b *RowBuilder) podLoop(loop Looper, info BuilderInformation, pod v1.Pod, indentLevel int) ([][]Cell, error) {
	var total int
	var podRowsOut [][]Cell

	log := logger{location: "RowBuilder:PodLoop"}
	log.Debug("Start")

	if b.ShowInitContainers {
		log.Debug("loop init ContainerStatuses")
		info.ContainerType = "I"
		info.TypeName = "InitContainer"
		if b.LoopStatus {
			for _, container := range pod.Status.InitContainerStatuses {
				// should the container be processed
				log.Debug("processing -", container.Name)
				if skipContainerName(b.CommonFlags, container.Name) {
					continue
				}
				// info.ContainerName = container.Name
				info.Name = container.Name
				allRows, err := loop.BuildContainerStatus(container, info)
				if err != nil {
					return [][]Cell{}, err
				}
				for _, row := range allRows {
					rowsOut := b.makeFullRow(&info, indentLevel, row)
					total += len(rowsOut)
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
				// info.ContainerName = container.Name
				info.Name = container.Name
				allRows, err := loop.BuildContainerSpec(container, info)
				if err != nil {
					return [][]Cell{}, err
				}
				for _, row := range allRows {
					rowsOut := b.makeFullRow(&info, indentLevel, row)
					total += len(rowsOut)
					b.Table.AddRow(rowsOut...)
				}
				podRowsOut = append(podRowsOut, allRows...)
			}
		}
	}

	//now show the container line
	log.Debug("loop standard ContainerStatuses")
	info.ContainerType = "C"
	info.TypeName = "Container"
	if b.LoopStatus {
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(b.CommonFlags, container.Name) {
				continue
			}
			log.Debug("processing -", container.Name)
			info.Name = container.Name
			allRows, err := loop.BuildContainerStatus(container, info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeFullRow(&info, indentLevel, row)
				total += len(rowsOut)
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
			info.Name = container.Name
			allRows, err := loop.BuildContainerSpec(container, info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeFullRow(&info, indentLevel, row)
				total += len(rowsOut)
				b.Table.AddRow(rowsOut...)
			}
			podRowsOut = append(podRowsOut, allRows...)
		}
	}

	log.Debug("loop ephemeral ContainerStatuses")
	info.ContainerType = "E"
	info.TypeName = "EphemeralContainer"
	if b.LoopStatus {
		for _, container := range pod.Status.EphemeralContainerStatuses {
			// should the container be processed
			if skipContainerName(b.CommonFlags, container.Name) {
				continue
			}
			log.Debug("processing -", container.Name)
			// info.ContainerName = container.Name
			info.Name = container.Name
			allRows, err := loop.BuildContainerStatus(container, info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeFullRow(&info, indentLevel, row)
				total += len(rowsOut)
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
			// info.ContainerName = container.Name
			info.Name = container.Name
			allRows, err := loop.BuildEphemeralContainerSpec(container, info)
			if err != nil {
				return [][]Cell{}, err
			}
			for _, row := range allRows {
				rowsOut := b.makeFullRow(&info, indentLevel, row)
				total += len(rowsOut)
				b.Table.AddRow(rowsOut...)
			}
			podRowsOut = append(podRowsOut, allRows...)
		}
	}

	return podRowsOut, nil
}

// makeFullRow adds the listed columns to the default columns, outputs
//  the complete row as a list of columns
func (b *RowBuilder) makeFullRow(info *BuilderInformation, indentLevel int, columns ...[]Cell) []Cell {
	log := logger{location: "RowBuilder:makeFullRow"}
	log.Debug("Start")

	rowList := b.getDefaultCells(info)

	if b.LabelNodeName != "" {
		rowList = append(rowList, NewCellText(b.labelNodeValue))
	}

	if b.LabelPodName != "" {
		rowList = append(rowList, NewCellText(b.labelPodValue))
	}

	if b.AnnotationPodName != "" {
		rowList = append(rowList, NewCellText(b.annotationPodValue))
	}

	if info.TreeView {
		name := ""
		//default cells dont have name column, need to add it in tree view
		if len(info.TypeName) == 0 {
			name = info.Name
		} else {
			name = info.TypeName + "/" + info.Name
		}
		if !b.ShowNodeTree {
			// rowList = append(rowList, NewCellText(indentText(indentLevel-1, name)))
			rowList = append(rowList, NewCellTextIndent(name, indentLevel-1))
		} else {
			// rowList = append(rowList, NewCellText(indentText(indentLevel, name)))
			rowList = append(rowList, NewCellTextIndent(name, indentLevel))
		}

	}

	for _, c := range columns {
		rowList = append(rowList, c...)
	}

	log.Debug("len(rowList) =", len(rowList))
	return rowList
}

// GetDefaultHead: returns the common headers in order
func (b *RowBuilder) getDefaultHead(info *BuilderInformation) []string {
	log := logger{location: "RowBuilder:GetDefaultHead"}
	log.Debug("Start")

	var headList []string

	log.Debug("b.info.TreeView =", info.TreeView)
	if info.TreeView {
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

	if info.TreeView {
		headList = append(headList, "NAME")
	}

	log.Debug("headList =", headList)
	return headList
}

// GetDefaultCells: returns an array of cells prepopulated with the common information
func (b *RowBuilder) getDefaultCells(info *BuilderInformation) []Cell {
	log := logger{location: "RowBuilder:GetDefaultCells"}
	log.Debug("Start")

	if info.TreeView {
		return []Cell{
			NewCellText(info.ContainerType),
			NewCellText(info.Namespace),
			NewCellText(info.NodeName),
		}
	} else {
		return []Cell{
			NewCellText(info.ContainerType),
			NewCellText(info.Namespace),
			NewCellText(info.NodeName),
			NewCellText(info.PodName),
			NewCellText(info.Name),
		}
	}
}
