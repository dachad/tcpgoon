package mtcpclient

import (
	"github.com/dachad/tcpgoon/tcpclient"
	"testing"
	"time"
)

func TestCalculateStdDev(t *testing.T) {
	var stdDevScenariosChecks = []struct {
		scenarioDescription string
		durationsInSecs     []int
		expectedStdDev      int
	}{
		{
			scenarioDescription: "Empty group of connections should report 0 as stats values",
			durationsInSecs:     []int{},
			expectedStdDev:      0,
		},
		{
			scenarioDescription: "Single connection should report a std dev of 0",
			durationsInSecs:     []int{1},
			expectedStdDev:      0,
		},
		{
			scenarioDescription: "Several connections with same durations should report a std dev of 0",
			durationsInSecs:     []int{1, 1, 1, 1, 1},
			expectedStdDev:      0,
		},
		{
			scenarioDescription: "A known set of durations should report a known std dev",
			durationsInSecs:     []int{1, 2, 3},
			expectedStdDev:      1,
		},
	}
	for _, test := range stdDevScenariosChecks {
		var gc GroupOfConnections = []tcpclient.Connection{}

		var sum int
		for i, connectionDuration := range test.durationsInSecs {
			gc = append(gc, tcpclient.NewConnection(i, tcpclient.ConnectionEstablished,
				time.Duration(connectionDuration)*time.Second))
			sum += connectionDuration
		}

		mr := metricsCollectionStats{}
		if len(test.durationsInSecs) != 0 {
			mr = metricsCollectionStats{
				avg: time.Duration(sum/len(test.durationsInSecs)) * time.Second,
			}
		}

		stddev := gc.calculateStdDev(tcpclient.ConnectionEstablished, mr)

		if stddev != time.Duration(test.expectedStdDev)*time.Second {
			t.Error(test.scenarioDescription+", and its", stddev)
		}
	}
}
