package plugin

import (
	"fmt"
)

var logDebug bool
var dontUseColour bool

type logger struct {
	location string
}

//given log number (int) return the prefix string (ERROR,INFO,DEBUG,etc) and a colour map
// that can be used with fmt.Printf and friends
func logGetType(logType int) (string, string) {

	// Black: [30m
	// Red: [31m
	// Green: [32m
	// Yellow: [33m
	// Blue: [34m
	// Magenta: [35m
	// Cyan: [36m
	// White: [37m

	// Reset: [0m

	const (
		NoColour     = "%s%s"
		SayColour    = "\033[1;32m%s%s\033[0m"
		AskColour    = "\033[0;33m%s%s\033[0m"
		InfoColour   = "\033[0;32m%s%s\033[0m"
		ErrorColour  = "\033[1;31m%s%s\033[0m"
		DebugColour  = "\033[0;36m%s%s\033[0m"
		StdinColour  = "\033[0;35m%s\033[0m%s"
		StdoutColour = "\033[0;36m%s\033[0m%s"
		StderrColour = "\033[0;31m%s%s\033[0m"
	)

	switch logType {
	case 1: //error
		return "ERROR", ErrorColour
	case 2: //info
		return "INFO", InfoColour
	case 3: //debug
		return "DEBUG", DebugColour
	case 6: //stdin
		return "STDIN", StdinColour
	case 7: //stdout
		return "STDOUT", StdoutColour
	case 8: //stderr
		return "STRERR", StderrColour
	}

	//default catch-all
	return "UNKNOWN", NoColour
}

func (l *logger) Error(message ...interface{}) {
	id := 1
	logPrefix, logColour := logGetType(id)

	msg := fmt.Sprintln(message...)
	//dump the message out to the screen
	if dontUseColour {
		l.showLog("", logPrefix+": ", msg)
	} else {
		l.showLog(logColour, logPrefix+": ", msg)
	}
}

func (l *logger) Debug(message ...interface{}) {
	//need to set colours here
	if logDebug {
		id := 3
		logPrefix, logColour := logGetType(id)

		msg := fmt.Sprintln(message...)
		prefix := logPrefix + ":" + l.location + ": "

		//dump the message out to the screen
		if dontUseColour {
			l.showLog("", prefix, msg)
		} else {
			l.showLog(logColour, prefix, msg)
		}
	}
}

//print the log to stdout
func (l *logger) showLog(format string, prefix string, message string) {
	if len(format) == 0 {
		format = "%s%s"
	}

	colourMsg := fmt.Sprintf(format, prefix, message)
	fmt.Print(colourMsg)

}
