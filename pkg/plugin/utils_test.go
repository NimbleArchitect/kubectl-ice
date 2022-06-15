package plugin

import (
	"math"
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"
)

//*****************
//skipContainerName
//*****************
type skipContainerNameTest struct {
	arg1     commonFlags
	arg2     string
	expected bool
}

var skipContainerNameTests = []skipContainerNameTest{
	{commonFlags{container: ""}, "thisname", false},
	{commonFlags{container: "notthis"}, "thisname", true},
	{commonFlags{container: "thisname"}, "thisname", false},
	{commonFlags{container: "notthis"}, "", true},
}

func TestSkipContainerName(t *testing.T) {

	for _, test := range skipContainerNameTests {
		if output := skipContainerName(test.arg1, test.arg2); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}

}

//********************
//skipmemoryGetUnitLst
//********************
type skipMemoryGetUnitLstTest struct {
	arg1      string
	expected1 int64
	expected2 string
}

var skipMemoryGetUnitLstTests = []skipMemoryGetUnitLstTest{
	//all empty or invalid tests
	{"", 1000000, "M"},
	{"toobig", 1000000, "M"},
	{"1", 1000000, "M"},
	{"A", 1000000, "M"},
	{"Aa", 1000000, "M"},
	{"aA", 1000000, "M"},
	{"aa", 1000000, "M"},
	{"AA", 1000000, "M"},

	//valid tests
	{"k", 1000, "k"},
	{"K", 1000, "k"},
	{"M", 1000000, "M"},
	{"m", 1000000, "M"},
	{"G", 1000000000, "G"},
	{"g", 1000000000, "G"},
	{"T", 1000000000000, "T"},
	{"t", 1000000000000, "T"},
	{"P", 1000000000000000, "P"},
	{"p", 1000000000000000, "P"},
	{"E", 1000000000000000000, "E"},
	{"e", 1000000000000000000, "E"},

	{"ki", 1024, "Ki"},
	{"Ki", 1024, "Ki"},
	{"Mi", 1048576, "Mi"},
	{"mi", 1048576, "Mi"},
	{"Gi", 1073741824, "Gi"},
	{"gi", 1073741824, "Gi"},
	{"Ti", 1099511627776, "Ti"},
	{"ti", 1099511627776, "Ti"},
	{"Pi", 1125899906842624, "Pi"},
	{"pi", 1125899906842624, "Pi"},
	{"Ei", 1152921504606846976, "Ei"},
	{"ei", 1152921504606846976, "Ei"},
}

func TestMemoryGetUnitLst(t *testing.T) {

	for _, test := range skipMemoryGetUnitLstTests {
		output1, output2 := memoryGetUnitLst(test.arg1)
		if output1 != test.expected1 {
			t.Errorf("Output1 %d not equal to expected %d", output1, test.expected1)
		}
		if output2 != test.expected2 {
			t.Errorf("Output2 %s not equal to expected %s", output2, test.expected2)
		}
	}

}

//*******************
//memoryHumanReadable
//*******************
type memoryHumanReadableTest struct {
	arg1     int64
	arg2     string
	expected string
}

var skipMemoryHumanReadableTests = []memoryHumanReadableTest{
	//all empty or invalid tests
	{0, "", "0"},
	{1, "b", "0.00M"},
	{-1, "b", "-0.00M"},

	//valid tests
	{1000000, "K", "1000.00k"},
	{1678, "k", "1.68k"},
	{1678 * 1024 * 1024, "m", "1759.51M"},
	{5678 * 1024, "m", "5.81M"},
	{5678 * 1024, "g", "0.01G"},
}

func TestMemoryHumanReadable(t *testing.T) {

	for _, test := range skipMemoryHumanReadableTests {
		output := memoryHumanReadable(test.arg1, test.arg2)
		if output != test.expected {
			t.Errorf("Output1 %s not equal to expected %s, using input %d,%s", output, test.expected, test.arg1, test.arg2)
		}
	}
}

//*******************
//validateFloat64
//*******************
func TestValidateFloat64(t *testing.T) {
	input := 1234.1234
	output := validateFloat64(input)
	if output != input {
		t.Errorf("Output1 %f not equal to expected %f", output, input)
	}

	input = math.NaN()
	output = validateFloat64(input)
	if output != 0.0 {
		t.Errorf("Output1 %f not equal to expected %f", output, input)
	}
}

//*******************
//portAsString
//*******************
type portAsStringTest struct {
	arg1     intstr.IntOrString
	expected string
}

var portAsStringTests = []portAsStringTest{
	{intstr.IntOrString{Type: 0, IntVal: 0, StrVal: ""}, ""},
	{intstr.IntOrString{Type: 0, IntVal: 80, StrVal: ""}, ":80"},
	{intstr.IntOrString{Type: 0, IntVal: 81, StrVal: "http"}, ":81"},
	{intstr.IntOrString{Type: 1, IntVal: 82, StrVal: "http"}, ":http"},
	{intstr.IntOrString{Type: 1, IntVal: 83, StrVal: ""}, ""},
	{intstr.IntOrString{Type: 2, IntVal: 84, StrVal: ""}, ""},
}

func TestDemo(t *testing.T) {
	for _, test := range portAsStringTests {
		output := portAsString(test.arg1)
		if output != test.expected {
			t.Errorf("Output %s not equal to expected %s", output, test.expected)
		}
	}
}
