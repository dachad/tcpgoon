package mtcpclient

import (
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"time"
)

func StartBackgroundClosureTrigger(connections []tcpclient.Connection) (<-chan bool){
	closureCh := make(chan bool)
	go closureMonitor(connections, closureCh)
	return closureCh
}
func closureMonitor(connections []tcpclient.Connection, closureCh chan bool) {
	const pullingPeriodInMs  = 500
	for  {
		if !tcpclient.PendingConnections(connections) {
			close(closureCh)
			return
		}
		time.Sleep(pullingPeriodInMs * time.Millisecond)
	}
}
