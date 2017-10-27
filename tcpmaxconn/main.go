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
)

func connectionHandler(id int, host string, port int, wg *sync.WaitGroup) error {
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

func runThreads(numberConnections int, delay int, host string, port int) {
	runtime.GOMAXPROCS(numberConnections)

	var wg sync.WaitGroup
	wg.Add(numberConnections)

	for runner := 1; runner <= numberConnections; runner++ {
		fmt.Println("Initiating runner # " + strconv.Itoa(runner))
		go connectionHandler(runner, host, port, &wg)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Println("Runner " + strconv.Itoa(runner) + " initated. Remaining: " + strconv.Itoa(numberConnections-runner))
	}

	fmt.Println("Waiting runners to finish")
	wg.Wait()
}

var versionString = "unset"

func main() {
	hostPtr := flag.String("host", "localhost", "Host you want to open tcp connections against")
	portPtr := flag.Int("port", 8888, "Port you want to open tcp connections against")
	numberConnectionsPtr := flag.Int("connections", 100, "Number of connections you want to open")
	delayPtr := flag.Int("delay", 10, "Number of ms you want to sleep between each connection creation")

	flag.Parse()

	fmt.Println("Running Version: ", versionString)

	runThreads(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr)

	fmt.Println("\nTerminating Program")
}
