package tcpserver

import (
	"net"
	"sync"
	"testing"
	"time"
)

var once sync.Once

func tcpServer(t *testing.T) func() {
	return func() {
		t.Log("Starting TCP server")
		dispatcher := &Dispatcher{
			Handlers: make(map[string]*Handler),
			Lock: sync.RWMutex{},
		}
		if err := dispatcher.ListenHandlers(8888); err != nil {
			t.Error("Could not start the TCP server", err)
			return
		}
		t.Log("TCP server started")
	}
}

func TestTcpServer(t *testing.T) {
	t.Log("Testing TCP server")
	f := tcpServer(t)
	go once.Do(f)
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal("Could not connect to TCP server", err)
	}
	conn.Close()
}
