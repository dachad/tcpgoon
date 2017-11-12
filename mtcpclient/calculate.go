package mtcpclient

import (
	"math"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

type metricsCollectionStats struct {
	avg    time.Duration
	min    time.Duration
	max    time.Duration
	total  time.Duration
	stdDev time.Duration
}

func (gc GroupOfConnections) calculateMetricsReport(status tcpclient.ConnectionStatus) (mr metricsCollectionStats) {
	mr.total = 0
	for _, item := range gc {
		if item.GetConnectionStatus() == status {
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
