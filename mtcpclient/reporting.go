package mtcpclient

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

func collectConnectionsStatus(connectionsStatusRegistry *GroupOfConnections, statusChannel <-chan tcpclient.Connection) {
	concurrentEstablished := 0
	for {
		newConnectionStatusReported := <-statusChannel
		concurrentEstablished = updateConcurrentEstablished(concurrentEstablished, newConnectionStatusReported, connectionsStatusRegistry)
		connectionsStatusRegistry.connections[newConnectionStatusReported.ID] = newConnectionStatusReported
	}
}

func updateConcurrentEstablished(concurrentEstablished int, newConnectionStatusReported tcpclient.Connection, connectionsStatusRegistry *GroupOfConnections) int {
	if tcpclient.IsOk(newConnectionStatusReported) {
		concurrentEstablished++
		connectionsStatusRegistry.metrics.maxConcurrentEstablished = int(math.Max(float64(concurrentEstablished),
			float64(connectionsStatusRegistry.metrics.maxConcurrentEstablished)))
	} else if tcpclient.IsOk(connectionsStatusRegistry.connections[newConnectionStatusReported.ID]) {
		concurrentEstablished--
	}
	return concurrentEstablished
}

// ReportConnectionsStatus keeps printing on screen the summary of connections states
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

	connStatusTracker := newGroupOfConnections(numberConnections)

	go ReportConnectionsStatus(*connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	return connStatusCh, connStatusTracker
}

// FinalMetricsReport creates the final reporting summary
func FinalMetricsReport(gc GroupOfConnections) (output string) {
	// Report Established Connections
	output += "--- tcpgoon execution statistics ---\n" +
		"Total established connections: " +
		strconv.Itoa(len(gc.getConnectionsThatWentWell(true).connections)) + "\n" +
		"Max concurrent established connections: " +
		strconv.Itoa(gc.metrics.maxConcurrentEstablished) + "\n" +
		"Number of established connections on closure: " +
		strconv.Itoa(len(gc.getConnectionsThatAreOk().connections)) + "\n"

	if gc.atLeastOneConnectionOK() {
		output += gc.getConnectionsThatWentWell(true).pingStyleReport(successfulExecution)
	}
	if gc.AtLeastOneConnectionInError() {
		output += gc.getConnectionsThatWentWell(false).pingStyleReport(failedExecution)
	}

	return output
}
