package mtcpclient

import (
	"fmt"
	"math"
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

func FinalMetricsReport(connectionDescriptions []tcpclient.Connection) string {
	mr := calculateMetricsReport(connectionDescriptions)
	return "Time to establish TCP connections min/avg/max/stdDev = " +
		mr.minToEstablished.String() + "/" +
		mr.avgToEstablished.String() + "/" +
		mr.maxToEstablished.String() + "/" +
		mr.stdDevToEstablished.String()
}

type metricsCollectionStats struct {
	avgToEstablished    time.Duration
	minToEstablished    time.Duration
	maxToEstablished    time.Duration
	totalToEstablished  time.Duration
	stdDevToEstablished time.Duration
}

func calculateMetricsReport(c []tcpclient.Connection) metricsCollectionStats {
	var mr metricsCollectionStats

	for i, item := range c {
		if i == 0 {
			mr.totalToEstablished = tcpclient.TCPProcessingTime(item)
			mr.minToEstablished = tcpclient.TCPProcessingTime(item)
			mr.maxToEstablished = tcpclient.TCPProcessingTime(item)
		} else {
			mr.minToEstablished = time.Duration(math.Min(float64(mr.minToEstablished), float64(tcpclient.TCPProcessingTime(item))))
			mr.maxToEstablished = time.Duration(math.Max(float64(mr.maxToEstablished), float64(tcpclient.TCPProcessingTime(item))))
			mr.totalToEstablished += tcpclient.TCPProcessingTime(item)
		}
	}
	mr.avgToEstablished = mr.totalToEstablished / time.Duration(len(c))

	var sd float64
	for _, item := range c {
		sd += math.Pow(float64(tcpclient.TCPProcessingTime(item))-float64(mr.avgToEstablished), 2)
	}
	mr.stdDevToEstablished = time.Duration(math.Sqrt(sd / float64(time.Duration(len(c)))))

	return mr
}
