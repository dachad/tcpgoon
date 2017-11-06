package mtcpclient

import (
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"sync"
	"fmt"
	"time"
	"io"
	"strconv"
)

// MultiTCPConnect tries to open us many TCP connections as numberConnections against
// host:port, with a delay between them of delay (ms). You can supply a
// debugOut to get debugging messages, while connStatusCh will be streaming
// tcpclient.Connection descriptions on each status update of the connections.
// closureCh will interrupt execution when closed
func MultiTCPConnect(numberConnections int, delay int, host string, port int,
connStatusCh chan <- tcpclient.Connection, closureCh <-chan bool,
debugOut io.Writer) {
	var wg sync.WaitGroup
	for runner := 0; runner < numberConnections; runner++ {
		select {
		case <-closureCh:
			fmt.Fprintln(debugOut, "MultiTCPConnect routine got the closure request")
			break
		default:
			fmt.Fprintln(debugOut, "Initiating gothread # " + strconv.Itoa(runner) + " to start a new connection")
			wg.Add(1)
			go tcpclient.TCPConnect(runner, host, port, &wg, debugOut, connStatusCh, closureCh)
			fmt.Fprintln(debugOut, "Gothread # " + strconv.Itoa(runner) +
				" initated. Remaining: " + strconv.Itoa(numberConnections - runner))
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
	fmt.Fprintln(debugOut, "Waiting gothreads to finish")
	wg.Wait()
}
