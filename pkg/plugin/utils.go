package plugin

import (
	"fmt"
	"math"
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

//returns a list of memory sizes with their multipacation amount
func memoryGetUnitLst() map[string]int64 {
	// Ki | Mi | Gi | Ti | Pi | Ei = 1024 = 1Ki
	// m "" k | M | G | T | P | E = 1000 = 1k
	var d int64 = 1000 // decimal
	var b int64 = 1024 // binary

	return map[string]int64{
		"Ki": b, "Mi": b * b, "Gi": b * b * b, "Ti": b * b * b * b, "Pi": b * b * b * b * b, "Ei": b * b * b * b * b * b,
		"k": d, "M": d * d, "G": d * d * d, "T": d * d * d * d, "P": d * d * d * d * d, "E": d * d * d * d * d * d,
	}
}

// takes a float and converts to a nearest size with unit discriptor as a string
func memoryHumanReadable(memorySize int64) string {
	var floatfmt string
	power := 100.0
	outVal := ""

	if memorySize == 0 {
		return "0"
	}

	byteList := memoryGetUnitLst()

	for k, v := range byteList {
		if len(k) == 2 {
			size := float64(memorySize) / float64(v)
			val := math.Round(size*power) / power

			remain := int64(math.Round(size*power)) % int64(power)
			if remain == 0 {
				floatfmt = "%.2f%s"
			} else {
				floatfmt = "%.2f%s"
			}

			// TODO: it works but its clunky and a bit crap, needs work :(
			if val > 0.0 && val <= 900 {
				outVal = fmt.Sprintf(floatfmt, val, k)
			}
			if val > 0.9 && val <= 900 {
				outVal = fmt.Sprintf(floatfmt, val, k)
			}
		}
	}
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
	case "json":
		t.PrintJson()
	}

}
