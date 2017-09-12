package tcpserver

import (
	"testing"
	"time"
	"sync"
	"net"
)

var once sync.Once

func tcpServer(t *testing.T) func () {
	return func() {
		t.Log("Starting TCP server")
		dispatcher := &Dispatcher{make(map[string]*Handler)}		
		if err := dispatcher.ListenHandlers(8888); err != nil {
			t.Error("Could not start the TCP server", err)
			return
		} else {
			t.Log("TCP server started")
		}
	}
}
func TestTcpServer(t *testing.T) {
	t.Log("Testing TCP server")
	f := tcpServer(t)
	go once.Do(f)
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal("Coudl not connect to TCP server", err)
	}
	conn.Close()}
