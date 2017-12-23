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
		connectionsStatusRegistry.connections[newConnectionStatusReported.ID] = newConnectionStatusReported
		if len(connectionsStatusRegistry.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionEstablished})) > connectionsStatusRegistry.metrics.maxConcurrentEstalished{
			connectionsStatusRegistry.metrics.maxConcurrentEstalished = len(connectionsStatusRegistry.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionEstablished}))
		}
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
	connStatusTracker := GroupOfConnections{
		connections: make([]tcpclient.Connection, numberConnections),
		metrics: gcMetrics{
			maxConcurrentEstalished: 0,
		},
	}

	go ReportConnectionsStatus(connStatusTracker, rinterval)
	go collectConnectionsStatus(&connStatusTracker, connStatusCh)

	return connStatusCh, &connStatusTracker
}

func FinalMetricsReport(gc GroupOfConnections) (output string) {
	// Report Established Connections	
	output += "--- Summary of Established connections --- \n" +
		"Total established connections: " +
		strconv.Itoa(len(gc.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionEstablished, tcpclient.ConnectionClosed}))) + "\n" +
		"Max concurrent established connections: " +
		strconv.Itoa(gc.metrics.maxConcurrentEstalished) + "\n" +
		"Last number of established connections: " +
		strconv.Itoa(len(gc.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionEstablished}))) + "\n"

	output += "--- Summary of Timing statistics --- \n"
	// Report for Established connections and also Closed ones
	if gc.AtLeastOneConnectionOK() {
		output += gc.pingStyleReport(tcpclient.ConnectionEstablished)
	}

	// Report for Errored connections
	if gc.AtLeastOneConnectionInError() {
		output += gc.pingStyleReport(tcpclient.ConnectionError)
	}
	return output
}
