package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/coder/websocket"
)

type broadcastServer struct {
	connections map[*websocket.Conn]struct{}
	join        chan *websocket.Conn
	leave       chan *websocket.Conn
	message     chan []byte
	ctx         context.Context
}

func NewBroadcastServer(ctx context.Context) *broadcastServer {
	return &broadcastServer{
		connections: map[*websocket.Conn]struct{}{},
		join:        make(chan *websocket.Conn),
		leave:       make(chan *websocket.Conn),
		message:     make(chan []byte),
		ctx:         ctx,
	}
}

func (bs *broadcastServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil) // don't require custom AcceptOptions
	if err != nil {
		log.Printf("upgrade HTTP conn to websocket: %v", err)
		return
	}
	defer conn.CloseNow()

	bs.join <- conn
	for {
		_, msg, err := conn.Read(bs.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				os.Exit(0)
			}
			log.Printf("read message from connection with ctx: %v", err)
			bs.leave <- conn
			return
		}
		log.Printf("read message from connection: %s", msg)
		bs.message <- msg
	}
}

func (bs *broadcastServer) ConnectionHub() {
	for {
		select {
		case joinConn := <-bs.join:
			bs.connections[joinConn] = struct{}{}
		case leaveConn := <-bs.leave:
			delete(bs.connections, leaveConn)
		case msg := <-bs.message:
			for clientConn := range bs.connections {
				log.Printf("writing message to connection: %s", msg) // TODO: identify connections for logging purposes
				err := clientConn.Write(bs.ctx, websocket.MessageText, msg)
				if err != nil {
					log.Printf("write to client connection with ctx: %v", err)
				}
			}
		}
	}
}
