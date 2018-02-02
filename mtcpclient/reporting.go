package mtcpclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

func collectConnectionsStatus(connectionsStatusRegistry *GroupOfConnections, statusChannel <-chan tcpclient.Connection) {
	for {
		newConnectionStatusReported := <-statusChannel
		if newConnectionStatusReported.GetConnectionStatus() == tcpclient.ConnectionEstablished {
			connectionsStatusRegistry.metrics.maxConcurrentEstablished++
		} else if connectionsStatusRegistry.connections[newConnectionStatusReported.ID].GetConnectionStatus() == tcpclient.ConnectionEstablished {
			connectionsStatusRegistry.metrics.maxConcurrentEstablished--
		}
		connectionsStatusRegistry.connections[newConnectionStatusReported.ID] = newConnectionStatusReported
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
func StartBackgroundReporting(numberConnections int, rinterval int) (chan tcpclient.Connection, *GroupOfConnections) {
	// A connection may report up to 3 messages: Dialing -> Established -> Closed
	const maxMessagesWeMayGetPerConnection = 3
	connStatusCh := make(chan tcpclient.Connection, numberConnections*maxMessagesWeMayGetPerConnection)

	var connStatusTracker *GroupOfConnections
	connStatusTracker = newGroupOfConnections(numberConnections)

	go ReportConnectionsStatus(*connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	return connStatusCh, connStatusTracker
}

func FinalMetricsReport(gc GroupOfConnections) (output string) {
	// Report Established Connections
	output += "--- tcpgoon execution statistics --- \n" +
		"Total established connections: " +
		strconv.Itoa(len(gc.getConnectionsThatWentWell().connections)) + "\n" +
		"Max concurrent established connections: " +
		strconv.Itoa(gc.metrics.maxConcurrentEstablished) + "\n" +
		"Number of established connections on closure: " +
		strconv.Itoa(len(gc.getConnectionsThatWentWell().connections)) + "\n"

	output += gc.pingStyleReport()

	return output
}
