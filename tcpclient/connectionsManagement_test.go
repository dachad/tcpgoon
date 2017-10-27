package tcpclient

import (
	"strconv"
	"sync"
	"testing"
	"time"
	"github.com/dachad/check-max-tcp-connections/tcpserver"
)

func TestTcpConnect(t *testing.T) {
	var numberConnections = 2
	var host = "127.0.0.1"
	var port = 55555
	var delay = 1

	dispatcher := &tcpserver.Dispatcher{make(map[string]*tcpserver.Handler)}

	run := func() {
		if err := dispatcher.ListenHandlers(port); err != nil {
			t.Error("Could not start the TCP server", err)
			return
		}
		t.Log("TCP server started")
	}
	go run()

	var wg sync.WaitGroup
	wg.Add(numberConnections)

	for runner := 1; runner <= numberConnections; runner++ {
		t.Log("Initiating runner # ", strconv.Itoa(runner))
		go TcpConnect(runner, host, port, &wg)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		t.Logf("Runner %s initated. Remaining: %s", strconv.Itoa(runner), strconv.Itoa(numberConnections-runner))
	}

	t.Log("Waiting runners to finish")
	time.Sleep(time.Duration(delay) * time.Second)

	for runner := 1; runner <= numberConnections; runner++ {
		t.Log("Closing runner #", strconv.Itoa(runner))
		wg.Done()
	}

}
