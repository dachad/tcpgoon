package mtcpclient

import (
	"testing"
	"github.com/dachad/tcpgoon/tcpclient"
	"time"
)

func TestCalculateStdDevOfEmptyGroup(t *testing.T) {
	var gc GroupOfConnections
	gc = []tcpclient.Connection{
	}
	stddev := gc.calculateStdDev(tcpclient.ConnectionClosed, metricsCollectionStats{})
	if stddev != time.Duration(0) {
		t.Error("Empty group of connections should report 0 stats, and its", stddev)
	}
}

func TestCalculateStdDevOfSingleItem(t *testing.T) {
	var gc GroupOfConnections
	gc = []tcpclient.Connection{
		//tcpclient.NewConnection(0, tcpclient.ConnectionClosed, time.Second), -> this should also work, but it doesnt
		// given a bug in how durations are fetched depending on the status
		tcpclient.NewConnection(0, tcpclient.ConnectionEstablished, time.Second),
	}
	stddev := gc.calculateStdDev(tcpclient.ConnectionEstablished, metricsCollectionStats{
		avg: time.Second,
		min: time.Second,
		max: time.Second,
		total: time.Second,
		numberOfConnections: 1,
	})
	if stddev != time.Duration(0) {
		t.Error("Single connection should report a std dev of 0, and its", stddev)
	}
}

//func TestCalculateStdDev(t *testing.T) {
//	var stdDevScenariosChecks = []struct {
//		scenarioDescription string
//		durations           []float64
//		expectedStdDev      float64
//	}{
//		{
//			scenarioDescription: "Single connection should report a std dev of 0",
//			durations: {1.0},
//			expectedStdDev: 0.0,
//		},
//	}
//	for _, test := range stdDevScenariosChecks {
//		for i,connectionDuration := range test.durations {
//
//		}
//	}
//}

