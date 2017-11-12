package mtcpclient

import (
	"fmt"
	"math"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

type GroupOfConnections []tcpclient.Connection

type metricsCollectionStats struct {
	avg    time.Duration
	min    time.Duration
	max    time.Duration
	total  time.Duration
	stdDev time.Duration
}

func (gc GroupOfConnections) calculateMetricsReport(status tcpclient.ConnectionStatus) (mr metricsCollectionStats) {
	for i, item := range gc {
		if i == 0 {
			mr.total = item.GetTCPProcessingDuration(status)
			mr.min = item.GetTCPProcessingDuration(status)
			mr.max = item.GetTCPProcessingDuration(status)
		} else {
			mr.min = time.Duration(math.Min(float64(mr.min), float64(item.GetTCPProcessingDuration(status))))
			mr.max = time.Duration(math.Max(float64(mr.max), float64(item.GetTCPProcessingDuration(status))))
			mr.total += item.GetTCPProcessingDuration(status)
		}
	}
	mr.avg = mr.total / time.Duration(len(gc))
	mr.stdDev = gc.calculateStdDev(status, mr.avg)

	return mr
}

func (gc GroupOfConnections) calculateStdDev(status tcpclient.ConnectionStatus, average time.Duration) time.Duration {
	var sd float64
	for _, item := range gc {
		sd += math.Pow(float64(item.GetTCPProcessingDuration(status))-float64(average), 2)
	}
	return time.Duration(math.Sqrt(sd / float64(time.Duration(len(gc)))))

}

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

func (gc GroupOfConnections) ConnectionInError() bool {
	return gc.isIn(tcpclient.ConnectionError)
}
