package mtcpclient

import (
	"fmt"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

func collectConnectionsStatus(connectionsStatusRegistry []tcpclient.Connection, statusChannel <-chan tcpclient.Connection) {
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
func StartBackgroundReporting(numberConnections int, rinterval int) (chan tcpclient.Connection, []tcpclient.Connection) {
	// A connection may report up to 3 messages: Dialing -> Established -> Closed
	const maxMessagesWeMayGetPerConnection = 3
	connStatusCh := make(chan tcpclient.Connection, numberConnections*maxMessagesWeMayGetPerConnection)
	connStatusTracker := make([]tcpclient.Connection, numberConnections)

	go ReportConnectionsStatus(connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	return connStatusCh, connStatusTracker
}

func FinalMetricsReport(gc GroupOfConnections) (output string) {
	// Report for Estalished connections
	if gc.AtLeastOneConnectionEstablished() {
		mr := gc.calculateMetricsReport(tcpclient.ConnectionEstablished)
		output = "Time to establish TCP successful connections min/avg/max/stdDev = " +
			mr.min.String() + "/" +
			mr.avg.String() + "/" +
			mr.max.String() + "/" +
			mr.stdDev.String()

	}

	// Report for Errored connections
	if gc.AtLeastOneConnectionInError() {
		mr := gc.calculateMetricsReport(tcpclient.ConnectionError)
		output += "Time spent in failed connections min/avg/max/stdDev = " +
			mr.min.String() + "/" +
			mr.avg.String() + "/" +
			mr.max.String() + "/" +
			mr.stdDev.String()
	}

	return output
}
