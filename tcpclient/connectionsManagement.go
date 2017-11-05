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
		status: connectionDialing,
	}
	reportConnectionStatus(debugOut, statusChannel, connectionDescription)
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port),
		time.Duration(DefaultDialTimeoutInSecs)*time.Second)
	if err != nil {
		connectionDescription.status = connectionError
		reportConnectionStatus(debugOut, statusChannel, connectionDescription)
		fmt.Fprintln(debugOut, "Connection", id, "was unable to open the connection. Error:")
		fmt.Fprintln(debugOut, err)
		wg.Done()
		return err
	}

	defer conn.Close()
	connectionDescription.status = connectionEstablished
	reportConnectionStatus(debugOut, statusChannel, connectionDescription)
	connBuf := bufio.NewReader(conn)
	for {
		select {
		case <-closeRequest:
			fmt.Fprintln(debugOut, "Connection", id, "is being requested to close")
			// we don't mark connection as closed, as its us closing cleanly at the end of the execution,
			//  so final report can consider it was established when finishing and not closed by the other end
			wg.Done()
			return nil
		default:
			conn.SetReadDeadline(time.Now().Add(time.Duration(defaultReadTimeoutInSecs) * time.Second))
			str, err := connBuf.ReadString('\n')
			if len(str) > 0 {
				fmt.Fprintln(debugOut, "Connection", id, "got", str)
			} else if terr, ok := err.(net.Error); ok && terr.Timeout() {
				fmt.Fprintln(debugOut, "No info from connection", id, "before timing out. Reading again..")
			} else if err != nil {
				fmt.Fprintln(debugOut, "Connection", id, "looks closed. Error", reflect.TypeOf(err),"when reading:")
				fmt.Fprintln(debugOut, err)
				connectionDescription.status = connectionClosed
				reportConnectionStatus(debugOut, statusChannel, connectionDescription)
				wg.Done()
				return err
			}
		}

	}
}
