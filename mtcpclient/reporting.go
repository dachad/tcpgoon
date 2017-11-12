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

func ReportConnectionsStatus(connectionDescriptions []tcpclient.Connection, intervalBetweenUpdates int) {
	var gc groupOfConnections
	gc = connectionDescriptions
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

func FinalMetricsReport(connectionDescriptions []tcpclient.Connection) string {
	var gc groupOfConnections
	gc = connectionDescriptions
	mr := gc.calculateMetricsReport()
	return "Time to establish TCP connections min/avg/max/stdDev = " +
		mr.minToEstablished.String() + "/" +
		mr.avgToEstablished.String() + "/" +
		mr.maxToEstablished.String() + "/" +
		mr.stdDevToEstablished.String()
}
