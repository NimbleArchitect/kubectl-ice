package plugin

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// sets the maximum number of spaces allowed in a column, spaces are clipped to this number
const maxLineLength = 80

type headerRow struct {
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
	typ    int // 0=string, 1=int64, 2=float64, 3=placeholder
	phRef  int // placeholder reference id, used to track the row thats used as a placeholder
	indent int // the number of indents required in the output
	colour int
}

type Table struct {
	currentRow    int
	headCount     int
	columnOrder   []int
	rowOrder      []int
	head          []headerRow
	data          [][]Cell
	hideRow       []bool
	placeHolder   map[int][]Cell
	placeHolderID int
	ColourOutput  int
}

// SetHeader sets the header row to the specified array of strings
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

// AddRow Adds a new row to the end of the table, accepts an array of strings
func (t *Table) AddRow(row ...Cell) {
	log := logger{location: "Table:AddRow"}
	log.Debug("Start")

	if t.headCount > len(row) {
		panic("not enough columns in provided row")
	}

	for i := 0; i < t.headCount; i++ {
		strLen := len([]rune(row[i].text))
		if row[i].indent > 0 {
			strLen += t.indentLen(row[i].indent)
		}
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

// Order changes the order of columns displayed in the table, specifying a subset of the column
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

// HideColumn select the column number to hide, columns numbers are the unsorted column number
func (t *Table) HideColumn(columnNumber int) {
	log := logger{location: "Table:HideColumn"}
	log.Debug("Start")

	log.Debug("columnNumber =", columnNumber)
	log.Debug("len(t.head) =", len(t.head))
	if len(t.head) > columnNumber {
		log.Debug("hide =", t.head[columnNumber].title)
		t.head[columnNumber].hidden = true
	} else {
		panic(fmt.Sprintln("invalid column number", columnNumber))
	}
}

// HideTheseColumns hides the column number to hide, columns numbers are the unsorted column number
func (t *Table) HideOnlyNamedColumns(columnName []string) error {
	var found bool
	var validNames []string

	log := logger{location: "Table:HideOnlyNamedColumns"}
	log.Debug("Start")

	log.Debug("len(columnName) =", len(columnName))
	// // unhide every column
	for i := range t.head {
		t.head[i].hidden = true
		validNames = append(validNames, t.head[i].title)
	}

	// hide only the listed columns
	for _, c := range columnName {
		found = false
		for i, h := range t.head {
			if c == h.title {
				log.Debug("hide =", h.title)
				t.head[i].hidden = false
				found = true
			}
		}
		if !found {
			// 	t.head[i].hidden = true
			return fmt.Errorf("error: invalid column \"%s\" current valid column names are as following %s", c, validNames)
		}
	}
	return nil
}

// Print outputs the table on the terminal, taking the column order and visibiliy into account
func (t *Table) Print() {
	var cellcolour int
	var withColour bool
	var visibleColumns int
	headLine := ""
	colourArray := make([]int, t.headCount)

	if t.ColourOutput != COLOUR_NONE {
		withColour = true

		for i := 0; i < t.headCount; i++ {
			colourArray[i] = int(math.Mod(float64(i), float64(7))) + 30
		}
	}

	// loop through all headers and make a single line properly spaced
	for col := 0; col < t.headCount; col++ {
		// columnOrder contains the actual column number to use next
		idx := t.columnOrder[col]
		if t.head[idx].hidden {
			continue
		}

		cellcolour = colourArray[visibleColumns]
		visibleColumns += 1

		word := t.head[idx].title
		runelen := len([]rune(word))

		if len(word) == 0 {
			word = "-"
		}

		if t.ColourOutput == COLOUR_MIX || t.ColourOutput == COLOUR_COLUMNS {
			word = fmt.Sprintf("\033[%dm%s%s", cellcolour, word, colourEnd)
		}
		pad := strings.Repeat(" ", t.head[idx].columnLength-runelen)

		headLine += fmt.Sprint(word, pad)
	}
	// print the header in one long line
	fmt.Println(strings.TrimRight(headLine, " "))

	// loop through each row
	for r := 0; r < len(t.data); r++ {
		var row []Cell

		visibleColumns = 0
		line := ""
		excludeRow := false
		rowNum := t.rowOrder[r]

		if t.hideRow[rowNum] {
			continue
		}

		if t.data[rowNum][0].typ == 3 {
			row = t.placeHolder[t.data[rowNum][0].phRef]
		} else {
			row = t.data[rowNum]
		}
		// now loop through each column in the currentl selected row
		for col := 0; col < t.headCount; col++ {
			idx := t.columnOrder[col]
			cell := row[idx]

			if t.head[idx].hidden {
				// dont process the row if its hidden
				continue
			}

			cellcolour := colourArray[visibleColumns]
			visibleColumns += 1

			if len(cell.text) == 0 {
				cell.text = "-"
			}

			origtxt := t.indentText(cell.indent, cell.text)
			celltxt := origtxt
			spaceCount := t.head[idx].columnLength - len([]rune(origtxt))
			if spaceCount <= 0 {
				spaceCount = maxLineLength
			}
			pad := strings.Repeat(" ", spaceCount)

			// colour output has been set and the cell has data
			if withColour {
				if t.ColourOutput == COLOUR_MIX || t.ColourOutput == COLOUR_COLUMNS {
					celltxt = fmt.Sprintf("\033[%dm%s%s", cellcolour, origtxt, colourEnd)
				}

				// we check for errors last so it can overwrite the column colours when we are using the mix colour set
				if cell.colour > -1 && (t.ColourOutput == COLOUR_ERRORS || t.ColourOutput == COLOUR_MIX) {
					// error colour set uses red/yellow/green for ok/warning/problem
					if cell.colour == 0 && t.ColourOutput == COLOUR_MIX {
						celltxt = fmt.Sprintf("\033[%dm%s%s", cellcolour, origtxt, colourEnd)
					} else {
						celltxt = fmt.Sprintf("\033[%dm%s%s", cell.colour, origtxt, colourEnd)
					}
				}
			}

			line += fmt.Sprint(celltxt, pad)
		}
		if !excludeRow {
			fmt.Println(strings.TrimRight(line, " "))
		}
	}

}

// PrintJson outputs the table on the terminal as json, all fileds are shown and all are unsorted as
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

// PrintYaml outputs the table on the terminal as yaml, all fileds are shown and all are unsorted as
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

// PrintList outputs the key and value on a single line by its self. all fileds are shown and all are unsorted as
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

// PrintCsv outputs the table as a csv including the header row. all fileds are shown and all are unsorted as
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

// sort Sorts via the column number, uses the full column count including hidden columns
//
//	function can be run multiple times and is cumalitive
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

// SortByNames given a , seperated list of names match them to actual headers and sort each one in order
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

// strMatch run a pattten match, accepts * and ?
func strMatch(str string, pattern string) bool {
	// shamelessly converted from c++ code on web as I was too lazy to work it out myself
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

func NewCellEmpty() Cell {
	return Cell{
		typ:    -1,
		colour: -1,
	}
}

// NewCellText quick wrapper to return a cell object containing the given string
func NewCellText(text string) Cell {

	temp := strings.Replace(text, "\r", "\\r", -1)
	temp = strings.Replace(temp, "\f", "\\f", -1)
	temp = strings.Replace(temp, "\n", "\\n", -1)
	temp = strings.Replace(temp, "\t", "\\t", -1)

	return Cell{
		text:   temp,
		colour: -1,
	}
}

// NewCellTextIndent creates a text cell with an indentation indicator, this dosen't actually indent the cell it just
//
//	tells table.go Print to indent it for us
func NewCellTextIndent(text string, indentLevel int) Cell {

	temp := strings.Replace(text, "\r", "\\r", -1)
	temp = strings.Replace(temp, "\f", "\\f", -1)
	temp = strings.Replace(temp, "\n", "\\n", -1)
	temp = strings.Replace(temp, "\t", "\\t", -1)

	return Cell{
		text:   temp,
		indent: indentLevel,
		colour: -1,
	}
}

// NewCellInt quick wrapper to return a cell object containing the given string and int
func NewCellInt(text string, value int64) Cell {
	return Cell{
		text:   text,
		number: value,
		typ:    1,
		colour: -1,
	}
}

// NewCellFloat quick wrapper to return a cell object containing the given string float
func NewCellFloat(text string, value float64) Cell {
	return Cell{
		text:   text,
		float:  value,
		typ:    2,
		colour: -1,
	}
}

// NewCellColourText quick wrapper to return a cell object containing the given string and the colour to be used
func NewCellColourText(colour int, text string) Cell {

	temp := strings.Replace(text, "\r", "\\r", -1)
	temp = strings.Replace(temp, "\f", "\\f", -1)
	temp = strings.Replace(temp, "\n", "\\n", -1)
	temp = strings.Replace(temp, "\t", "\\t", -1)

	return Cell{
		text:   temp,
		colour: colour,
	}
}

// NewCellColorInt quick wrapper to return a cell object containing the given colour, string and int
func NewCellColourInt(colour int, text string, value int64) Cell {
	return Cell{
		text:   text,
		number: value,
		typ:    1,
		colour: colour,
	}
}

// NewCellFloat quick wrapper to return a cell object containing the given colour, string and float
func NewCellColourFloat(colour int, text string, value float64) Cell {
	return Cell{
		text:   text,
		float:  value,
		typ:    2,
		colour: colour,
	}
}

// ListOutOfRange when given a columnID to work with it will calculate a range and
// returns a list of rows with values outside that range
func (t *Table) ListOutOfRange(columnID int) ([]int, error) {
	var upperFenceInt, lowerFenceInt int64
	var upperFenceFloat, lowerFenceFloat float64

	cellType := t.data[0][columnID].typ

	if cellType == 0 {
		return []int{}, errors.New("error: unable to creaate a range with strings")
	}

	orderList := make([]int, len(t.data))

	visibleRows := 0
	for i, v := range t.data {
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
		upperFenceInt, lowerFenceInt = t.getFencesInt(orderList, columnID, t.data)
	} else {
		upperFenceFloat, lowerFenceFloat = t.getFencesFloat(orderList, columnID, t.data)
	}

	out := []int{}

	for k, v := range t.data {
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

// GetRows does what it says on the tin
func (t *Table) GetRows() [][]Cell {
	return t.data
}

// HideRows just sets the hide row flag, used by the print function to exclude the row from the output
func (t *Table) HideRows(rowID []int) {
	for _, v := range rowID {
		t.hideRow[v] = true
	}
}

// getFencesInt given the current order and a list of rows caluclate the upper and lower boundy exclusion limit for the selected columnID
func (t *Table) getFencesInt(orderList []int, columnID int, rows [][]Cell) (int64, int64) {
	upper, lower := t.getFencesBoundarys(orderList, columnID, rows, 1)
	return upper.(int64), lower.(int64)
}

// getFencesFloat given the current order and a list of rows caluclate the upper and lower boundy exclusion limit for the selected columnID
func (t *Table) getFencesFloat(orderList []int, columnID int, rows [][]Cell) (float64, float64) {
	upper, lower := t.getFencesBoundarys(orderList, columnID, rows, 2)
	return upper.(float64), lower.(float64)
}

// getFencesBoundarys the actual function to caluclate the upper and lower boundy exclusion limit
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

// AddPlaceHolderRow - Adds an updatable row to the table, returns an update id that can be used with UpdatePlaceHolderRow
func (t *Table) AddPlaceHolderRow() int {
	var cellRow []Cell

	id := t.placeHolderID
	t.placeHolderID++

	for i := 0; i < t.headCount; i++ {
		cellRow = append(cellRow, Cell{
			// text:  "PH" + fmt.Sprint(id),
			typ:   3,
			phRef: id,
		})
	}

	t.AddRow(cellRow...)
	if len(t.placeHolder) == 0 {
		t.placeHolder = make(map[int][]Cell, 1)
	}
	t.placeHolder[id] = cellRow

	return id
}

// UpdatePlaceHolderRow - updates the given placeholder at id with the contents of cellList
func (t *Table) UpdatePlaceHolderRow(id int, cellList []Cell) {

	for i := 0; i < t.headCount; i++ {
		strLen := len([]rune(cellList[i].text))
		if cellList[i].indent > 0 {
			strLen += t.indentLen(cellList[i].indent)
		}
		if strLen >= t.head[i].columnLength {
			if (strLen + 2) > maxLineLength {
				t.head[i].columnLength = maxLineLength
			} else {
				t.head[i].columnLength = strLen + 2
			}
		}
	}
	t.placeHolder[id] = cellList
}

// HidePlaceHolderRow matches the placeholder id to an actual row number and calls HideRows to hide the row
func (t *Table) HidePlaceHolderRow(id int) {
	for r := 0; r < len(t.data); r++ {
		rowNum := t.rowOrder[r]

		if t.data[rowNum][0].phRef == id {
			t.HideRows([]int{r})
		}
	}
}

// indentText indents the text to the specified level adds └─ for every level above 0
func (t *Table) indentText(level int, data string) string {
	var indent string

	if level == 0 {
		return data
	}

	if level == 1 {
		indent = "└─"
	}

	if level >= 2 {
		indent = strings.Repeat(" ", level) + "└─"
	}

	return fmt.Sprint(indent, data)
}

// indentLen returns the number of characters that would be indented at the provided level
func (t *Table) indentLen(level int) int {
	var indent int

	if level == 0 {
		return 0
	}

	if level == 1 {
		indent = 2
	}

	if level >= 2 {
		indent = level + 2
	}

	return indent
}
