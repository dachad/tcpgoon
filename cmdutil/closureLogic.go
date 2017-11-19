package cmdutil

import (
	"fmt"
	"io"

	"github.com/dachad/tcpgoon/mtcpclient"
)

const (
	okExitStatus                     = 0
	incompleteExecutionExitStatus    = 1
	completedButConnErrorsExitStatus = 2
)

func CloseNicely(host string, port int, gc mtcpclient.GroupOfConnections, debugOut io.Writer) int {
	printClosureReport(host, port, gc)
	if gc.PendingConnections() {
		fmt.Fprintln(debugOut, "We detected some connections did not complete")
		return incompleteExecutionExitStatus
	}
	if gc.AtLeastOneConnectionInError() {
		fmt.Fprintln(debugOut, "We detected connection errors")
		return completedButConnErrorsExitStatus
	}
	fmt.Fprintln(debugOut, "Metrics point to a clean execution. Successful exit")
	return okExitStatus
}

func CloseAbruptly() int {
	fmt.Println("\n*** Execution aborted as prompted by the user ***")
	return incompleteExecutionExitStatus
}
