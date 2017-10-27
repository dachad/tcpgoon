package tcpclient

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

func reportConnectionStatus(debugOut io.Writer, statusChannel chan<- Connection, connectionDescription Connection)  {
	statusChannel <- connectionDescription
	fmt.Fprintln(debugOut, "\t", connectionDescription)
}

func TcpConnect(id int, host string, port int, wg *sync.WaitGroup, debugOut io.Writer, statusChannel chan<- Connection) error {
	connectionDescription := Connection{
		Id: id,
		status: ConnectionDialing,
	}
	reportConnectionStatus(debugOut, statusChannel, connectionDescription)
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
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
