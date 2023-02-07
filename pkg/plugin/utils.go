package plugin

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/intstr"
)

const colourEnd = "\033[0m"
const colourNone = -1

// [0] = colour, [1] = modifier // bold,flashing,underline, etc
var colourBad = [2]int{31, 0}
var colourOk = [2]int{32, 0}
var colourWarn = [2]int{33, 0}

// always returns false if the flagList.container is empty as we expect to show all containers
// returns true if we dont have a match
func skipContainerName(flagList commonFlags, containerName string) bool {
	log := logger{location: "Resource"}
	log.Debug("Start")

	if len(flagList.container) == 0 {
		return false
	}

	if flagList.container == containerName {
		return false
	}

	log.Debug("skipping -", containerName)
	return true

}

// returns a memory multiplier that matches the byteType string
func memoryGetUnitLst(byteType string) (int64, string) {
	// Ki | Mi | Gi | Ti | Pi | Ei = 1024 = 1Ki
	// m "" k | M | G | T | P | E = 1000 = 1k
	var d int64 = 1000 // decimal
	var b int64 = 1024 // binary

	memSizes := map[string]int64{
		"Ki": b, "Mi": b * b, "Gi": b * b * b, "Ti": b * b * b * b, "Pi": b * b * b * b * b, "Ei": b * b * b * b * b * b,

		"k": d, "M": d * d, "G": d * d * d, "T": d * d * d * d, "P": d * d * d * d * d, "E": d * d * d * d * d * d,
		"KB": d, "MB": d * d, "GB": d * d * d, "TB": d * d * d * d, "PB": d * d * d * d * d, "EB": d * d * d * d * d * d,
	}

	// limit to two characters
	if len([]rune(byteType)) > 2 {
		byteType = byteType[0:2]
	}

	if len(byteType) > 0 {
		for k, v := range memSizes {
			a := strings.ToLower(k)
			b := strings.ToLower(byteType)
			if a == b {
				return v, k
			}
		}
	}

	return memSizes["M"], "M"
}

// takes a float and converts to a nearest size with unit discriptor as a string
func memoryHumanReadable(memorySize int64, displayAs string) string {
	power := 100.0
	outVal := ""
	if memorySize == 0 {
		return "0"
	}

	multiplier, identifier := memoryGetUnitLst(displayAs)

	size := float64(memorySize) / float64(multiplier)
	val := math.Round(size*power) / power

	outVal = fmt.Sprintf("%.2f%s", val, identifier)

	return outVal
}

// checks if number is NaN, always returns a valid number
func validateFloat64(number float64) float64 {
	if number != number {
		return 0.0
	}
	return number
}

// prints a table on the terminal of a given outType
func outputTableAs(t Table, outType string) {

	switch outType {

	case "":
		t.Print()
	case "csv":
		t.PrintCsv()
	case "list":
		t.PrintList()
	case "json":
		t.PrintJson()
	case "yaml":
		t.PrintYaml()
	}

}

// takes a port object and returns either the number or the name as a string with a proceeding :
// returns empty string if port is empty
func portAsString(port intstr.IntOrString) string {
	// port number provided
	if port.Type == 0 {
		if port.IntVal > 0 {
			return fmt.Sprintf(":%d", port.IntVal)
		} else {
			return ""
		}
	}

	// port name provided
	if port.Type == 1 {
		if len(port.StrVal) > 0 {
			return ":" + port.StrVal
		} else {
			return ""
		}
	}

	return ""
}

// setColourValue set the colour by value, currently 0-74=good, 75-89=warning, 76-100=bad
func setColourValue(value int) [2]int {
	var colour [2]int

	colour = colourOk
	if value > 90 {
		colour = colourBad
	} else if value > 75 {
		colour = colourWarn
	}

	return colour
}

// setColourBoolean set the colour form bool currently: true=good, false=bad
func setColourBoolean(value bool) [2]int {
	var colour [2]int

	if value {
		colour = colourOk
	} else {
		colour = colourBad
	}

	return colour
}

// splitColourString decodes a given colour string item (0.0 or x0.0) into its component parts
//
//	returns state srting (g, w, b) if found, colour, modifier and error state
func splitColourString(colour string) (string, int, int, error) {
	var colourCode int
	var colourMod int

	log := logger{location: "splitColourString"}
	log.Debug("Start")

	rawColour := colour
	rawColourArray := strings.Split(rawColour, "")

	prefixChar := rawColourArray[0]

	_, err := strconv.Atoi(prefixChar)
	if err == nil {
		// we only have a number.number to deal with
		prefixChar = ""
	} else {
		colourArray := rawColourArray[1:len(rawColourArray)]
		rawColour = strings.Join(colourArray, "")
	}

	// we only have a number.number to deal with
	rawColourString := strings.Split(rawColour, ".")

	colourMod, err = strconv.Atoi(rawColourString[0])
	if err != nil {
		return "", 0, 0, errors.New("invalid custom colour modifier")
	}
	colourCode, err = strconv.Atoi(rawColourString[1])
	if err != nil {
		return "", 0, 0, errors.New("invalid custom colour")
	}

	return prefixChar, colourCode, colourMod, nil
}

// getColourSetFromString splits the colour string into colour parts, seperated by ;
//
//	the colours for good, warning and bad are also set when found
//	colourset is set to COLOUR_CUSTOM by default, if g, w or b is found in the colour then COLOUR_CUSTOMMIX is returned instead
//	returns the colours as an array[x][2], colourset and error state
func getColourSetFromString(colours []string) ([][2]int, int, error) {
	var colourArray [][2]int
	colourset := COLOUR_CUSTOM

	log := logger{location: "getColourSetFromString"}
	log.Debug("Start")

	for _, v := range colours {
		if len(v) == 0 {
			continue
		}
		if len(v) <= 3 {
			return [][2]int{}, COLOUR_NONE, errors.New("invalid custom colour detected")
		}

		prefix, code, mod, err := splitColourString(v)
		if err != nil {
			return [][2]int{}, COLOUR_NONE, err
		}

		switch prefix {
		case "g": // good colour code
			colourOk = [2]int{code, mod}
			colourset = COLOUR_CUSTOMMIX
		case "w": // warning colour code
			colourWarn = [2]int{code, mod}
			colourset = COLOUR_CUSTOMMIX
		case "b": // bad colour code
			colourBad = [2]int{code, mod}
			colourset = COLOUR_CUSTOMMIX
		case "":
			colourCode := [2]int{code, mod}
			colourArray = append(colourArray, colourCode)
		}
	}

	if len(colourArray) == 0 {
		colourArray = append(colourArray, [2]int{-1, 0})
	}
	return colourArray, colourset, nil
}
