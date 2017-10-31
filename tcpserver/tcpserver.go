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

// Handler : Struct
type Handler struct {
	conn   net.Conn
	closed chan bool
}

// Listen : listen connection for incomming data
func (h *Handler) listen() {
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

// Dispatcher : Struct with all the hanlders
type Dispatcher struct {
	Handlers map[string]*Handler //`type:"map[ip]*Handler"`
}

func (d *Dispatcher) addHandler(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	handler := &Handler{conn, make(chan bool, 1)}
	d.Handlers[addr] = handler

	go handler.listen()

	<-handler.closed // when connection closed, remove handler from handlers
	delete(d.Handlers, addr)
}

// ListenHandlers : start listening on the handler
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

		go d.addHandler(conn)
	}
}
