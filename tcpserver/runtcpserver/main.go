package main

import (
    "github.com/dachad/check-max-tcp-connections/tcpserver"
)

func main() {
    dispatcher := &tcpserver.Dispatcher{make(map[string]*tcpserver.Handler)}
    dispatcher.ListenHandlers(8888)
}