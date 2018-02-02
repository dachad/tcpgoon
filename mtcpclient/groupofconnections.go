package mtcpclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

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

func (gc GroupOfConnections) isIn(status tcpclient.ConnectionStatus) bool {
	for _, item := range gc.connections {
		if item.GetConnectionStatus() == status {
			return true
		}
	}
	return false
}

func (gc GroupOfConnections) PendingConnections() bool {
	return gc.isIn(tcpclient.ConnectionNotInitiated) || gc.isIn(tcpclient.ConnectionDialing)
}

func (gc GroupOfConnections) AtLeastOneConnectionInError() bool {
	return gc.isIn(tcpclient.ConnectionError)
}

func (gc GroupOfConnections) AtLeastOneConnectionOK() bool {
	return gc.isIn(tcpclient.ConnectionEstablished) || gc.isIn(tcpclient.ConnectionClosed)
}

func appendConnections(previous_connections []tcpclient.Connection, new_connections ...tcpclient.Connection) []tcpclient.Connection {
	// TODO implement as a method using as a reference
	m := len(previous_connections)
	n := m + len(new_connections)
	if n > cap(previous_connections) {
		// allocate double what's needed, for future growth.
		newSlice := make([]tcpclient.Connection, (n+1)*2)
		copy(newSlice, previous_connections)
		previous_connections = newSlice
	}
	previous_connections = previous_connections[0:n]
	copy(previous_connections[m:n], new_connections)
	return previous_connections
}

func (gc GroupOfConnections) getFilteredListByStatus(statuses []tcpclient.ConnectionStatus) (filteredConnections []tcpclient.Connection) {
	for _, connection := range gc.connections {
		if connection.IsStatusIn(statuses) {
			filteredConnections = appendConnections(filteredConnections, connection)
		}
	}
	return filteredConnections
}

func (gc GroupOfConnections) pingStyleReport() (output string) {
	var filteredConnections GroupOfConnections
	var mr metricsCollectionStats

	if gc.AtLeastOneConnectionOK() {
		filteredConnections.connections = gc.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionEstablished, tcpclient.ConnectionClosed})
		mr = filteredConnections.calculateMetricsReport()
		output += "Response time stats for " + strconv.Itoa(mr.numberOfConnections) +
			" successful connections min/avg/max/dev = " + printStats(mr)
	}

	if gc.AtLeastOneConnectionInError() {
		filteredConnections.connections = gc.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionError})
		mr = filteredConnections.calculateMetricsReport()
		output += "Time to error stats for " + strconv.Itoa(mr.numberOfConnections) +
			" failed connections min/avg/max/dev = " + printStats(mr)
	}

	return output
}

func printStats(mr metricsCollectionStats) string {
	return mr.min.Truncate(time.Microsecond).String() + "/" +
		mr.avg.Truncate(time.Microsecond).String() + "/" +
		mr.max.Truncate(time.Microsecond).String() + "/" +
		mr.stdDev.Truncate(time.Microsecond).String() + "\n"
}
