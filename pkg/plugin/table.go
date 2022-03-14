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
	currentRow int
	headCount  int
	head       []headerRow
	data       [][]string
}

// sets the header row to the specified array of strings
// headerRow is always reinitilized to empty before headers are added
func (t *Table) SetHeader(headItem ...string) {

	t.head = make([]headerRow, len(headItem))

	for i := 0; i < len(headItem); i++ {
		tmpHead := headerRow{}
		tmpHead.title = headItem[i]
		tmpHead.columnLength = len(headItem[i]) + 2
		tmpHead.sort = 0

		t.head[i] = tmpHead
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
	for idx := 0; idx < t.headCount; idx++ {
		// columnOrder contains the actual column number to use next
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
		for idx := 0; idx < t.headCount; idx++ {
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
