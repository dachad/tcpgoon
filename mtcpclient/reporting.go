package mtcpclient

import (
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"fmt"
	"time"
)

func collectConnectionsStatus(connectionDescriptions []tcpclient.Connection, statusChannel <-chan tcpclient.Connection) {
	for {
		connectionStatus := <-statusChannel
		connectionDescriptions[connectionStatus.Id] = connectionStatus

	}
}

func reportConnectionsStatus(connectionDescriptions []tcpclient.Connection, intervalBetweenUpdates int) {
	for {
		fmt.Println(tcpclient.PrintGroupOfConnections(connectionDescriptions))
		time.Sleep(time.Duration(intervalBetweenUpdates) * time.Second)
	}
}

// startReportingLogic starts some goroutines to capture and report data from the tcpclient
// routines. It initializes and returns the channel that will be used for these communications
func StartReportingLogic(numberConnections int, rinterval int) chan tcpclient.Connection {
	// A connection may report up to 3 messages: Dialing -> Established -> Closed
	const maxMessagesWeMayGetPerConnection = 3
	connStatusCh := make(chan tcpclient.Connection, numberConnections*maxMessagesWeMayGetPerConnection)
	connStatusTracker := make([]tcpclient.Connection, numberConnections)

	go reportConnectionsStatus(connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	return connStatusCh
}
