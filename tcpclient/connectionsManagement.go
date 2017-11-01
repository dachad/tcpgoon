package tcpclient

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

// this can be an argument in the future
var defaultDialTimeoutInSecs = 10

func reportConnectionStatus(debugOut io.Writer, statusChannel chan<- Connection, connectionDescription Connection) {
	statusChannel <- connectionDescription
	fmt.Fprintln(debugOut, "\t", connectionDescription)
}

// TCPConnect just opens a TCP connection against the target described by
// the host:port, and considers the id to report back status changes through the
// status goChannel with descriptors matching the Connection struct supplied in this
// same package.
// It also needs an iowriter to print debugging information.
func TCPConnect(id int, host string, port int, wg *sync.WaitGroup, debugOut io.Writer, statusChannel chan<- Connection) error {
	connectionDescription := Connection{
		Id:     id,
		status: ConnectionDialing,
	}
	reportConnectionStatus(debugOut, statusChannel, connectionDescription)
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port),
		time.Duration(defaultDialTimeoutInSecs)*time.Second)
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
		str, err := connBuf.ReadString('\n')
		if len(str) > 0 {
			fmt.Fprintln(debugOut, "Connection", id, "got", str)
		}
		if err != nil {
			connectionDescription.status = ConnectionClosed
			reportConnectionStatus(debugOut, statusChannel, connectionDescription)
			wg.Done()
			return err
		}
	}
}
