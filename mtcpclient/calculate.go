package mtcpclient

import (
	"math"
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
	"fmt"
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
	// TODO: There's something i don't like... initiatlizing values in the first loop, and the standard deviation
	//  requiring an extra pass considering all items... i'd move initialization out of the loop, and maybe iterate
	//  over a filtered list rather than several loops over the original one, and maybe use specific generic functions...
	//  requires further thinking in any case
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
	// TODO: passing the whole mr struct looks overkilling, given we only want a single value, the average, and maybe
	//  we can actually use a version of the algorithm that calculates it (and the number of items)
	var nitems int = 0
	var sd float64
	for _, item := range gc {
		if item.GetConnectionStatus() == status {
			nitems += 1
			sd += math.Pow(float64(item.GetTCPProcessingDuration(status))-float64(mr.avg), 2)
		}
	}
	if nitems == 0 {
		return 0
	}
	fmt.Println("sd:", sd, "nitems:", nitems)
	return time.Duration(math.Sqrt(sd / float64(nitems)))

}
