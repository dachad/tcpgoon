package tcpserver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Handler struct { // Handler : describe what this function does
	conn   net.Conn
	closed chan bool
}

func (h *Handler) Listen() { // Listen : listen connection for incomming data
	defer h.conn.Close()
	bf := bufio.NewReader(h.conn)
	for {
		line, _, err := bf.ReadLine()
		log.Println(line)
		if err != nil {
			if err == io.EOF {
				log.Println("End connection")
			}
			h.closed <- true // send to dispatcher, that connection is closed
			return
		}
	}
}

type Dispatcher struct { // Dispatcher : describe what this function does
	Handlers map[string]*Handler //`type:"map[ip]*Handler"`
}

func (d *Dispatcher) AddHandler(conn net.Conn) { // AddHandler : describe what this function does
	addr := conn.RemoteAddr().String()
	handler := &Handler{conn, make(chan bool, 1)}
	d.Handlers[addr] = handler

	go handler.Listen()

	<-handler.closed // when connection closed, remove handler from handlers
	delete(d.Handlers, addr)
}

func (d *Dispatcher) ListenHandlers(port int) error {
	sport := strconv.Itoa(port)

	ln, err := net.Listen("tcp", ":"+sport)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(conn.RemoteAddr())

		tcpconn := conn.(*net.TCPConn)
		tcpconn.SetKeepAlive(true)
		tcpconn.SetKeepAlivePeriod(10 * time.Second)

		go d.AddHandler(conn)
	}
}
