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
