package mtcpclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

type GroupOfConnections []tcpclient.Connection

func (gc GroupOfConnections) String() string {
	var nDialing, nEstablished, nClosed, nNotInitiated, nError, nTotal int = 0, 0, 0, 0, 0, 0
	for _, item := range gc {
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
	for _, item := range gc {
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

func appendConnections(gc GroupOfConnections, connections ...tcpclient.Connection) GroupOfConnections {
	// TODO implement as a method using as a reference
	m := len(gc)
	n := m + len(connections)
	if n > cap(gc) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]tcpclient.Connection, (n+1)*2)
		copy(newSlice, gc)
		gc = newSlice
	}
	gc = gc[0:n]
	copy(gc[m:n], connections)
	return gc
}

func (gc GroupOfConnections) getFilteredListByStatus(status []tcpclient.ConnectionStatus) (filteredConnections GroupOfConnections) {
	for _, connection := range gc {
		for _, s := range status {
			if connection.GetConnectionStatus() == s {
				filteredConnections = appendConnections(filteredConnections, connection)
			}
		}
	}
	return filteredConnections
}

func (gc GroupOfConnections) pingStyleReport(status tcpclient.ConnectionStatus) (output string) {
	var introduction string
	var filteredConnections GroupOfConnections
	var mr metricsCollectionStats

	switch status {
	case tcpclient.ConnectionEstablished:
		filteredConnections = gc.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionEstablished, tcpclient.ConnectionClosed})
		mr = filteredConnections.calculateMetricsReport()
		introduction = "Response time stats for " + strconv.Itoa(mr.numberOfConnections) +
			" successful connections min/avg/max/dev = "

	case tcpclient.ConnectionError:
		filteredConnections = gc.getFilteredListByStatus([]tcpclient.ConnectionStatus{tcpclient.ConnectionError})
		mr = filteredConnections.calculateMetricsReport()
		introduction = "Time to error stats for " + strconv.Itoa(mr.numberOfConnections) +
			" failed connections min/avg/max/dev = "
	}
	output = introduction +
		mr.min.Truncate(time.Microsecond).String() + "/" +
		mr.avg.Truncate(time.Microsecond).String() + "/" +
		mr.max.Truncate(time.Microsecond).String() + "/" +
		mr.stdDev.Truncate(time.Microsecond).String() + "\n"

	return output
}
