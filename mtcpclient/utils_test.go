package mtcpclient

import (
	"time"

	"github.com/dachad/tcpgoon/tcpclient"
)

func newSampleSingleConnection() *GroupOfConnections {
	var gc *GroupOfConnections
	gc = newGroupOfConnections(0)
	gc.connections = append(gc.connections, tcpclient.NewConnection(0, tcpclient.ConnectionEstablished,
		time.Duration(500)*time.Millisecond))
	gc.metrics.maxConcurrentEstablished = 1
	return gc
}

func newSampleMultipleConnections() *GroupOfConnections {
	var gc *GroupOfConnections
	gc = newGroupOfConnections(0)
	gc.connections = append(gc.connections, tcpclient.NewConnection(0, tcpclient.ConnectionEstablished,
		time.Duration(500)*time.Millisecond))
	gc.connections = append(gc.connections, tcpclient.NewConnection(1, tcpclient.ConnectionError,
		time.Duration(1)*time.Second))
	gc.connections = append(gc.connections, tcpclient.NewConnection(2, tcpclient.ConnectionError,
		time.Duration(3)*time.Second))
	gc.metrics.maxConcurrentEstablished = 1
	return gc
}
