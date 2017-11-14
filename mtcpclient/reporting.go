package mtcpclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

func collectConnectionsStatus(connectionsStatusRegistry GroupOfConnections, statusChannel <-chan tcpclient.Connection) {
	for {
		newConnectionStatusReported := <-statusChannel
		connectionsStatusRegistry[newConnectionStatusReported.ID] = newConnectionStatusReported
	}
}

func ReportConnectionsStatus(gc GroupOfConnections, intervalBetweenUpdates int) {
	for {
		fmt.Println(gc)
		if intervalBetweenUpdates == 0 {
			break
		}
		time.Sleep(time.Duration(intervalBetweenUpdates) * time.Second)
	}
}

// StartBackgroundReporting starts some goroutines (so it's not blocking) to capture and report data from the tcpclient
// routines. It initializes and returns the channel that will be used for these communications
func StartBackgroundReporting(numberConnections int, rinterval int) (chan tcpclient.Connection, GroupOfConnections) {
	// A connection may report up to 3 messages: Dialing -> Established -> Closed
	const maxMessagesWeMayGetPerConnection = 3
	connStatusCh := make(chan tcpclient.Connection, numberConnections*maxMessagesWeMayGetPerConnection)
	connStatusTracker := make(GroupOfConnections, numberConnections)

	go ReportConnectionsStatus(connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	return connStatusCh, connStatusTracker
}

func FinalMetricsReport(gc GroupOfConnections) (output string) {
	// Report for Estalished connections
	if gc.AtLeastOneConnectionEstablished() {
		mr := gc.calculateMetricsReport(tcpclient.ConnectionEstablished)
		output = "Timing stats for " + strconv.Itoa(mr.numberOfConnections) +
			" established connections min/avg/max/dev = " +
			mr.min.Truncate(time.Microsecond).String() + "/" +
			mr.avg.Truncate(time.Microsecond).String() + "/" +
			mr.max.Truncate(time.Microsecond).String() + "/" +
			mr.stdDev.Truncate(time.Microsecond).String() + "\n"

	}

	// Report for Errored connections
	if gc.AtLeastOneConnectionInError() {
		mr := gc.calculateMetricsReport(tcpclient.ConnectionError)
		output += "Timing stats for " + strconv.Itoa(mr.numberOfConnections) +
			" failed connections min/avg/max/dev = " +
			mr.min.Truncate(time.Microsecond).String() + "/" +
			mr.avg.Truncate(time.Microsecond).String() + "/" +
			mr.max.Truncate(time.Microsecond).String() + "/" +
			mr.stdDev.Truncate(time.Microsecond).String() + "\n"
	}

	return output
}
