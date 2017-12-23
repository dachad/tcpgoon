package mtcpclient

import (
	"fmt"
	"github.com/dachad/tcpgoon/debugging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartBackgroundClosureTrigger(gc *GroupOfConnections) <-chan bool {
	closureCh := make(chan bool)

	signalsCh := make(chan os.Signal, 1)
	registerProperSignals(signalsCh)

	go closureMonitor(gc, signalsCh, closureCh)
	return closureCh
}

func registerProperSignals(signalsCh chan os.Signal) {
	// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
	// SIGINT is the most common mechanism to the user to stop the process (Ctrl^C),
	//  but we also register SIGTERM, given its quite generic. SIGHUP could be
	//  another candidate, but in that case, probably there no terminal / user
	//  in the other end...
	signal.Notify(signalsCh, syscall.SIGINT, syscall.SIGTERM)
}

// closureMonitor polls a connections slice, to see if there's connections pending
//  to be triggered, and a signal channel, in case execution is interrupted
func closureMonitor(gc *GroupOfConnections, signalsCh chan os.Signal,
	closureCh chan bool) {
	const pullingPeriodInMs = 500
	for {
		select {
		case signal := <-signalsCh:
			fmt.Fprintln(debugging.DebugOut, "We captured a closure signal:", signal)
			close(closureCh)
			return
		case <-time.After(pullingPeriodInMs * time.Millisecond):
			if !gc.PendingConnections() {
				close(closureCh)
				return
			}
		}
	}
}
