package mtcpclient

import (
	"testing"
)

func TestFinalMetricsReport(t *testing.T) {
	var finalMetricsReportScenariosChecks = []struct {
		scenarioDescription        string
		groupOfConnectionsToReport *GroupOfConnections
		expectedReport             string
	}{
		{
			scenarioDescription:        "Empty group of connections should report empty values",
			groupOfConnectionsToReport: newGroupOfConnections(0),
			expectedReport: "--- tcpgoon execution statistics --- \n" +
				"Total established connections: 0\n" +
				"Max concurrent established connections: 0\n" +
				"Number of established connections on closure: 0\n",
		},
		{
			scenarioDescription:        "Single connection should generate a report that describes its associated metric",
			groupOfConnectionsToReport: newSampleSingleConnection(),
			expectedReport: "--- tcpgoon execution statistics --- \n" +
				"Total established connections: 1\n" +
				"Max concurrent established connections: 1\n" +
				"Number of established connections on closure: 1\n" +
				"Response time stats for 1 successful connections min/avg/max/dev = 500ms/500ms/500ms/0s\n",
		},
		{
			// TODO: We will need to extend this to cover a mix connections closed + established on closure, when the code supports it
			scenarioDescription:        "Multiple connections with different statuses should generate a report that describes the metrics of the right subset",
			groupOfConnectionsToReport: newSampleMultipleConnections(),
			expectedReport: "--- tcpgoon execution statistics --- \n" +
				"Total established connections: 1\n" +
				"Max concurrent established connections: 1\n" +
				"Number of established connections on closure: 1\n" +
				"Response time stats for 1 successful connections min/avg/max/dev = 500ms/500ms/500ms/0s\n" +
				"Time to error stats for 2 failed connections min/avg/max/dev = 1s/2s/3s/1s\n",
		},
	}

	for _, test := range finalMetricsReportScenariosChecks {
		resultingReport := FinalMetricsReport(*test.groupOfConnectionsToReport)
		if resultingReport != test.expectedReport {
			t.Error(test.scenarioDescription+", and it is:", resultingReport)
		}
	}
}
