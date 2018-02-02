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

func (gc GroupOfConnections) pingStyleReport() (output string) {
	if gc.AtLeastOneConnectionOK() {
		var connectionsThatWentWell GroupOfConnections
		connectionsThatWentWell = gc.getConnectionsThatWentWell()
		output += "Response time stats for " + strconv.Itoa(len(connectionsThatWentWell.connections)) +
			" successful connections min/avg/max/dev = " + printStats(connectionsThatWentWell.calculateMetricsReport())
	}
	if gc.AtLeastOneConnectionInError() {
		var getConnectionsThatWentBad GroupOfConnections
		getConnectionsThatWentBad = gc.getConnectionsThatWentBad()
		output += "Time to error stats for " + strconv.Itoa(len(getConnectionsThatWentBad.connections)) +
			" failed connections min/avg/max/dev = " + printStats(getConnectionsThatWentBad.calculateMetricsReport())
	}
	return output
}

func (gc GroupOfConnections) getConnectionsThatWentWell() (connectionsThatWentWell GroupOfConnections) {
	for _, connection := range gc.connections {
		if connection.WentOk() {
			connectionsThatWentWell.connections = append(connectionsThatWentWell.connections, connection)
		}
	}
	return connectionsThatWentWell
}

func (gc GroupOfConnections) getConnectionsThatWentBad() (connectionsThatWentBad GroupOfConnections) {
	for _, connection := range gc.connections {
		if connection.WentOk() == false {
			connectionsThatWentBad.connections = append(connectionsThatWentBad.connections, connection)
		}
	}
	return connectionsThatWentBad
}

func printStats(mr metricsCollectionStats) string {
	return mr.min.Truncate(time.Microsecond).String() + "/" +
		mr.avg.Truncate(time.Microsecond).String() + "/" +
		mr.max.Truncate(time.Microsecond).String() + "/" +
		mr.stdDev.Truncate(time.Microsecond).String() + "\n"
}
