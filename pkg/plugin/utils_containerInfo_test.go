package plugin

// var ciinfo containerInfomation

//*****************
//GetDefaultHead
//*****************
// type getDefaultHeadTest struct {
// 	treeview      bool
// 	labelNodeName string
// 	labelPodName  string
// 	expected      []string
// }

// var getDefaultHead = []getDefaultHeadTest{
// 	{false, "", "", []string{"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER"}},
// 	{true, "", "", []string{"NAMESPACE", "NODE"}},
// 	{false, "nodename", "", []string{"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER", "nodename"}},
// 	{false, "", "podname", []string{"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER", "podname"}},
// 	{false, "nodename", "podname", []string{"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER", "nodename", "podname"}},
// 	{true, "nodename", "", []string{"NAMESPACE", "NODE", "nodename"}},
// 	{true, "", "podname", []string{"NAMESPACE", "NODE", "podname"}},
// 	{true, "nodename", "podname", []string{"NAMESPACE", "NODE", "nodename", "podname"}},
// 	{false, "", "", []string{"T", "NAMESPACE", "NODE", "PODNAME", "CONTAINER"}},
// }

// func TestGetDefaultHead(t *testing.T) {

// 	for _, test := range getDefaultHead {
// 		ciinfo.treeView = test.treeview
// 		ciinfo.labelNodeName = test.labelNodeName
// 		ciinfo.labelPodName = test.labelPodName

// 		if output := ciinfo.GetDefaultHead(); !reflect.DeepEqual(output, test.expected) {
// 			t.Errorf("Output %v not equal to expected \"%v\"", output, test.expected)
// 		}
// 	}

// }

//*****************
//GetDefaultCells
//*****************
// type getDefaultCellsTest struct {
// 	treeview bool
// 	expected []Cell
// }

// var getDefaultCells = []getDefaultCellsTest{
// 	{false, []Cell{{"I", 0, 0, 0}, {"ice-ns", 0, 0, 0}, {"main-node", 0, 0, 0}, {"trialpod", 0, 0, 0}, {"mainapp", 0, 0, 0}}},
// 	{true, []Cell{{"ice-ns", 0, 0, 0}, {"main-node", 0, 0, 0}}},
// }

// func TestGetDefaultCells(t *testing.T) {

// 	for _, test := range getDefaultCells {
// 		ciinfo.namespace = "ice-ns"
// 		ciinfo.nodeName = "main-node"
// 		ciinfo.containerType = "I"
// 		ciinfo.podName = "trialpod"
// 		ciinfo.containerName = "mainapp"

// 		ciinfo.treeView = test.treeview

// 		if output := ciinfo.GetDefaultCells(); !reflect.DeepEqual(output, test.expected) {
// 			t.Errorf("Output %v not equal to expected \"%v\"", output, test.expected)
// 		}
// 	}
// }
