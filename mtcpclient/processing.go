package mtcpclient

import (
	"fmt"
	"math"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

type groupOfConnections []tcpclient.Connection

type metricsCollectionStats struct {
	avgToEstablished    time.Duration
	minToEstablished    time.Duration
	maxToEstablished    time.Duration
	totalToEstablished  time.Duration
	stdDevToEstablished time.Duration
}

func (gc groupOfConnections) calculateMetricsReport() (mr metricsCollectionStats) {
	for i, item := range gc {
		if i == 0 {
			mr.totalToEstablished = item.GetTCPEstablishedDuration()
			mr.minToEstablished = item.GetTCPEstablishedDuration()
			mr.maxToEstablished = item.GetTCPEstablishedDuration()
		} else {
			mr.minToEstablished = time.Duration(math.Min(float64(mr.minToEstablished), float64(item.GetTCPEstablishedDuration())))
			mr.maxToEstablished = time.Duration(math.Max(float64(mr.maxToEstablished), float64(item.GetTCPEstablishedDuration())))
			mr.totalToEstablished += item.GetTCPEstablishedDuration()
		}
	}
	mr.avgToEstablished = mr.totalToEstablished / time.Duration(len(gc))
	mr.stdDevToEstablished = gc.calculateStdDev(mr.avgToEstablished)

	return mr
}

func (gc groupOfConnections) calculateStdDev(average time.Duration) time.Duration {
	var sd float64
	for _, item := range gc {
		sd += math.Pow(float64(item.GetTCPEstablishedDuration())-float64(average), 2)
	}
	return time.Duration(math.Sqrt(sd / float64(time.Duration(len(gc)))))

}

func (gc groupOfConnections) String() string {
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
