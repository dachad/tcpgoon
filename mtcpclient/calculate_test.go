package mtcpclient

import (
	"testing"
	"github.com/dachad/tcpgoon/tcpclient"
	"fmt"
)

func TestcalculateStdDev(t *testing.T) {
	var gc GroupOfConnections
	gc = []tcpclient.Connection{
		{},
	}
	fmt.Println(gc.calculateStdDev(tcpclient.ConnectionNotInitiated, metricsCollectionStats{}))
}

