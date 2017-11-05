package tcpclient

import "fmt"

type Connection struct {
	Id     int
	status connectionStatus
}

type connectionStatus int

const (
	connectionNotInitiated connectionStatus = 0
	connectionDialing connectionStatus = 1
	connectionEstablished connectionStatus = 2
	connectionClosed connectionStatus = 3
	connectionError connectionStatus = 4
)

func (cs connectionStatus) isIn(connections []Connection) bool {
	for _, item := range connections {
		if item.status == cs {
			return true
		}
	}
	return false
}

func (c Connection) String() string {
	var status string
	switch c.status {
	case connectionNotInitiated:
		status = "not initiated"
	case connectionDialing:
		status = "dialing"
	case connectionEstablished:
		status = "established"
	case connectionClosed:
		status = "closed"
	case connectionError:
		status = "errored"
	}
	return fmt.Sprintf("Connection %d is %s", c.Id, status)
}

func PendingConnections(c []Connection) bool {
	return connectionNotInitiated.isIn(c) || connectionDialing.isIn(c)
}

func ConnectionInError(c []Connection) bool {
	return connectionError.isIn(c)
}

func PrintGroupOfConnections(c []Connection) string {
	var nDialing, nEstablished, nClosed, nNotInitiated, nError, nTotal int = 0, 0, 0, 0, 0, 0
	for _, item := range c {
		switch item.status {
		case connectionDialing:
			nDialing++
		case connectionEstablished:
			nEstablished++
		case connectionClosed:
			nClosed++
		case connectionError:
			nError++
		case connectionNotInitiated:
			nNotInitiated++
		}
		nTotal++
	}
	return fmt.Sprintf("Total: %d, Dialing: %d, Established: %d, Closed: %d, Error: %d, NotInitiated: %d",
		nTotal, nDialing, nEstablished, nClosed, nError, nNotInitiated)
}
