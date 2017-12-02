package cmd

import (
	"fmt"

	"github.com/dachad/tcpgoon/cmdutil"
	"github.com/dachad/tcpgoon/debugging"
	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/tcpclient"
)

func run(params tcpgoonParams) {
	tcpclient.DefaultDialTimeoutInMs = params.connDialTimeoutPtr

	// TODO: we should decouple the caller from the mtcpclient package (too many structures being moved from
	//  one side to the other.. everything in a single structure, or applying something like the builder pattern,
	//  may help
	connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(params.numberConnectionsPtr, params.reportingIntervalPtr)
	closureCh := mtcpclient.StartBackgroundClosureTrigger(connStatusTracker)
	mtcpclient.MultiTCPConnect(params.numberConnectionsPtr, params.delayPtr, params.hostPtr, params.portPtr, connStatusCh, closureCh)
	fmt.Fprintln(debugging.DebugOut, "Tests execution completed")

	cmdutil.CloseNicely(params.hostPtr, params.portPtr, connStatusTracker)
}
