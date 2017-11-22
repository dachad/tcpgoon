package mtcpclient

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dachad/tcpgoon/debugging"
	"github.com/dachad/tcpgoon/tcpclient"
)

// MultiTCPConnect tries to open us many TCP connections as numberConnections against
// host:port, with a delay between them of delay (ms). ConnStatusCh will be streaming
// tcpclient.Connection descriptions on each status update of the connections.
// closureCh will interrupt execution when closed
func MultiTCPConnect(numberConnections int, delay int, host string, port int,
	connStatusCh chan<- tcpclient.Connection, closureCh <-chan bool) {
	var wg sync.WaitGroup
	for runner := 0; runner < numberConnections; runner++ {
		select {
		case <-closureCh:
			fmt.Fprintln(debugging.DebugOut, "MultiTCPConnect routine got the closure request")
			break
		default:
			fmt.Fprintln(debugging.DebugOut, "Initiating gothread # "+strconv.Itoa(runner)+" to start a new connection")
			wg.Add(1)
			go tcpclient.TCPConnect(runner, host, port, &wg, connStatusCh, closureCh)
			fmt.Fprintln(debugging.DebugOut, "Gothread # "+strconv.Itoa(runner)+
				" initated. Remaining: "+strconv.Itoa(numberConnections-runner))
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
	fmt.Fprintln(debugging.DebugOut, "Waiting gothreads to finish")
	wg.Wait()
}
