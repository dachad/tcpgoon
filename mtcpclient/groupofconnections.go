package mtcpclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

// GroupOfConnections aggregates all the running connections plus some general metrics
type GroupOfConnections struct {
	connections []tcpclient.Connection
	metrics     gcMetrics
}

type gcMetrics struct {
	maxConcurrentEstablished int
}

func newGroupOfConnections(numberConnections int) *GroupOfConnections {
	gc := new(GroupOfConnections)
	gc.connections = make([]tcpclient.Connection, numberConnections)
	gc.metrics = gcMetrics{
		maxConcurrentEstablished: 0,
	}
	return gc
}

func (gc GroupOfConnections) String() string {
	var nDialing, nEstablished, nClosed, nNotInitiated, nError, nTotal int = 0, 0, 0, 0, 0, 0
	for _, item := range gc.connections {
		switch item.GetConnectionStatus() {
		case tcpclient.ConnectionDialing:
			nDialing++
		case tcpclient.ConnectionEstablished:
			nEstablished++
		case tcpclient.ConnectionClosed:
			nClosed++
		case tcpclient.ConnectionError:
			nError++
		case tcpclient.ConnectionNotInitiated:
			nNotInitiated++
		}
		nTotal++
	}
	return fmt.Sprintf("Total: %d, Dialing: %d, Established: %d, Closed: %d, Error: %d, NotInitiated: %d",
		nTotal, nDialing, nEstablished, nClosed, nError, nNotInitiated)
}

func (gc GroupOfConnections) containsAConnectionWithStatus(status string) bool {
	for _, connection := range gc.connections {
		switch status {
		case "pending":
			return connection.PendingToProcess()
		case "error":
			return connection.WithError()
		case "established":
			return connection.WentOk()
		}
	}
	return false
}

// PendingConnections retuns True if at least one connection is being processed
func (gc GroupOfConnections) PendingConnections() bool {
	return gc.containsAConnectionWithStatus("pending")
}

// AtLeastOneConnectionInError returns True is at least one connection establishment failed
func (gc GroupOfConnections) AtLeastOneConnectionInError() bool {
	return gc.containsAConnectionWithStatus("error")
}

func (gc GroupOfConnections) atLeastOneConnectionOK() bool {
	return gc.containsAConnectionWithStatus("established")
}

func (gc GroupOfConnections) pingStyleReport(typeOfReport string) (output string) {
	var headerline string
	switch typeOfReport {
	case "successful":
		headerline = "Response time"
	case "failed":
		headerline = "Time to error"
	}

	output += headerline + " stats for " + strconv.Itoa(len(gc.connections)) + " " + typeOfReport +
		" connections min/avg/max/dev = " + printStats(gc.calculateMetricsReport())

	return output
}

func (gc GroupOfConnections) getConnectionsThatWentWell(itWentWell bool) (connectionsThatWent GroupOfConnections) {
	for _, connection := range gc.connections {
		if connection.WentOk() == itWentWell {
			connectionsThatWent.connections = append(connectionsThatWent.connections, connection)
		}
	}
	return connectionsThatWent
}

func (gc GroupOfConnections) getConnectionsThatAreOk() (connectionsThatAreOk GroupOfConnections) {
	for _, connection := range gc.connections {
		if connection.IsOk() {
			connectionsThatAreOk.connections = append(connectionsThatAreOk.connections, connection)
		}
	}
	return connectionsThatAreOk
}

func printStats(mr *metricsCollectionStats) string {
	return mr.min.Truncate(time.Microsecond).String() + "/" +
		mr.avg.Truncate(time.Microsecond).String() + "/" +
		mr.max.Truncate(time.Microsecond).String() + "/" +
		mr.stdDev.Truncate(time.Microsecond).String() + "\n"
}
