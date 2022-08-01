package plugin

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	duration "k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var timestampFormat = "2006-01-02 15:04:05"

var statusShort = "List status of each container in a pod"

var statusDescription = ` Prints container status information from pods, current and previous exit code, reason and signal
are shown slong with current ready and running state. Pods and containers can also be selected
by name. If no name is specified the container state of all pods in the current namespace is
shown.

The T column in the table output denotes S for Standard and I for init containers`

var statusExample = `  # List individual container status from pods
  %[1]s status

  # List conttainers status from pods output in JSON format
  %[1]s status -o json

  # List status from all container in a single pod
  %[1]s status my-pod-4jh36

  # List previous container status from a single pod
  %[1]s status -p my-pod-4jh36

  # List status of all containers named web-container searching all 
  # pods in the current namespace
  %[1]s status -c web-container

  # List status of containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s status -c web-container --sort '!CONTAINER'

  # List status of containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s status -c web-container --sort PODNAME

  # List container status from all pods where label app equals web
  %[1]s status -l app=web

  # List status from all containers where the pods label app is either web or mail
  %[1]s status -l "app in (web,mail)"`

func Status(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {

	log := logger{location: "Status"}
	log.Debug("Start")

	builder := RowBuilder{}
	builder.LoopStatus = true
	builder.ShowInitContainers = true
	builder.PodName = args

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList

	loopinfo := status{}
	builder.Connection = &connect
	builder.SetFlagsFrom(commonFlagList)

	if cmd.Flag("previous").Value.String() == "true" {
		log.Debug("loopinfo.ShowPrevious = true")
		loopinfo.ShowPrevious = true
	}

	if cmd.Flag("details").Value.String() == "true" {
		loopinfo.ShowDetails = true
		builder.ShowContainerType = true
	}

	table := Table{}
	builder.Table = &table
	log.Debug("commonFlagList.showTreeView =", commonFlagList.showTreeView)
	builder.ShowTreeView = commonFlagList.showTreeView

	builder.Build(&loopinfo)

	if !builder.ShowTreeView {
		if !loopinfo.ShowPrevious { // restart count dosent show up when using previous flag
			// do we need to find the outliers, we have enough data to compute a range
			if commonFlagList.showOddities {
				row2Remove, err := table.ListOutOfRange(builder.DefaultHeaderLen + 2) //3 = restarts column
				if err != nil {
					return err
				}
				table.HideRows(row2Remove)
			}
		}
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type status struct {
	ShowPrevious bool
	ShowDetails  bool

	pNotReady     bool //Ready - we use the inverted term so the code makes more sense
	pStopped      bool //Started - we use the inverted term so the code makes more sense
	pRestarts     int64
	pRestartsText string
}

func (s *status) Headers() []string {

	return []string{
		"READY",
		"STARTED",
		"RESTARTS",
		"STATE",
		"REASON",
		"EXIT-CODE",
		"SIGNAL",
		"TIMESTAMP",
		"AGE",
		"MESSAGE",
	}
}

func (s *status) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}
func (s *status) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s *status) HideColumns(info BuilderInformation) []int {
	//"READY","STARTED","RESTARTS","STATE","REASON","EXIT-CODE","SIGNAL","TIMESTAMP","AGE","MESSAGE",
	var hideColumns []int

	if s.ShowDetails {
		hideColumns = append(hideColumns, 8)
	}

	if s.ShowPrevious {
		// remove "READY STARTED RESTARTS AGE" leaving the following
		//  "STATE REASON EXIT-CODE SIGNAL TIMESTAMP MESSAGE"
		hideColumns = append(hideColumns, 0, 1, 2, 8)
	}

	if len(hideColumns) == 0 {
		// hide AGE, MESSAGE
		hideColumns = append(hideColumns, 7, 9)
	}
	// }

	return hideColumns
}

func (s *status) BuildBranch(info BuilderInformation) ([]Cell, error) {

	out := []Cell{
		NewCellText(""),   //ready
		NewCellText(""),   //started
		NewCellInt("", 0), //restarts
		NewCellText(""),   //state
		NewCellText(""),   //reason
		NewCellText(""),   //exit-code
		NewCellText(""),   //signal
		NewCellText(""),   //timestamp
		NewCellText(""),   //age
		NewCellText(""),   //message
	}

	return out, nil
}

func (s *status) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	var age string
	var timestamp string

	phase := string(pod.Status.Phase)
	if pod.Status.StartTime != nil {
		starttime := pod.Status.StartTime.Time
		timestamp = starttime.Format(timestampFormat)
		rawAge := time.Since(starttime)
		age = duration.HumanDuration(rawAge)
	}

	return []Cell{
		NewCellText(""),                       //ready
		NewCellText(""),                       //started
		NewCellInt("0", 0),                    //restarts
		NewCellText(strings.TrimSpace(phase)), //state
		NewCellText(pod.Status.Reason),        //reason
		NewCellText(""),                       //exit-code
		NewCellText(""),                       //signal
		NewCellText(timestamp),                //timestamp
		NewCellText(age),                      //age
		NewCellText(""),                       //message
	}, nil
}

