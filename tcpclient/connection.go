package tcpclient

import "fmt"

type Connection struct {
	Id     int
	status ConnectionStatus
}

type ConnectionStatus int

var (
	ConnectionNotInitiated ConnectionStatus = 0
	ConnectionDialing      ConnectionStatus = 1
	ConnectionEstablished  ConnectionStatus = 2
	ConnectionClosed       ConnectionStatus = 3
	ConnectionError        ConnectionStatus = 4
)

func (c Connection) String() string {
	var status string
	switch c.status {
	case ConnectionNotInitiated:
		status = "not initiated"
	case ConnectionDialing:
		status = "dialing"
	case ConnectionEstablished:
		status = "established"
	case ConnectionClosed:
		status = "closed"
	case ConnectionError:
		status = "errored"
	}
	return fmt.Sprintf("Connection %d is %s", c.Id, status)
}

func PendingConnections(c []Connection) bool {
	for _, item := range c {
		if item.status == ConnectionNotInitiated || item.status == ConnectionDialing {
			return true
		}
	}
	return false
}

func PrintGroupOfConnections(c []Connection) string {
	var nDialing, nEstablished, nClosed, nNotInitiated, nError, nTotal int = 0, 0, 0, 0, 0, 0
	for _, item := range c {
		switch item.status {
		case ConnectionDialing:
			nDialing += 1
		case ConnectionEstablished:
			nEstablished += 1
		case ConnectionClosed:
			nClosed += 1
		case ConnectionError:
			nError += 1
		case ConnectionNotInitiated:
			nNotInitiated += 1
		}
		nTotal += 1
	}
	return fmt.Sprintf("Total: %d, Dialing: %d, Established: %d, Closed: %d, Error: %d, NotInitiated: %d",
		nTotal, nDialing, nEstablished, nClosed, nError, nNotInitiated)
}
