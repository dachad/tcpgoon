package tcpclient

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
	"reflect"
)

// this first one can be a flag in the future
var DefaultDialTimeoutInSecs = 5
var defaultReadTimeoutInSecs = 1

func reportConnectionStatus(debugOut io.Writer, statusChannel chan<- Connection, connectionDescription Connection) {
	statusChannel <- connectionDescription
	fmt.Fprintln(debugOut, "\t", connectionDescription)
}

// TCPConnect just opens a TCP connection against the target described by
// the host:port, and considers the id to report back status changes through the
// status goChannel with descriptors matching the Connection struct supplied in this
// same package. It also needs an iowriter to print debugging information (it can be
// ioutil.Discard.
func TCPConnect(id int, host string, port int, wg *sync.WaitGroup, debugOut io.Writer,
		statusChannel chan<- Connection, closeRequest <-chan bool) error {
	connectionDescription := Connection{
		Id:     id,
		status: ConnectionDialing,
	}
	reportConnectionStatus(debugOut, statusChannel, connectionDescription)
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port),
		time.Duration(DefaultDialTimeoutInSecs)*time.Second)

	if err != nil {
		connectionDescription.status = ConnectionError
		reportConnectionStatus(debugOut, statusChannel, connectionDescription)
		fmt.Fprintln(debugOut, "Connection", id, "was unable to open the connection. Error:")
		fmt.Fprintln(debugOut, err)
		wg.Done()
		return err
	}

	connectionDescription.status = ConnectionEstablished
	reportConnectionStatus(debugOut, statusChannel, connectionDescription)
	connBuf := bufio.NewReader(conn)
	for {
		select {
		case <-closeRequest:
			// we need to close gracefully the connection
		default:
			// TODO: to review
			conn.SetReadDeadline(time.Now().Add(time.Duration(defaultReadTimeoutInSecs) * time.Second))
			str, err := connBuf.ReadString('\n')
			if len(str) > 0 {
				fmt.Fprintln(debugOut, "Connection", id, "got", str)
			}
			if err, ok := err.(net.Error); ok && err.Timeout() {
				fmt.Fprintln(debugOut, "Connection", id, "timed out reading. Trying again..")
			} else if err != nil {
				fmt.Fprintln(debugOut, "Connection", id, "looks closed. Error", reflect.TypeOf(err),"when reading:")
				fmt.Fprintln(debugOut, err)
				connectionDescription.status = ConnectionClosed
				reportConnectionStatus(debugOut, statusChannel, connectionDescription)
				wg.Done()
				return err
			}
		}

	}
}
