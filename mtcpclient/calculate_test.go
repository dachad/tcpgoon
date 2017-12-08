package mtcpclient

import (
	"testing"
	"github.com/dachad/tcpgoon/tcpclient"
	"time"
)

func TestCalculateStdDevOfEmptyGroup(t *testing.T) {
	var gc GroupOfConnections
	gc = []tcpclient.Connection{
		{},
	}
	stddev := gc.calculateStdDev(tcpclient.ConnectionNotInitiated, metricsCollectionStats{})
	if  stddev != time.Duration(0) {
		t.Error("Empty group of connections should report 0 stats, and its", stddev)
	}
}

func TestCalculateStdDevOfSingleItem(t *testing.T) {
	var gc GroupOfConnections
	gc = []tcpclient.Connection{
		tcpclient.NewConnection(0, tcpclient.ConnectionClosed, time.Second),
	}
	stddev := gc.calculateStdDev(tcpclient.ConnectionClosed, metricsCollectionStats{
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

