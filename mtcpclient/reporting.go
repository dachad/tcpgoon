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
	for {
		fmt.Println(tcpclient.PrintGroupOfConnections(connectionDescriptions))
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

func ReportExecutionSummary(connectionDescriptions []tcpclient.Connection) {
	fmt.Println(printFinalMetricsReport(connectionDescriptions))
}

func printFinalMetricsReport(c []tcpclient.Connection) string {
	var avgToEstablished, minToEstablished, maxToEstablished time.Duration
	var totalToEstalished time.Duration
	for i, item := range c {
		if i == 0 {
			totalToEstalished = tcpclient.TCPProcessingTime(item)
			minToEstablished = tcpclient.TCPProcessingTime(item)
			maxToEstablished = tcpclient.TCPProcessingTime(item)
		} else {
			switch {
			case tcpclient.TCPProcessingTime(item) < minToEstablished:
				minToEstablished = tcpclient.TCPProcessingTime(item)
			case tcpclient.TCPProcessingTime(item) > maxToEstablished:
				maxToEstablished = tcpclient.TCPProcessingTime(item)
			}
			totalToEstalished += tcpclient.TCPProcessingTime(item)
		}
	}
	avgToEstablished = totalToEstalished / time.Duration(len(c))

	return "Time to establish TCP connections min/avg/max = " + minToEstablished.String() + "/" + avgToEstablished.String() + "/" + maxToEstablished.String()
}
