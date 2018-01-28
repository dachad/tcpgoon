package mtcpclient_test

import (
	"testing"
	"github.com/dachad/tcpgoon/tcpclient"
	"time"
	"github.com/dachad/tcpgoon/mtcpclient"
)

func TestFinalMetricsReport(t *testing.T)  {
	var finalMetricsReportScenariosChecks = []struct {
		scenarioDescription        string
		groupOfConnectionsToReport mtcpclient.GroupOfConnections
		expectedReport             string
	}{
		{
			scenarioDescription:        "Empty group of connections should report nothing",
			groupOfConnectionsToReport: mtcpclient.GroupOfConnections{},
			expectedReport: "\n",
		},
		{
			scenarioDescription:        "Single connection should generate a report that describes its associated metric",
			groupOfConnectionsToReport: mtcpclient.GroupOfConnections{
				tcpclient.NewConnection(0, tcpclient.ConnectionEstablished, time.Duration(500) * time.Millisecond),
			},
			expectedReport: "Response time stats for 1 established connections min/avg/max/dev = 500ms/500ms/500ms/0s\n",
		},
		{
			// TODO: We will need to extend this to cover a mix connections closed + established on closure, when the code supports it
			scenarioDescription:        "Multiple connections with different statuses should generate a report that describes the metrics of the right subset",
			groupOfConnectionsToReport: mtcpclient.GroupOfConnections{
				tcpclient.NewConnection(0, tcpclient.ConnectionEstablished, time.Duration(500) * time.Millisecond),
				tcpclient.NewConnection(1, tcpclient.ConnectionError, time.Duration(1) * time.Second),
				tcpclient.NewConnection(2, tcpclient.ConnectionError, time.Duration(3) * time.Second),
			},
			expectedReport: "Response time stats for 1 established connections min/avg/max/dev = 500ms/500ms/500ms/0s\n" +
				"Time to error stats for 2 failed connections min/avg/max/dev = 1s/2s/3s/1s\n",
		},
	}

	for _, test := range finalMetricsReportScenariosChecks {
		resultingReport := mtcpclient.FinalMetricsReport(test.groupOfConnectionsToReport)
		if resultingReport != test.expectedReport {
			t.Error(test.scenarioDescription + ", and it is:", resultingReport)
		}
	}
}
