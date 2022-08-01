package plugin

import (
	"reflect"
	"testing"
)

var table Table

//*****************
//SetHeader
//*****************
type tableSetHeaderTest struct {
	arg1        []string
	count       int
	columnorder []int
	expected    []headerRow
}

var tableSetHeaderTests = []tableSetHeaderTest{
	{[]string{"a", "b"}, 2, []int{0, 1}, []headerRow{{matchFilter{"", 0, false, false}, 3, 0, false, 0, "a"}, {matchFilter{"", 0, false, false}, 3, 0, false, 0, "b"}}},
}

func TestTableSetHeader(t *testing.T) {

	for _, test := range tableSetHeaderTests {

		table.SetHeader(test.arg1...)
		if table.headCount != test.count {
			t.Errorf("Output %v not equal to expected \"%v\"", table.headCount, test.count)
		}
		if !reflect.DeepEqual(table.head, test.expected) {
			t.Errorf("Output %v not equal to expected \"%v\"", table.head, test.expected)
		}
		if !reflect.DeepEqual(table.columnOrder, test.columnorder) {
			t.Errorf("Output %v not equal to expected \"%v\"", table.columnOrder, test.columnorder)
		}
	}

}

//*****************
//AddRow
//*****************
type addRowTest struct {
	arg1      []Cell
	rowCount  int
	columnLen int
	expected  [][]Cell
}

var addRowTests = []addRowTest{
	{[]Cell{NewCellText("one")}, 1, 5, [][]Cell{{Cell{"one", 0, 0, 0, 0, 0}}}},
	{[]Cell{NewCellText("two")}, 2, 5, [][]Cell{{Cell{"one", 0, 0, 0, 0, 0}}, {Cell{"two", 0, 0, 0, 0, 0}}}},
	{[]Cell{NewCellText("three")}, 3, 7, [][]Cell{{Cell{"one", 0, 0, 0, 0, 0}}, {Cell{"two", 0, 0, 0, 0, 0}}, {Cell{"three", 0, 0, 0, 0, 0}}}},
	{[]Cell{NewCellText("four"), NewCellText("extra"), NewCellText("larger")}, 4, 7, [][]Cell{{Cell{"one", 0, 0, 0, 0, 0}}, {Cell{"two", 0, 0, 0, 0, 0}}, {Cell{"three", 0, 0, 0, 0, 0}}, {Cell{"four", 0, 0, 0, 0, 0}, Cell{"extra", 0, 0, 0, 0, 0}, Cell{"larger", 0, 0, 0, 0, 0}}}},
}

func TestAddRow(t *testing.T) {
	table.SetHeader("A")

	for _, test := range addRowTests {

		table.AddRow(test.arg1...)
		if table.currentRow != test.rowCount {
			t.Errorf("Output %v not equal to expected \"%v\"", table.currentRow, test.rowCount)
		}
		if table.head[0].columnLength != test.columnLen {
			t.Errorf("Output %v not equal to expected \"%v\"", table.head[0].columnLength, test.columnLen)
		}
		if !reflect.DeepEqual(table.data, test.expected) {
			t.Errorf("Output %v not equal to expected \"%v\"", table.data, test.expected)
		}
	}

}

//*****************
//Order
//*****************
type orderTest struct {
	arg1     []int
	expected []int
}

var orderTests = []orderTest{
	{[]int{0, 1}, []int{0, 1, 2}},
	{[]int{1, 0}, []int{1, 0, 2}},
	{[]int{2, 1}, []int{2, 1, 0}},
	{[]int{2, 0, 1}, []int{2, 0, 1}},
	{[]int{3, 0, 1, 3}, []int{3, 0, 1, 3, 2}},
}

func TestOrder(t *testing.T) {
	table.SetHeader("A", "B", "C")

	for _, test := range orderTests {
		table.Order(test.arg1...)
		if !reflect.DeepEqual(table.columnOrder, test.expected) {
			t.Errorf("Output %v not equal to expected \"%v\"", table.columnOrder, test.expected)
		}
	}

}

//*****************
//HideColumn
//*****************
type hideColumnTest struct {
	arg1     int
	expected []headerRow
}

// []headerRow{{matchFilter{"", 0, false, false}, 3, 0, false, 0, "a"}, {matchFilter{"", 0, false, false}, 3, 0, false, 0, "b"}}}
var hideColumnTests = []hideColumnTest{
	{2, []headerRow{{matchFilter{"", 0, false, false}, 3, 0, false, 0, "A"}, {matchFilter{"", 0, false, false}, 3, 0, false, 0, "B"}, {matchFilter{"", 0, false, false}, 3, 0, true, 0, "C"}}},
	{2, []headerRow{{matchFilter{"", 0, false, false}, 3, 0, false, 0, "A"}, {matchFilter{"", 0, false, false}, 3, 0, false, 0, "B"}, {matchFilter{"", 0, false, false}, 3, 0, true, 0, "C"}}},
	{0, []headerRow{{matchFilter{"", 0, false, false}, 3, 0, true, 0, "A"}, {matchFilter{"", 0, false, false}, 3, 0, false, 0, "B"}, {matchFilter{"", 0, false, false}, 3, 0, true, 0, "C"}}},
}

func TestHideColumn(t *testing.T) {
	table.SetHeader("A", "B", "C")

	for _, test := range hideColumnTests {
		table.HideColumn(test.arg1)
		if !reflect.DeepEqual(table.head, test.expected) {
			t.Errorf("Output %v not equal to expected \"%v\"", table.head, test.expected)
		}
	}

}

func TestHideColumnPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// The following is the code under test
	table.HideColumn(4)
}
