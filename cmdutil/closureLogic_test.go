package cmdutil

import (
	"io/ioutil"
	"testing"

	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/tcpclient"
)

func TestCloselyNicelyPendingConnections(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := something.NewMockMyInterface(mockCtrl)

	var host = "127.0.0.1"
	var port = 55555
	var gc = mtcpclient.GroupOfConnections{tcpclient.Connection{
		ID: 0,
	}}
	gc[0].SetConnectionStatus(tcpclient.ConnectionDialing)

	if CloseNicely(host, port, gc, ioutil.Discard) != incompleteExecutionExitStatus {
		t.Error("Some connection is still dialing")
	}
}

func TestCloselyNicelyErroredConnections(t *testing.T) {
	var host = "127.0.0.1"
	var port = 55555
	var gc = mtcpclient.GroupOfConnections{tcpclient.Connection{
		ID: 0,
	}}
	gc[0].SetConnectionStatus(tcpclient.ConnectionError)

	if CloseNicely(host, port, gc, ioutil.Discard) != completedButConnErrorsExitStatus {
		t.Error("Some connection is still errored")
	}
}

func TestCloselyNicelyClean(t *testing.T) {
	var host = "127.0.0.1"
	var port = 55555
	var gc = mtcpclient.GroupOfConnections{tcpclient.Connection{
		ID: 0,
	}}
	gc[0].SetConnectionStatus(tcpclient.ConnectionEstablished)

	if CloseNicely(host, port, gc, ioutil.Discard) != okExitStatus {
		t.Error("All connections are established")
	}
}

func TestCloseAbruptly(t *testing.T) {

	if CloseAbruptly() != incompleteExecutionExitStatus {
		t.Error("Failed abruptly closing")
	}
}
