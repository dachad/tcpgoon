package mtcpclient

import (
	"math"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

type metricsCollectionStats struct {
	avg                 time.Duration
	min                 time.Duration
	max                 time.Duration
	total               time.Duration
	stdDev              time.Duration
	numberOfConnections int
}

func (gc GroupOfConnections) calculateMetricsReport(status tcpclient.ConnectionStatus) (mr metricsCollectionStats) {
	mr.total = 0
	mr.numberOfConnections = 0
	for _, item := range gc {
		if item.GetConnectionStatus() == status {
			mr.numberOfConnections++
			if mr.total == 0 {
				mr.total = item.GetTCPProcessingDuration(status)
				mr.min = item.GetTCPProcessingDuration(status)
				mr.max = item.GetTCPProcessingDuration(status)
			} else {
				mr.min = time.Duration(math.Min(float64(mr.min), float64(item.GetTCPProcessingDuration(status))))
				mr.max = time.Duration(math.Max(float64(mr.max), float64(item.GetTCPProcessingDuration(status))))
				mr.total += item.GetTCPProcessingDuration(status)
			}
		}
	}
	mr.avg = mr.total / time.Duration(mr.numberOfConnections)
	mr.stdDev = gc.calculateStdDev(status, mr)
	return mr
}

func (gc GroupOfConnections) calculateStdDev(status tcpclient.ConnectionStatus, mr metricsCollectionStats) time.Duration {
	var sd float64

	for _, item := range gc {
		if item.GetConnectionStatus() == status {
			sd += math.Pow(float64(item.GetTCPProcessingDuration(status))-float64(mr.avg), 2)
		}
	}
	return time.Duration(math.Sqrt(sd / float64(time.Duration(mr.numberOfConnections))))

}
