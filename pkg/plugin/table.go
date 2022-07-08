package plugin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//sets the maximum number of spaces allowed in a column, spaces are clipped to this number
const maxLineLength = 80

type matchFilter struct {
	value      string
	comparison int  // 1:>, 2:<, 3:!
	compareEql bool // true:==, true:<=, true:>=
	set        bool
}

type matchValue struct {
	operator string
	value    string
}

type headerRow struct {
	filter       matchFilter
	columnLength int
	columnType   int // 0:string, 1:int
	hidden       bool
	sort         int // 0:no-sort, 1:sort-forward, 2:sort-backward
	title        string
}

type Cell struct {
	text   string
	number int64
	float  float64
	typ    int // 0=string, 1=int64, 2=float64
}

type Table struct {
	currentRow  int
	headCount   int
	columnOrder []int
	rowOrder    []int
	head        []headerRow
	data        [][]Cell
	hideRow     []bool
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
func (t *Table) AddRow(row ...Cell) {
	for i := 0; i < t.headCount; i++ {
		strLen := len([]rune(row[i].text))
		if strLen >= t.head[i].columnLength {
			if (strLen + 2) > maxLineLength {
				t.head[i].columnLength = maxLineLength
			} else {
				t.head[i].columnLength = strLen + 2
			}
		}

		if row[i].typ == 2 {
			t.head[i].columnType = 2
		} else if row[i].typ == 1 {
			t.head[i].columnType = 1
		}
	}

	t.data = append(t.data, row)                  // add data to row
	t.rowOrder = append(t.rowOrder, t.currentRow) // add row number to end of sort list
	t.hideRow = append(t.hideRow, false)
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
			orderedList = append(orderedList, t.columnOrder[i])
		}
	}
	orderedList = append(items, orderedList...)

	t.columnOrder = orderedList

}

