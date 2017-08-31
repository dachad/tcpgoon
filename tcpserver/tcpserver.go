package main
import (
    "strconv"
    "net"
    "log"
    "time"
    "bufio"
	"io"
	"fmt"
)

type Handler struct {
    conn   net.Conn
    closed chan bool
}

func (h *Handler) Listen() { // listen connection for incomming data
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

        // ... some business logic with data
    }
}

type Dispatcher struct {
    handlers map[string]*Handler `map[ip]*Handler`
}

func (d *Dispatcher) AddHandler(conn net.Conn) {
    addr := conn.RemoteAddr().String()
    handler := &Handler{conn, make(chan bool, 1)}
    d.handlers[addr] = handler

    go handler.Listen()

    <-handler.closed // when connection closed, remove handler from handlers
    delete(d.handlers, addr)
}

func (d *Dispatcher) ListenHandlers(port int) {
    sport := strconv.Itoa(port)

    ln, err := net.Listen("tcp", ":" + sport)
    if err != nil {
        log.Println(err)
        return
    }

    defer ln.Close()

    for {
        conn, err := ln.Accept() // accept connection
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

func main() {
    dispatcher := &Dispatcher{make(map[string]*Handler)}
    dispatcher.ListenHandlers(8888)
}
