package main

import (
	"fmt"
	"os"

	"github.com/dachad/tcpgoon/cmdutil"
	"github.com/dachad/tcpgoon/debugging"
	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/tcpclient"
	"github.com/spf13/pflag"
)

func main() {
	hostPtr := pflag.StringP("host", "h", "", "Host you want to open tcp connections against (Required)")
	portPtr := pflag.IntP("port", "p", 0, "Port you want to open tcp connections against (Required)")
	numberConnectionsPtr := pflag.IntP("connections", "c", 100, "Number of connections you want to open")
	delayPtr := pflag.IntP("sleep", "s", 10, "Time you want to sleep between connections, in ms")
	connDialTimeoutPtr := pflag.IntP("dial-timeout", "t", 5000, "Connection dialing timeout, in ms")
	debugPtr := pflag.BoolP("debug", "d", false, "Print debugging information to the standard error")
	reportingIntervalPtr := pflag.IntP("interval", "i", 1, "Interval, in seconds, between stats updates")
	assumeyesPtr := pflag.BoolP("assume-yes", "y", false, "Force execution without asking for confirmation")
	pflag.Parse()

	// Target host and port are mandatory to run the TCP check
	if *hostPtr == "" || *portPtr == 0 {
		pflag.PrintDefaults()
		os.Exit(1)
	}

	tcpclient.DefaultDialTimeoutInMs = *connDialTimeoutPtr

	if *debugPtr {
		debugging.EnableDebug()
	}

	if !(*assumeyesPtr || cmdutil.AskForUserConfirmation(*hostPtr, *portPtr, *numberConnectionsPtr)) {
		fmt.Fprintln(debugging.DebugOut, "Execution not approved by the user")
		cmdutil.CloseAbruptly()
	}

	// TODO: we should decouple the caller from the mtcpclient package (too many structures being moved from
	//  one side to the other.. everything in a single structure, or applying something like the builder pattern,
	//  may help
	connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(*numberConnectionsPtr, *reportingIntervalPtr)
	closureCh := mtcpclient.StartBackgroundClosureTrigger(connStatusTracker)
	mtcpclient.MultiTCPConnect(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr, connStatusCh, closureCh)
	fmt.Fprintln(debugging.DebugOut, "Tests execution completed")

	cmdutil.CloseNicely(*hostPtr, *portPtr, connStatusTracker)
}
