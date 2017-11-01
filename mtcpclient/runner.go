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
// tcpclient.Connection descriptions on each status update of the connections
func MultiTCPConnect(numberConnections int, delay int, host string, port int,
		connStatusCh chan tcpclient.Connection, debugOut io.Writer) {
	var wg sync.WaitGroup
	wg.Add(numberConnections)
	for runner := 0; runner < numberConnections; runner++ {
		fmt.Fprintln(debugOut, "Initiating runner # " + strconv.Itoa(runner))
		go tcpclient.TCPConnect(runner, host, port, &wg, debugOut, connStatusCh)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Fprintln(debugOut, "Runner " + strconv.Itoa(runner) +
			" initated. Remaining: " + strconv.Itoa(numberConnections - runner))
	}
	fmt.Fprintln(debugOut, "Waiting runners to finish")
	wg.Wait()
}
