package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"github.com/dachad/check-max-tcp-connections/mtcpclient"
	"github.com/spf13/pflag"
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"time"
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

	if *assumeyesPtr || askForUserConfirmation(*hostPtr, *portPtr, *numberConnectionsPtr) {
		connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(*numberConnectionsPtr, *reportingIntervalPtr)
		closureCh := mtcpclient.StartBackgroundClosureTrigger(connStatusTracker)
		mtcpclient.MultiTCPConnect(*numberConnectionsPtr, *delayPtr, *hostPtr, *portPtr, connStatusCh, closureCh, debugOut)

		printClosureReport(*hostPtr, *portPtr, connStatusTracker)
		fmt.Fprintln(debugOut, "\nTerminating Program")
	} else {
		fmt.Println("\n*** Execution aborted as prompted by the user ***")
	}
}
func printClosureReport(host string, port int, connections []tcpclient.Connection) {
	// workaround to allow last status updates to be collected properly
	time.Sleep(time.Duration(50) * time.Millisecond)
	fmt.Println(strings.Repeat("-", 3), host + ":" + string(port), "tcp test statistics", strings.Repeat("-", 3))
	mtcpclient.ReportConnectionsStatus(connections, 0)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func askForUserConfirmation(host string, port int, connections int) bool {

	fmt.Println("****************************** WARNING ******************************")
	fmt.Println("* You are going to  run a TCP stress check with these arguments:")
	fmt.Println("*	- Target: " + host)
	fmt.Println("*	- TCP Port: " + strconv.Itoa(port))
	fmt.Println("*	- # of concurrent connections: " + strconv.Itoa(connections))
	fmt.Println("*********************************************************************")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to continue? (y/N): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Response not processed")
			os.Exit(1)
		}

		response = strings.TrimSuffix(response, "\n")
		response = strings.ToLower(response)
		switch {
		case stringInSlice(response, []string{"yes", "y"}):
			return true
		case stringInSlice(response, []string{"no", "n", ""}):
			return false
		default:
			fmt.Println("\nSorry, response not recongized. Try again, please")
		}
	}
}
