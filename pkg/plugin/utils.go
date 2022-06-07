package plugin

import (
	"fmt"
	"math"
	"strings"

	"k8s.io/apimachinery/pkg/util/intstr"
)

// always returns false if the flagList.container is empty as we expect to show all containers
// returns true if we dont have a match
func skipContainerName(flagList commonFlags, containerName string) bool {
	if len(flagList.container) == 0 {
		return false
	}

	if flagList.container == containerName {
		return false
	}

	return true

}

//returns a memory multiplier that matches the byteType string
func memoryGetUnitLst(byteType string) (int64, string) {
	// Ki | Mi | Gi | Ti | Pi | Ei = 1024 = 1Ki
	// m "" k | M | G | T | P | E = 1000 = 1k
	var d int64 = 1000 // decimal
	var b int64 = 1024 // binary

	memSizes := map[string]int64{
		"Ki": b, "Mi": b * b, "Gi": b * b * b, "Ti": b * b * b * b, "Pi": b * b * b * b * b, "Ei": b * b * b * b * b * b,
		"k": d, "M": d * d, "G": d * d * d, "T": d * d * d * d, "P": d * d * d * d * d, "E": d * d * d * d * d * d,
	}

	if len(byteType) > 0 && len(byteType) <= 2 {
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
	var floatfmt string

	power := 100.0
	outVal := ""

	if memorySize == 0 {
		return "0"
	}

	multiplier, identifier := memoryGetUnitLst(displayAs)

	// convertTo = "M"
	// if len(displayAs) <= 2 {
	// 	convertTo = displayAs
	// }

	// if val, ok := byteList[convertTo]; ok {
	// 	multiplier = val
	// } else {
	// 	multiplier = byteList["M"]
	// }

	size := float64(memorySize) / float64(multiplier)
	val := math.Round(size*power) / power

	remain := int64(math.Round(size*power)) % int64(power)
	if remain == 0 {
		floatfmt = "%.2f%s"
	} else {
		floatfmt = "%.2f%s"
	}

	outVal = fmt.Sprintf(floatfmt, val, identifier)

	return outVal
}

//checks if number is NaN, always returns a valid number
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
	//port number provided
	if port.Type == 0 {
		if port.IntVal > 0 {
			return fmt.Sprintf(":%d", port.IntVal)
		} else {
			return ""
		}
	}

	//port name provided
	if port.Type == 1 {
		if len(port.StrVal) > 0 {
			return ":" + port.StrVal
		} else {
			return ""
		}
	}

	return ""
}
