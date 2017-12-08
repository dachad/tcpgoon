package tcpclient

import (
	"strings"
	"testing"
	"time"
)

func TestConnectionString(t *testing.T) {

	connectionDescription := Connection{
		ID:      0,
		status:  ConnectionEstablished,
		metrics: connectionMetrics{tcpEstablishedDuration: 1 * time.Second},
	}
	if strings.Compare(connectionDescription.String(), "Connection 0 has become established after 1s") != 0 {
		t.Error("Connection string interface is not as expected")
	}
}
