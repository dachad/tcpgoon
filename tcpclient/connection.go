package tcpclient

import (
	"fmt"
	"time"
)

type Connection struct {
	ID      int
	status  connectionStatus
	metrics connectionMetrics
}

type connectionStatus int

type connectionMetrics struct {
	tcpEstablishedDuration time.Duration
	// packets lost, retransmissions and other metrics could come
}

const (
	ConnectionNotInitiated connectionStatus = 0
	ConnectionDialing      connectionStatus = 1
	ConnectionEstablished  connectionStatus = 2
	ConnectionClosed       connectionStatus = 3
	ConnectionError        connectionStatus = 4
)

func (c Connection) GetConnectionStatus() connectionStatus {
	return c.status
}

func (cs connectionStatus) isIn(connections []Connection) bool {
	for _, item := range connections {
		if item.status == cs {
			return true
		}
	}
	return false
}

func (c Connection) String() string {
	var status string
	switch c.status {
	case ConnectionNotInitiated:
		status = "not initiated"
	case ConnectionDialing:
		status = "dialing"
	case ConnectionEstablished:
		status = "established"
	case ConnectionClosed:
		status = "closed"
	case ConnectionError:
		status = "errored"
	}

	switch c.status {
	case ConnectionEstablished:
		return fmt.Sprintf("Connection %d has become %s after %s", c.ID, status, c.metrics.tcpEstablishedDuration)
	default:
		return fmt.Sprintf("Connection %d is %s", c.ID, status)
	}

}

func PendingConnections(c []Connection) bool {
	return ConnectionNotInitiated.isIn(c) || ConnectionDialing.isIn(c)
}

func ConnectionInError(c []Connection) bool {
	return ConnectionError.isIn(c)
}

func (c Connection) GetTCPEstablishedDuration() time.Duration {
	return c.metrics.tcpEstablishedDuration
}
