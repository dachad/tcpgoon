package mtcpclient

import (
	"fmt"
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"time"
)

func collectConnectionsStatus(connectionDescriptions []tcpclient.Connection, statusChannel <-chan tcpclient.Connection) {
	for {
		connectionStatus := <-statusChannel
		connectionDescriptions[connectionStatus.Id] = connectionStatus
	}
}

func ReportConnectionsStatus(connectionDescriptions []tcpclient.Connection, intervalBetweenUpdates int) {
	for {
		fmt.Println(tcpclient.PrintGroupOfConnections(connectionDescriptions))
		if intervalBetweenUpdates == 0 { break }
		time.Sleep(time.Duration(intervalBetweenUpdates) * time.Second)
	}
}

// startReportingLogic starts some goroutines (so it's not blocking) to capture and report data from the tcpclient
// routines. It initializes and returns the channel that will be used for these communications
func StartBackgroundReporting(numberConnections int, rinterval int) (chan tcpclient.Connection, []tcpclient.Connection) {
	// A connection may report up to 3 messages: Dialing -> Established -> Closed
	const maxMessagesWeMayGetPerConnection = 3
	connStatusCh := make(chan tcpclient.Connection, numberConnections * maxMessagesWeMayGetPerConnection)
	connStatusTracker := make([]tcpclient.Connection, numberConnections)

	go ReportConnectionsStatus(connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	return connStatusCh, connStatusTracker
}
