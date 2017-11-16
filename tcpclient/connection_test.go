package tcpclient

import (
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
