package main

import (
	"flag"
	"fmt"
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var debugOut io.Writer = ioutil.Discard

func runTCPConnectionsInParallel(numberConnections int, delay int, host string, port int, rinterval int) {
	const numberOfGoRoutinesToCollectAndReportStatus = 2
	runtime.GOMAXPROCS(numberConnections + numberOfGoRoutinesToCollectAndReportStatus)
	// A connection may report up to 3 messages: Dialing -> Established -> Connected
	const MaxMessagesWeMayGetPerConnection = 3
	connStatusCh := make(chan tcpclient.Connection, numberConnections*MaxMessagesWeMayGetPerConnection)
	connStatusTracker := make([]tcpclient.Connection, numberConnections)

	// moving these outside of runTcpCOnnectionsInParallel may have a lot of sense.. now this is
	//  doing too much staff. We should also think about moving to the client package..
	go reportConnectionsStatus(connStatusTracker, rinterval)
	go collectConnectionsStatus(connStatusTracker, connStatusCh)

	var wg sync.WaitGroup
	wg.Add(numberConnections)
	for runner := 0; runner < numberConnections; runner++ {
		fmt.Fprintln(debugOut, "Initiating runner # "+strconv.Itoa(runner))
		go tcpclient.TCPConnect(runner, host, port, &wg, debugOut, connStatusCh)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Fprintln(debugOut, "Runner "+strconv.Itoa(runner)+
			" initated. Remaining: "+strconv.Itoa(numberConnections-runner))
	}
	fmt.Fprintln(debugOut, "Waiting runners to finish")
	wg.Wait()
}
func collectConnectionsStatus(connectionDescriptions []tcpclient.Connection, statusChannel <-chan tcpclient.Connection) {
	for {
		connectionStatus := <-statusChannel
		connectionDescriptions[connectionStatus.Id] = connectionStatus

	}
}
func reportConnectionsStatus(connectionDescriptions []tcpclient.Connection, intervalBetweenUpdates int) {
	for {
		fmt.Println(tcpclient.PrintGroupOfConnections(connectionDescriptions))
		time.Sleep(time.Duration(intervalBetweenUpdates) * time.Second)
	}
}

func main() {
	hostPtr := flag.String("host", "localhost", "Host you want to open tcp connections against")
	// according to https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers, you are probably not using this
	portPtr := flag.Int("port", 9998, "Port you want to open tcp connections against")
	numberConnectionsPtr := flag.Int("connections", 100, "Number of connections you want to open")
	delayPtr := flag.Int("delay", 10, "Number of ms you want to sleep between each connection creation")
	debugPtr := flag.Bool("debug", false, "Print debugging information to standard error")
	reportingIntervalPtr := flag.Int("interval", 1, "Interval, in seconds, between updating connections status")
	flag.Parse()
	if *debugPtr {
		debugOut = os.Stderr
	}

	runTCPConnectionsInParallel(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr, *reportingIntervalPtr)

	fmt.Fprintln(debugOut, "\nTerminating Program")
}
