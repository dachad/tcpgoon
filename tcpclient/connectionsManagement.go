package tcpclient

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/dachad/tcpgoon/debugging"
)

var DefaultDialTimeoutInMs = 5000

func reportConnectionStatus(statusChannel chan<- Connection, connectionDescription Connection) {
	statusChannel <- connectionDescription
	fmt.Fprintln(debugging.DebugOut, "\t", connectionDescription)
}

// TCPConnect just opens a TCP connection against the target described by
// the host:port, and considers the id to report back status changes through the
// status goChannel with descriptors matching the Connection struct supplied in this
// same package.
func TCPConnect(id int, host string, port int, wg *sync.WaitGroup,
	statusChannel chan<- Connection, closeRequest <-chan bool) error {
	connectionDescription := Connection{
		ID:      id,
		status:  ConnectionDialing,
		metrics: connectionMetrics{},
	}
	reportConnectionStatus(statusChannel, connectionDescription)
	timeTCPInitiatied := time.Now()
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port),
		time.Duration(DefaultDialTimeoutInMs)*time.Millisecond)
	if err != nil {
		connectionDescription.status = ConnectionError
		connectionDescription.metrics.tcpErroredDuration = time.Now().Sub(timeTCPInitiatied)
		reportConnectionStatus(statusChannel, connectionDescription)
		fmt.Fprintln(debugging.DebugOut, "Connection", id, "was unable to open the connection. Error:")
		fmt.Fprintln(debugging.DebugOut, err)
		wg.Done()
		return err
	}
	defer conn.Close()
	connectionDescription.status = ConnectionEstablished
	connectionDescription.metrics.tcpEstablishedDuration = time.Now().Sub(timeTCPInitiatied)
	reportConnectionStatus(statusChannel, connectionDescription)
	connBuf := bufio.NewReader(conn)
	for {
		select {
		case <-closeRequest:
			fmt.Fprintln(debugging.DebugOut, "Connection", id, "is being requested to close")
			// we don't mark connection as closed, as its us closing cleanly at the end of the execution,
			//  so final report can consider it was established when finishing and not closed by the other end
			wg.Done()
			return nil
		default:
			const ReadTimeoutAndBetweenPollsInMs = 1000
			conn.SetReadDeadline(time.Now().Add(time.Duration(ReadTimeoutAndBetweenPollsInMs) * time.Millisecond))
			str, err := connBuf.ReadString('\n')
			if terr, ok := err.(net.Error); ok && terr.Timeout() {
				fmt.Fprintln(debugging.DebugOut, "No info from connection", id, "before timing out. Reading again...")
			} else if err != nil {
				fmt.Fprintln(debugging.DebugOut, "Connection", id, "looks closed. Error", reflect.TypeOf(err), "when reading:")
				fmt.Fprintln(debugging.DebugOut, err)
				connectionDescription.status = ConnectionClosed
				reportConnectionStatus(statusChannel, connectionDescription)
				wg.Done()
				return err
			} else if len(str) > 0 {
				fmt.Fprintln(debugging.DebugOut, "Connection", id, "got", str)
			}
		}

	}
}
