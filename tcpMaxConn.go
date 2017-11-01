package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"github.com/dachad/check-max-tcp-connections/mtcpclient"
)

func main() {
	hostPtr := flag.String("host", "localhost", "Host you want to open tcp connections against")
	// according to https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers, you are probably not using this
	portPtr := flag.Int("port", 9998, "Port you want to open tcp connections against")
	numberConnectionsPtr := flag.Int("connections", 100, "Number of connections you want to open")
	delayPtr := flag.Int("delay", 10, "Number of ms you want to sleep between each connection creation")
	debugPtr := flag.Bool("debug", false, "Print debugging information to standard error")
	reportingIntervalPtr := flag.Int("interval", 1, "Interval, in seconds, between updating connections status")
	flag.Parse()

	var debugOut io.Writer = ioutil.Discard
	if *debugPtr {
		debugOut = os.Stderr
	}

	connStatusCh := mtcpclient.StartReportingLogic(*numberConnectionsPtr, *reportingIntervalPtr)

	mtcpclient.MultiTCPConnect(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr, connStatusCh, debugOut)

	fmt.Fprintln(debugOut, "\nTerminating Program")
}


