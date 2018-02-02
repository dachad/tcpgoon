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

// NewConnection initializes a connection given all values that are actually stored internally. This is just being used
// as the first (dirty) approach for tests
func NewConnection(id int, status ConnectionStatus, procTime time.Duration) Connection {
	return Connection{
		ID:     id,
		status: status,
		metrics: connectionMetrics{
			tcpEstablishedDuration: procTime,
			tcpErroredDuration:     procTime,
		},
	}

}

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

func (c Connection) GetTCPProcessingDuration() time.Duration {
	if c.WentOk() {
		return c.metrics.tcpEstablishedDuration
	}
	return c.metrics.tcpErroredDuration
}

func (c Connection) IsStatusIn(statuses []ConnectionStatus) bool {
	for _, s := range statuses {
		if c.GetConnectionStatus() == s {
			return true
		}
	}
	return false
}

func (c Connection) WentOk() bool {
	if c.IsStatusIn([]ConnectionStatus{ConnectionEstablished, ConnectionClosed}) {
		return true
	}
	return false
}
