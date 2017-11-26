package cmd

import (
	"fmt"

	"github.com/dachad/tcpgoon/cmdutil"
	"github.com/dachad/tcpgoon/debugging"
	"github.com/dachad/tcpgoon/mtcpclient"
	"github.com/dachad/tcpgoon/tcpclient"
)

func run(flags tcpgoonFlags) {
	tcpclient.DefaultDialTimeoutInMs = flags.connDialTimeoutPtr

	// TODO: we should decouple the caller from the mtcpclient package (too many structures being moved from
	//  one side to the other.. everything in a single structure, or applying something like the builder pattern,
	//  may help
	connStatusCh, connStatusTracker := mtcpclient.StartBackgroundReporting(flags.numberConnectionsPtr, flags.reportingIntervalPtr)
	closureCh := mtcpclient.StartBackgroundClosureTrigger(connStatusTracker)
	mtcpclient.MultiTCPConnect(flags.numberConnectionsPtr, flags.delayPtr, flags.hostPtr, flags.portPtr, connStatusCh, closureCh)
	fmt.Fprintln(debugging.DebugOut, "Tests execution completed")

	cmdutil.CloseNicely(flags.hostPtr, flags.portPtr, connStatusTracker)
}
