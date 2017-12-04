package tcpclient

import (
	"strings"
	"testing"
	"time"
)

func TestGetConnectionStatus(t *testing.T) {
	connectionDescription := Connection{
		ID:      0,
		status:  ConnectionDialing,
		metrics: connectionMetrics{},
	}
	if connectionDescription.GetConnectionStatus() != ConnectionDialing {
		t.Error("GetConnectionStatus is not working")
	}
}

func TestGetTCPProcessingDurationEstablished(t *testing.T) {

	connectionDescription := Connection{
		ID:      0,
		status:  ConnectionEstablished,
		metrics: connectionMetrics{tcpEstablishedDuration: 1 * time.Second},
	}
	if connectionDescription.GetTCPProcessingDuration(ConnectionEstablished) != 1*time.Second {
		t.Error("GetTCPProcessingDuration is not working")
	}
}

func TestGetTCPProcessingDurationErrored(t *testing.T) {

	connectionDescription := Connection{
		ID:      0,
		status:  ConnectionError,
		metrics: connectionMetrics{tcpErroredDuration: 1 * time.Second},
	}
	if connectionDescription.GetTCPProcessingDuration(ConnectionError) != 1*time.Second {
		t.Error("GetTCPProcessingDuration is not working")
	}
}

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
