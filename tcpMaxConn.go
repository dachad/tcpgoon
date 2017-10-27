package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
	"github.com/dachad/check-max-tcp-connections/tcpclient"
)

func runTcpConnectionsInParallel(numberConnections int, delay int, host string, port int) {
	runtime.GOMAXPROCS(numberConnections)

	var wg sync.WaitGroup
	wg.Add(numberConnections)

	for runner := 1; runner <= numberConnections; runner++ {
		fmt.Println("Initiating runner # " + strconv.Itoa(runner))
		go tcpclient.TcpConnect(runner, host, port, &wg)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Println("Runner " + strconv.Itoa(runner) +
			" initated. Remaining: " + strconv.Itoa(numberConnections-runner))
	}

	fmt.Println("Waiting runners to finish")
	wg.Wait()
}

func main() {
	hostPtr := flag.String("host", "localhost", "Host you want to open tcp connections against")
	// according to https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers, you are probably not using this
	portPtr := flag.Int("port", 9998, "Port you want to open tcp connections against")
	numberConnectionsPtr := flag.Int("connections", 100, "Number of connections you want to open")
	delayPtr := flag.Int("delay", 10, "Number of ms you want to sleep between each connection creation")

	flag.Parse()

	runTcpConnectionsInParallel(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr)

	fmt.Println("\nTerminating Program")
}
