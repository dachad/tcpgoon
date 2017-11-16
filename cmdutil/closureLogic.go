package cmdutil

import (
	"fmt"
	"io"
	"os"

	"github.com/dachad/tcpgoon/mtcpclient"
)

const (
	okExitStatus                     = 0
	incompleteExecutionExitStatus    = 1
	completedButConnErrorsExitStatus = 2
)

func CloseNicely(host string, port int, gc mtcpclient.GroupOfConnections, debugOut io.Writer) {
	printClosureReport(host, port, gc)
	if gc.PendingConnections() {
		fmt.Fprintln(debugOut, "We detected some connections did not complete")
		os.Exit(incompleteExecutionExitStatus)
	}
	if gc.AtLeastOneConnectionInError() {
		fmt.Fprintln(debugOut, "We detected connection errors")
		os.Exit(completedButConnErrorsExitStatus)
	}
	fmt.Fprintln(debugOut, "Metrics point to a clean execution. Successful exit")
	os.Exit(okExitStatus)
}

func CloseAbruptly() {
	fmt.Println("\n*** Execution aborted as prompted by the user ***")
	os.Exit(incompleteExecutionExitStatus)
}
