package tcpclient

import (
	"sync"
	"fmt"
	"net"
	"bufio"
)

func TcpConnect(id int, host string, port int, wg *sync.WaitGroup) error {
	fmt.Println("\t runner " + strconv.Itoa(id) + " is initiating a connection")
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return err
	}
	fmt.Println("\t runner " + strconv.Itoa(id) + " established the connection")
	connBuf := bufio.NewReader(conn)
	for {
		str, err := connBuf.ReadString('\n')
		if len(str) > 0 {
			fmt.Println(str)
		}
		if err != nil {
			fmt.Println("\t runner " + strconv.Itoa(id) + " got its connection closed")
			wg.Done()
			return err
		}
	}
}
