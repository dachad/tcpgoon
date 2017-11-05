package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"github.com/spf13/pflag"
	"github.com/dachad/check-max-tcp-connections/mtcpclient"
	"github.com/dachad/check-max-tcp-connections/cmdutil"
)



func main() {
	hostPtr := pflag.StringP("host", "h", "", "Host you want to open tcp connections against (Required)")
	// according to https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers, you are probably not using this
	portPtr := pflag.IntP("port", "p", 0, "Port you want to open tcp connections against (Required)")
	numberConnectionsPtr := pflag.IntP("connections", "c", 100, "Number of connections you want to open")
	delayPtr := pflag.IntP("sleep", "s", 10, "Time you want to sleep between connections, in ms")
	debugPtr := pflag.BoolP("debug", "d", false, "Print debugging information to the standard error")
	reportingIntervalPtr := pflag.IntP("interval", "i", 1, "Interval, in seconds, between stats updates")
	assumeyesPtr := pflag.BoolP("assume-yes", "y", false, "Force execution without asking for confirmation")
	pflag.Parse()

	// Target host and port are mandatory to run the TCP check
	if *hostPtr == "" || *portPtr == 0 {
		pflag.PrintDefaults()
		os.Exit(1)
	}

	var debugOut io.Writer = ioutil.Discard
	if *debugPtr {
		debugOut = os.Stderr
	}

	if !(*assumeyesPtr || cmdutil.AskForUserConfirmation(*hostPtr, *portPtr, *numberConnectionsPtr)) {
		fmt.Fprintln(debugOut, "Execution not approved by the user")
		cmdutil.CloseAbruptly()
	}

	connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(*numberConnectionsPtr, *reportingIntervalPtr)
	closureCh := mtcpclient.StartBackgroundClosureTrigger(connStatusTracker)
	mtcpclient.MultiTCPConnect(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr, connStatusCh, closureCh, debugOut)
	fmt.Fprintln(debugOut, "Tests execution completed")

	cmdutil.CloseNicely(*hostPtr, *portPtr, connStatusTracker, debugOut)
}
