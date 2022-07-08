package plugin

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// holds a set of columns that are common to every subcommand
type containerInfomation struct {
	table              *Table
	labelNodeName      string
	labelNodeValue     string
	labelPodName       string
	labelPodValue      string
	annotationPodName  string
	annotationPodValue string
	containerName      string
	containerType      string
	namespace          string
	nodeName           string
	podName            string
	treeView           bool
}

func (ci *containerInfomation) AddRow(cellList ...[]Cell) []Cell {
	rowList := ci.GetDefaultCells()

	if ci.labelNodeName != "" {
		rowList = append(rowList, NewCellText(ci.labelNodeValue))
	}

	if ci.labelPodName != "" {
		rowList = append(rowList, NewCellText(ci.labelPodValue))
	}

	for _, c := range cellList {
		rowList = append(rowList, c...)
	}

	return rowList
}

func (ci *containerInfomation) ApplyRow(table *Table, cellList ...[]Cell) {
	rowList := ci.GetDefaultCells()

	if ci.labelNodeName != "" {
		rowList = append(rowList, NewCellText(ci.labelNodeValue))
	}

	if ci.labelPodName != "" {
		rowList = append(rowList, NewCellText(ci.labelPodValue))
	}

	for _, c := range cellList {
		rowList = append(rowList, c...)
	}

	table.AddRow(rowList...)
}

// GetDefaultHead: returns the common headers in order
func (ci *containerInfomation) GetDefaultHead() []string {
	var headList []string

	if ci.treeView {
		//in tree view we only create the namespace and nodename columns, the name colume is created outside of this
		// function so we have full control over its contents
		headList = []string{
			"NAMESPACE", "NODE",
		}
	} else {
		headList = []string{
			"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER",
		}
	}

	if ci.labelNodeName != "" {
		fmt.Println("***")
		headList = append(headList, ci.labelNodeName)
	}

	if ci.labelPodName != "" {
		headList = append(headList, ci.labelPodName)
	}
	return headList
}

// GetDefaultCells: returns an array of cells prepopulated with the common information
func (ci *containerInfomation) GetDefaultCells() []Cell {
	if ci.treeView {
		return []Cell{
			NewCellText(ci.namespace),
			NewCellText(ci.nodeName),
		}
	} else {
		return []Cell{
			NewCellText(ci.containerType),
			NewCellText(ci.namespace),
			NewCellText(ci.nodeName),
			NewCellText(ci.podName),
			NewCellText(ci.containerName),
		}
	}
}

// LoadFromPod: the common information is read from a given pod object and stored internally
func (ci *containerInfomation) LoadFromPod(pod v1.Pod) {
	ci.podName = pod.Name
	ci.namespace = pod.Namespace
	ci.nodeName = pod.Spec.NodeName
}

// SetVisibleColumns: sets the visable columns based on properties from flags
func (ci *containerInfomation) SetVisibleColumns(table Table, flags commonFlags) {
	if ci.treeView {
		//only hide the nodename and namespace, podname is always show in tree view
		if !flags.showNamespaceName {
			table.HideColumn(0)
		}

		if !flags.showNodeName {
			table.HideColumn(1)
		}
		return
	}

	if !flags.showNamespaceName {
		table.HideColumn(1)
	}

	if !flags.showNodeName {
		table.HideColumn(2)
	}

	if !flags.showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(3)
	}

}
