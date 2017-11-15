package tcpclient

import (
	"fmt"
	"time"
)

type Connection struct {
	ID      int
	status  ConnectionStatus
	metrics connectionMetrics
}

type ConnectionStatus int

type connectionMetrics struct {
	tcpEstablishedDuration time.Duration
	tcpErroredDuration     time.Duration
	// packets lost, retransmissions and other metrics could come
}

const (
	ConnectionNotInitiated ConnectionStatus = 0
	ConnectionDialing      ConnectionStatus = 1
	ConnectionEstablished  ConnectionStatus = 2
	ConnectionClosed       ConnectionStatus = 3
	ConnectionError        ConnectionStatus = 4
)

func (c Connection) GetConnectionStatus() ConnectionStatus {
	return c.status
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

func (c Connection) GetTCPProcessingDuration(status ConnectionStatus) time.Duration {
	switch status {
	case ConnectionEstablished:
		return c.metrics.tcpEstablishedDuration
	case ConnectionError:
		return c.metrics.tcpErroredDuration
	default:
		return 0
	}
}
