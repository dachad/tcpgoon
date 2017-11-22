package cmdutil

import (
	"fmt"
	"os"

	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/debugging"
)

const (
	okExitStatus                     = 0
	incompleteExecutionExitStatus    = 1
	completedButConnErrorsExitStatus = 2
)

func CloseNicely(host string, port int, gc mtcpclient.GroupOfConnections) {
	printClosureReport(host, port, gc)
	if gc.PendingConnections() {
		fmt.Fprintln(debugging.DebugOut, "We detected some connections did not complete")
		os.Exit(incompleteExecutionExitStatus)
	}
	if gc.AtLeastOneConnectionInError() {
		fmt.Fprintln(debugging.DebugOut, "We detected connection errors")
		os.Exit(completedButConnErrorsExitStatus)
	}
	fmt.Fprintln(debugging.DebugOut, "Metrics point to a clean execution. Successful exit")
	os.Exit(okExitStatus)
}

func CloseAbruptly() {
	fmt.Println("\n*** Execution aborted as prompted by the user ***")
	os.Exit(incompleteExecutionExitStatus)
}
