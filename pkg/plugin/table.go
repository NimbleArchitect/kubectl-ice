package plugin

import (
	"fmt"
	"strings"
)

type headerRow struct {
	title        string
	sort         int // 0:no-sort, 1:sort-forward, 2:sort-backward
	columnLength int
	hidden       bool
}

type Table struct {
	currentRow  int
	headCount   int
	columnOrder []int
	rowOrder    []int
	head        []headerRow
	data        [][]string
}

// sets the header row to the specified array of strings
// headerRow is always reinitilized to empty before headers are added
func (t *Table) SetHeader(headItem ...string) {

	t.head = make([]headerRow, len(headItem))

	if len(t.columnOrder) == 0 {
		t.columnOrder = []int{}
	}

	for i := 0; i < len(headItem); i++ {
		tmpHead := headerRow{}
		tmpHead.title = headItem[i]
		tmpHead.columnLength = len(headItem[i]) + 2
		tmpHead.sort = 0

		t.head[i] = tmpHead

		t.columnOrder = append(t.columnOrder, i)
	}

	t.headCount = len(headItem)
}

// Adds a new row to the end of the table, accepts an array of strings
func (t *Table) AddRow(row ...string) {

	for i := 0; i < t.headCount; i++ {
		if len(row[i]) >= t.head[i].columnLength {
			t.head[i].columnLength = len(row[i]) + 2
		}
	}

	t.data = append(t.data, row)
	t.currentRow += 1
}

//  changes the order of columns displayed in the table, specifying a subset of the column
// numbers will place those at the front in the order specified all other columns remain untouched
func (t *Table) Order(items ...int) {
	// rather then reordering all columns we have an order array that we can loop through
	// order contains the actual column number to use next
	orderedList := []int{}

	for i := 0; i < len(t.columnOrder); i++ {
		found := false
		for c := 0; c < len(items); c++ {
			if items[c] == t.columnOrder[i] {
				found = true
			}
		}
		if !found {
			//fmt.Println(t.columnOrder[i])
			orderedList = append(orderedList, t.columnOrder[i])
		}
	}
	orderedList = append(items, orderedList...)

	t.columnOrder = orderedList

}

// select the column number to hide, columns numbers are the unsorted column number
func (t *Table) HideColumn(columnNumber int) {
	if len(t.head) >= columnNumber {
		t.head[columnNumber].hidden = true
	}
}

// prints the table on the terminal, taking the column order and visibiliy into account
func (t *Table) Print() {
	headLine := ""
	// loop through all headers and make a single line properly spaced
	for col := 0; col < t.headCount; col++ {
		// columnOrder contains the actual column number to use next
		idx := t.columnOrder[col]
		if t.head[idx].hidden {
			continue
		}

		word := t.head[idx].title
		if len(word) == 0 {
			word = "-"
		}
		pad := strings.Repeat(" ", t.head[idx].columnLength-len(word))
		headLine += fmt.Sprint(word, pad)
	}
	// print the header in one long line
	fmt.Println(strings.TrimRight(headLine, " "))

	// loop through each row
	for rowNum := 0; rowNum < len(t.data); rowNum++ {
		line := ""
		row := t.data[rowNum]
		// now loop through each column the the currentl selected row
		for col := 0; col < t.headCount; col++ {
			idx := t.columnOrder[col]
			if t.head[idx].hidden {
				continue
			}

			word := row[idx]
			if len(word) == 0 {
				word = "-"
			}
			pad := strings.Repeat(" ", t.head[idx].columnLength-len(word))
			line += fmt.Sprint(word, pad)
		}
		fmt.Println(strings.TrimRight(line, " "))
	}
}

// Prints the table on the terminal as json, all fileds are shown and all are unsorted as
// programs like jq can be used to filter and sort
func (t *Table) PrintJson() {
	// loop through each row
	fmt.Println("{\"data\":[")
	for rowNum := 0; rowNum < len(t.data); rowNum++ {
		line := "{"
		row := t.data[rowNum]
		// now loop through each column the the currentl selected row
		for col := 0; col < t.headCount; col++ {
			word := row[col]
			if len(word) == 0 {
				word = ""
			}
			line += fmt.Sprintf("\"%s\": \"%s\"", t.head[col].title, word)
			// add , to the end of every key/value except the last
			if col+1 < t.headCount {
				line += ", "
			}
		}

		line += "}"
		// again add the , to end of every line except the last
		if rowNum+1 < len(t.data) {
			line += ", "
		}

		fmt.Println(line)
	}
	fmt.Println("]}")

}

func (t *Table) Sort(columnNumber int, ascending bool) {
	t.rowOrder[0] = 0
}
