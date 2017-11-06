package cmdutil

import (
	"os"
	"fmt"
	"io"
	"github.com/dachad/check-max-tcp-connections/tcpclient"
)

const (
	okExitStatus = 0
	incompleteExecutionExitStatus = 1
	completedButConnErrorsExitStatus = 2
)

func CloseNicely(host string, port int, connections []tcpclient.Connection, debugOut io.Writer) {
	printClosureReport(host, port, connections)
	if tcpclient.PendingConnections(connections) {
		fmt.Fprintln(debugOut, "We detected some connections did not complete")
		os.Exit(incompleteExecutionExitStatus)
	}
	if tcpclient.ConnectionInError(connections) {
		fmt.Fprintln(debugOut, "We detected connection errors")
		os.Exit(completedButConnErrorsExitStatus)
	}
	fmt.Fprintln(debugOut, "Metrics point to a clean execution. Successful exit")
	os.Exit(okExitStatus)
}

func CloseAbruptly()  {
	fmt.Println("\n*** Execution aborted as prompted by the user ***")
	os.Exit(incompleteExecutionExitStatus)
}
