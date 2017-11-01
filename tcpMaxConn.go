package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"github.com/dachad/check-max-tcp-connections/mtcpclient"
	"github.com/spf13/pflag"
)

func main() {
	hostPtr := pflag.StringP("host", "h", "localhost", "Host you want to open tcp connections against")
	// according to https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers, you are probably not using this
	portPtr := pflag.IntP("port", "p", 9998, "Port you want to open tcp connections against")
	numberConnectionsPtr := pflag.IntP("connections", "c", 100, "Number of connections you want to open")
	delayPtr := pflag.IntP("sleep", "s", 10, "Time you want to sleep between connections, in ms")
	debugPtr := pflag.BoolP("debug", "d", false, "Print debugging information to the standard error")
	reportingIntervalPtr := pflag.IntP("interval", "i", 1, "Interval, in seconds, between stats updates")
	pflag.Parse()

	var debugOut io.Writer = ioutil.Discard
	if *debugPtr {
		debugOut = os.Stderr
	}

	connStatusCh := mtcpclient.StartReportingLogic(*numberConnectionsPtr, *reportingIntervalPtr)

	mtcpclient.MultiTCPConnect(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr, connStatusCh, debugOut)

	fmt.Fprintln(debugOut, "\nTerminating Program")
}