func (s *status) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	var cellList []Cell
	var reason string
	var exitCode string
	var signal string
	var message string
	var startedAt string
	var startTime time.Time
	var skipAgeCalculation bool
	var started string
	var strState string
	var age string
	var state v1.ContainerState
	var rawExitCode, rawSignal, rawRestarts int64

	log := logger{location: "Status:BuildContainerStatus"}
	log.Debug("Start")

	if s.ShowPrevious {
		state = container.LastTerminationState
	} else {
		state = container.State
	}

	if state.Waiting != nil {
		strState = "Waiting"
		reason = state.Waiting.Reason
		message = state.Waiting.Message
		// waiting state dosent have a start time so we skip setting the age variable, used further down
		skipAgeCalculation = true
	}

	if state.Terminated != nil {
		strState = "Terminated"
		exitCode = fmt.Sprintf("%d", state.Terminated.ExitCode)
		rawExitCode = int64(state.Terminated.ExitCode)
		signal = fmt.Sprintf("%d", state.Terminated.Signal)
		rawSignal = int64(state.Terminated.Signal)
		startTime = state.Terminated.StartedAt.Time
		startedAt = state.Terminated.StartedAt.Format(timestampFormat)
		reason = state.Terminated.Reason
		message = state.Terminated.Message
	}

	if state.Running != nil {
		strState = "Running"
		startedAt = state.Running.StartedAt.Format(timestampFormat)
		startTime = state.Running.StartedAt.Time
	}

	if container.Started != nil {
		started = fmt.Sprintf("%t", *container.Started)
		if !*container.Started {
			s.pStopped = true
		}
	}

	ready := fmt.Sprintf("%t", container.Ready)
	if !container.Ready {
		s.pNotReady = true
	}
	restarts := fmt.Sprintf("%d", container.RestartCount)
	rawRestarts = int64(container.RestartCount)

	s.pRestarts += rawRestarts
	s.pRestartsText = fmt.Sprintf("%d", s.pRestarts)

	// remove pod and container name from the message string
	message = s.trimStatusMessage(message, info.PodName, info.Name)

	//we can only show the age if we have a start time some states dont have said starttime so we have to skip them
	if skipAgeCalculation {
		age = ""
	} else {
		rawAge := time.Since(startTime)
		age = duration.HumanDuration(rawAge)
	}

	// READY STARTED RESTARTS STATE REASON EXIT-CODE SIGNAL TIMESTAMP AGE MESSAGE
	cellList = append(cellList,
		NewCellText(ready),
		NewCellText(started),
		NewCellInt(restarts, rawRestarts),
		NewCellText(strState),
		NewCellText(reason),
		NewCellInt(exitCode, rawExitCode),
		NewCellInt(signal, rawSignal),
		NewCellText(startedAt),
		NewCellText(age),
		NewCellText(message),
	)

	log.Debug("len(cellList) =", len(cellList))

	out := make([][]Cell, 1)
	out[0] = cellList
	return out, nil
}

func (s *status) Sum(rows [][]Cell) []Cell {
	rowOut := make([]Cell, 10)

	rowOut[0].text = "true"
	rowOut[1].text = "true"

	//loop through each row in podTotals and add the columns in each row
	for _, r := range rows {
		if r[0].text == "false" {
			fmt.Println("1>", r[0].text)
			// ready = false
			rowOut[0].text = "false" //ready
		}
		if r[1].text == "false" {
			fmt.Println("2>", r[1].text)
			rowOut[1].text = "false" //started
		}
		rowOut[2].number += r[2].number //restarts
	}

	rowOut[2].text = fmt.Sprintf("%d", rowOut[2].number)
	rowOut[3].text = "" //state
	rowOut[4].text = "" //reason
	rowOut[5].text = "" //exit-code
	rowOut[6].text = "" //signal
	rowOut[7].text = "" //timestamp
	rowOut[8].text = "" //age
	rowOut[9].text = "" //message

	return rowOut
}

// Removes the pod name and container name from the status message as its already in the output table
func (s *status) trimStatusMessage(message string, podName string, containerName string) string {

	if len(message) <= 0 {
		return ""
	}
	if len(podName) <= 0 {
		return ""
	}
	if len(containerName) <= 0 {
		return ""
	}

	newMessage := ""
	strArray := strings.Split(message, " ")
	for _, v := range strArray {
		if "container="+containerName == v {
			continue
		}
		if strings.HasPrefix(v, "pod="+podName+"_") {
			continue
		}
		newMessage += " " + v
	}
	return strings.TrimSpace(newMessage)
}