// select the column number to hide, columns numbers are the unsorted column number
func (t *Table) HideColumn(columnNumber int) {
	log := logger{location: "Table:HideColumn"}
	log.Debug("Start")

	log.Debug("columnNumber =", columnNumber)
	log.Debug("len(t.head) =", len(t.head))
	if len(t.head) > columnNumber {
		t.head[columnNumber].hidden = true
	} else {
		panic(fmt.Sprintln("invalid column number", columnNumber))
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
		pad := strings.Repeat(" ", t.head[idx].columnLength-len([]rune(word)))
		headLine += fmt.Sprint(word, pad)
	}
	// print the header in one long line
	fmt.Println(strings.TrimRight(headLine, " "))

	// loop through each row
	for r := 0; r < len(t.data); r++ {
		line := ""
		excludeRow := false
		rowNum := t.rowOrder[r]

		if t.hideRow[rowNum] {
			continue
		}

		row := t.data[rowNum]
		// now loop through each column the the currentl selected row
		for col := 0; col < t.headCount; col++ {
			idx := t.columnOrder[col]
			if t.head[idx].hidden {
				continue
			}
			cell := row[idx]
			if len(cell.text) == 0 {
				cell.text = "-"
			}

			// due to looping over every column in the row we only set excludeRow if it is still false
			if !excludeRow {
				// do we have an exclude filter set that we need to process
				excludeRow = t.exclusionFilter(cell, idx)
			}

			spaceCount := t.head[idx].columnLength - len([]rune(cell.text))
			if spaceCount <= 0 {
				spaceCount = maxLineLength
			}
			pad := strings.Repeat(" ", spaceCount)
			line += fmt.Sprint(cell.text, pad)
		}
		if !excludeRow {
			fmt.Println(strings.TrimRight(line, " "))
		}
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
		// now loop through each column for the currently selected row
		for col := 0; col < t.headCount; col++ {
			word := row[col].text
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

// Prints the table on the terminal as yaml, all fileds are shown and all are unsorted as
// other programs can be used to filter and sort
func (t *Table) PrintYaml() {
	// loop through each row
	fmt.Println("data:")
	for rowNum := 0; rowNum < len(t.data); rowNum++ {
		line := ""
		sep := "-"

		row := t.data[rowNum]
		// now loop through each column for the currently selected row
		for col := 0; col < t.headCount; col++ {
			word := row[col].text
			if len(word) == 0 {
				word = ""
			}
			line += fmt.Sprintf("%s %s: \"%s\"\n", sep, t.head[col].title, word)
			sep = " "
		}
		fmt.Print(line)
	}

}

// prints the key and value on a single line by its self. all fileds are shown and all are unsorted as
// other programs can be used to filter and sort
func (t *Table) PrintList() {
	// loop through each row
	for rowNum := 0; rowNum < len(t.data); rowNum++ {
		row := t.data[rowNum]
		// now loop through each column for the currently selected row
		for col := 0; col < t.headCount; col++ {
			word := row[col].text
			if len(word) == 0 {
				word = ""
			}
			fmt.Println(t.head[col].title+":", word)
		}
	}
}

//prints the table as a csv including the header row. all fileds are shown and all are unsorted as
// other programs can be used to filter and sort
func (t *Table) PrintCsv() {

	if len(t.data) <= 0 {
		return
	}

	line := ""
	row := t.data[0]
	// now loop through each column for the currently selected row
	for col := 0; col < t.headCount; col++ {
		word := row[col].text
		if len(word) == 0 {
			word = ""
		}
		line += fmt.Sprintf("\"%s\"", t.head[col].title)
		// add , to the end of every column name except the last
		if col+1 < t.headCount {
			line += ", "
		}
	}
	fmt.Println(line)

	// loop through each column to get the column names
	for rowNum := 0; rowNum < len(t.data); rowNum++ {
		line := ""
		row := t.data[rowNum]
		// now loop through each column for the currently selected row
		for col := 0; col < t.headCount; col++ {
			word := row[col].text
			if len(word) == 0 {
				word = ""
			}
			line += fmt.Sprintf("\"%s\"", word)
			// add , to the end of every key/value except the last
			if col+1 < t.headCount {
				line += ", "
			}
		}

		fmt.Println(line)
	}
}

// Sort via the column number, uses the full column count including hidden columns
// sort function can be run multiple times and is cumalitive
func (t *Table) sort(list []int, columnNumber int, ascending bool) {
	// rather then reordering all rows we have an order array that we can loop through
	// sort contains the actual row number to use next

	// basic bubble sort is used
	for i := 0; i < t.currentRow+1; i++ {
		hasMoved := false
		for j := 0; j < t.currentRow-1; j++ {
			var wordLow, wordHigh string
			var intLow, intHigh int64
			var floatHigh, floatLow float64

			switchOrder := false
			jLow := list[j]
			jHigh := list[j+1]

			switch t.data[jLow][columnNumber].typ {
			case 0:
				wordLow = t.data[jLow][columnNumber].text
				wordHigh = t.data[jHigh][columnNumber].text
			case 1:
				intLow = t.data[jLow][columnNumber].number
				intHigh = t.data[jHigh][columnNumber].number
			case 2:
				floatLow = t.data[jLow][columnNumber].float
				floatHigh = t.data[jHigh][columnNumber].float
			}

			if ascending {
				switch t.data[jLow][columnNumber].typ {
				case 0:
					if wordLow > wordHigh {
						switchOrder = true
					}
				case 1:
					if intLow > intHigh {
						switchOrder = true
					}
				case 2:
					if floatLow > floatHigh {
						switchOrder = true
					}
				}
			} else {
				switch t.data[jLow][columnNumber].typ {
				case 0:
					if wordLow < wordHigh {
						switchOrder = true
					}
				case 1:
					if intLow < intHigh {
						switchOrder = true
					}
				case 2:
					if floatLow < floatHigh {
						switchOrder = true
					}
				}
			}

			if switchOrder {
				hasMoved = true
				list[j] = jHigh
				list[j+1] = jLow
			}
		}
		if !hasMoved {
			break
		}
	}

}

// given a , seperated list of names match them to actual headers and sort each one in order
// by default sorts in ascending to revers use ! in front of the header name
// returns error on fail and nil otherwise
func (t *Table) SortByNames(name ...string) error {
	columnIds := make([]int, len(name))
	columnFound := make([]bool, len(name))
	columnDescend := make([]bool, len(name))

	if len(name) <= 0 {
		return nil
	}

	// scan and match all column names against headers
	for i := 0; i < len(name); i++ {
		rawName := strings.TrimSpace(name[i])
		if len(rawName) <= 0 {
			continue
		}

		// do we need to sort descending
		if strings.HasPrefix(rawName, "!") {
			if len(rawName) == 1 {
				continue
			}
			// remove ! from start of word
			rawName = rawName[1:]
			columnDescend[i] = true
		}

		// loop all header looking for a match
		for c := 0; c < len(t.head); c++ {
			if rawName != t.head[c].title {
				// skip if we dont have a name match
				continue
			}
			// save the matched column id to our array
			columnIds[i] = c
			columnFound[i] = true
		}
	}

	// sort each one in order
	for i := 0; i < len(columnIds); i++ {
		if columnFound[i] {
			// sort function uses ascending true and descending false so we
			// invert descending fLAG to create our ascending flag
			ascend := !columnDescend[i]
			t.sort(t.rowOrder, columnIds[i], ascend)
		}
	}

	return nil
}

// check if matchWord should be excluded using the given filter idx
// return true if matchWord should be excluded and false all other times
func (t *Table) exclusionFilter(matchCell Cell, idx int) bool {
	var fValue float64
	var iValue int64

	exclude := true
	filter := t.head[idx].filter

	// do we have an exclude filter set that we need to process
	if !filter.set {
		return false
	}

	if t.head[idx].columnType == 0 {
		exclude = canExcludeMatchString(filter, matchCell.text, filter.value)
	}

	if t.head[idx].columnType == 1 {
		//convert filter.value to number
		iValue, _ = strconv.ParseInt(filter.value, 10, 64)

		exclude = canExcludeMatchInt(filter, matchCell.number, iValue)
	}

	if t.head[idx].columnType == 2 {
		//convert filter.value to float
		fValue, _ = strconv.ParseFloat(filter.value, 64)

		exclude = canExcludeMatchFloat(filter, matchCell.float, fValue)
	}
	return exclude
}

func canExcludeMatchString(filter matchFilter, val1 string, val2 string) bool {
	//equals
	if filter.compareEql {
		if strMatch(val1, val2) {
			return false
		}
	}

	// not
	if filter.comparison == 3 {
		if !strMatch(val1, val2) {
			return false
		}
	}

	// bigger
	if filter.comparison == 1 {
		if val1 > val2 {
			return false
		}
	}

	// smaller
	if filter.comparison == 2 {
		if val1 < val2 {
			return false
		}
	}

	return true
}

func canExcludeMatchInt(filter matchFilter, val1 int64, val2 int64) bool {
	//equals
	if filter.compareEql {
		if val1 == val2 {
			return false
		}
	}

	//not equals
	if filter.comparison == 3 {
		if val1 != val2 {
			return false
		}
	}

	// bigger
	if filter.comparison == 1 {
		if val1 > val2 {
			return false
		}
	}

	// smaller
	if filter.comparison == 2 {
		if val1 < val2 {
			return false
		}
	}

	return true
}

func canExcludeMatchFloat(filter matchFilter, val1 float64, val2 float64) bool {
	//equals
	if filter.compareEql {
		if val1 == val2 {
			return false
		}
	}

	//not equals
	if filter.comparison == 3 {
		if val1 != val2 {
			return false
		}
	}

	// bigger
	if filter.comparison == 1 {
		if val1 > val2 {
			return false
		}
	}

	// smaller
	if filter.comparison == 2 {
		if val1 < val2 {
			return false
		}
	}

	return true
}

// takes a filter as a string to exclude matching rows from the Print function
// fulter is in the form COLUMN_NAME OPERATOR VALUE, where operator can be one of <,>,<=,>=,!=,=,==
func (t *Table) SetFilter(filter map[string]matchValue) error {

	for words, match := range filter {
		// the smallest header name is T making a valid string "T=0"
		if len(words) < 3 {
			continue
		}

		found := false
		columnName := ""
		operator := ""
		value := ""

		if len(match.operator) > 0 && len(words) > 0 {
			columnName = words
			operator = match.operator
			value = match.value
			found = true
		}

		if found {
			idx := -1
			for i := 0; i < len(t.head); i++ {
				if columnName == t.head[i].title {
					idx = i
					break
				}
			}

			if idx == -1 {
				return errors.New("invalid column name specified")
			}

			switch operator {
			case "=":
				fallthrough
			case "==":
				t.head[idx].filter.comparison = 0
				t.head[idx].filter.compareEql = true

			case "<=":
				t.head[idx].filter.comparison = 2
				t.head[idx].filter.compareEql = true

			case ">=":
				t.head[idx].filter.comparison = 1
				t.head[idx].filter.compareEql = true

			case "<":
				t.head[idx].filter.comparison = 2

			case ">":
				t.head[idx].filter.comparison = 1

			case "!=":
				t.head[idx].filter.comparison = 3
				t.head[idx].filter.compareEql = false

			default:
				return errors.New("invalid operator found")
			}

			if len(value) <= 0 {
				return errors.New("invalid value specified for filter")
			}

			t.head[idx].filter.value = value
			t.head[idx].filter.set = true
		}

	}

	return nil
}

//run a pattten match, accepts * and ?
func strMatch(str string, pattern string) bool {
	// shamelessly converted from c++ code on web as I was too laszy to work it out myself
	// source: https://www.geeksforgeeks.org/wildcard-pattern-matching/

	n := len(str)
	m := len(pattern)

	if m == 0 {
		return (n == 0)
	}

	lookup := make([][]bool, n+1)
	for i := range lookup {
		lookup[i] = make([]bool, m+1)
	}

	lookup[0][0] = true

	for i, char := range pattern {
		j := i + 1
		if char == []rune("*")[0] {
			lookup[0][j] = lookup[0][j-1]
		}
	}

	for q, s := range str {
		i := q + 1
		for w, char := range pattern {
			j := w + 1
			if char == []rune("*")[0] {
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]
			} else if char == []rune("?")[0] || s == char {
				lookup[i][j] = lookup[i-1][j-1]
			} else {
				lookup[i][j] = false
			}
		}
	}
	return lookup[n][m]
}

// quick wrapper to return a cell object containing the given string
func NewCellText(text string) Cell {

	temp := strings.Replace(text, "\r", "\\r", -1)
	temp = strings.Replace(temp, "\f", "\\f", -1)
	temp = strings.Replace(temp, "\n", "\\n", -1)
	temp = strings.Replace(temp, "\t", "\\t", -1)

	return Cell{
		text: temp,
	}
}

// quick wrapper to return a cell object containing the given string and int
func NewCellInt(text string, value int64) Cell {
	return Cell{
		text:   text,
		number: value,
		typ:    1,
	}
}

// quick wrapper to return a cell object containing the given string float
func NewCellFloat(text string, value float64) Cell {
	return Cell{
		text:  text,
		float: value,
		typ:   2,
	}
}

// when given a list of rows and a columnID to work with it will calculate a range and
// returns a list of rows with values outside that range
func (t *Table) ListOutOfRange(columnID int, rows [][]Cell) ([]int, error) {
	var upperFenceInt, lowerFenceInt int64
	var upperFenceFloat, lowerFenceFloat float64

	cellType := rows[0][columnID].typ

	if cellType == 0 {
		return []int{}, errors.New("error: unable to creaate a range with strings")
	}

	orderList := make([]int, len(rows))

	visibleRows := 0
	for i, v := range rows {
		cell := v[columnID]
		orderList[i] = i
		if cellType != cell.typ {
			return []int{}, errors.New("error: table cell types dont match")
		}
		if !t.hideRow[i] {
			visibleRows += 1
		}
	}

	if visibleRows <= 4 {
		return []int{}, errors.New("error: not enough visible rows to calculate useful range")
	}

	t.sort(orderList, columnID, true)
	if cellType == 1 {
		upperFenceInt, lowerFenceInt = t.getFencesInt(orderList, columnID, rows)
	} else {
		upperFenceFloat, lowerFenceFloat = t.getFencesFloat(orderList, columnID, rows)
	}

	out := []int{}

	for k, v := range rows {
		keep := false
		cell := v[columnID]
		if cellType == 1 {
			if upperFenceInt < cell.number {
				keep = true
			}
			if lowerFenceInt > cell.number {
				keep = true
			}
		} else {
			if upperFenceFloat < cell.float {
				keep = true
			}
			if lowerFenceFloat > cell.float {
				keep = true
			}
		}
		if !keep {
			out = append(out, k)
		}
	}

	return out, nil
}

// does what it says on the tin
func (t *Table) GetRows() [][]Cell {
	return t.data
}

// just sets the hide row flag, used by the print function to exclude the row from the output
func (t *Table) HideRows(rowID []int) {
	for _, v := range rowID {
		t.hideRow[v] = true
	}
}

// given the current order and a list of rows caluclate the upper and lower boundy exclusion limit for the selected columnID
func (t *Table) getFencesInt(orderList []int, columnID int, rows [][]Cell) (int64, int64) {
	upper, lower := t.getFencesBoundarys(orderList, columnID, rows, 1)
	return upper.(int64), lower.(int64)
}

// given the current order and a list of rows caluclate the upper and lower boundy exclusion limit for the selected columnID
func (t *Table) getFencesFloat(orderList []int, columnID int, rows [][]Cell) (float64, float64) {
	upper, lower := t.getFencesBoundarys(orderList, columnID, rows, 2)
	return upper.(float64), lower.(float64)
}

// the actual function to caluclate the upper and lower boundy exclusion limit
func (t *Table) getFencesBoundarys(orderList []int, columnID int, rows [][]Cell, cellType int) (interface{}, interface{}) {
	// find middle of the list
	var q1Int, q3Int, iqrInt int64
	var q1Float, q3Float, iqrFloat float64

	// find the middle point in the list so we can split the list into 3
	listLen := len(orderList) + 1
	pos2 := listLen / 2
	pos1 := (pos2 / 2) - 1
	pos3 := pos2 + (pos2 / 2) - 1

	if listLen&1 == 1 { //even list length
		// the middle is held by 2 items, so we grab 2 points for the 1st third
		// and 2 points for the 3rd third
		rowPos1 := orderList[pos1]
		rowPos2 := orderList[pos1+1]
		rowPos3 := orderList[pos3]
		rowPos4 := orderList[pos3+1]

		// grab the values of all 4 points as we need to calulate, to get a single half
		// way value for each third
		t1Cell := rows[rowPos1][columnID]
		t2Cell := rows[rowPos2][columnID]
		t3Cell := rows[rowPos3][columnID]
		t4Cell := rows[rowPos4][columnID]

		// we support both floats and ints so need to calculate based on type
		switch cellType {
		case 1:
			q1Int = (t1Cell.number + t2Cell.number) / 2
			q3Int = (t3Cell.number + t4Cell.number) / 2
		case 2:
			q1Float = (t1Cell.float + t2Cell.float) / 2
			q3Float = (t3Cell.float + t4Cell.float) / 2
		}
	} else { // odd list length
		// odds are eaiser as we have a single middle point, but we still need to deal with floats and ints
		rowPos1 := orderList[pos1]
		rowPos3 := orderList[pos3]
		t1Cell := rows[rowPos1][columnID]
		t3Cell := rows[rowPos3][columnID]

		switch cellType {
		case 1:
			q1Int = t1Cell.number
			q3Int = t3Cell.number
		case 2:
			q1Float = t1Cell.float
			q3Float = t3Cell.float
		}
	}

	// now we can work out the distance between the 1st and 3rd third of the list
	// we calculate 1.5% of that difference and use to create a lower and upper fence
	// these can then be used to exclude everything in side of the 2 fences
	if cellType == 1 {
		iqrInt = q3Int - q1Int
		pc := int64((15 * iqrInt) / 10)
		upperFenceInt := q3Int + pc
		lowerFenceInt := pc - q1Int
		return upperFenceInt, lowerFenceInt
	} else {
		iqrFloat = q3Float - q1Float
		pc := 1.5 * iqrFloat
		upperFenceFloat := q3Float + pc
		lowerFenceFloat := pc - q1Float
		return upperFenceFloat, lowerFenceFloat
	}
}
