package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coder/websocket"
)

type clientConn struct {
	conn     *websocket.Conn
	nickname string
}

type broadcastServer struct {
	connections map[*clientConn]struct{}
	ctx         context.Context
	token       string
	join        chan *clientConn
	leave       chan *clientConn
	message     chan []byte
}

func NewBroadcastServer(ctx context.Context, token string) *broadcastServer {
	return &broadcastServer{
		connections: map[*clientConn]struct{}{},
		join:        make(chan *clientConn),
		leave:       make(chan *clientConn),
		message:     make(chan []byte),
		ctx:         ctx,
		token:       token,
	}
}

func (bs *broadcastServer) verifyToken(conn *websocket.Conn) (bool, error) {
	_, clientToken, err := conn.Read(bs.ctx)
	if err != nil {
		return false, fmt.Errorf("verify connection token: %w", err)
	}

	if bs.token != string(clientToken) {
		return false, fmt.Errorf("failed to verify the token %s != %s", bs.token, clientToken)
	}

	return true, nil
}

func (bs *broadcastServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientNickname := r.Header.Get("Client-Nickname")

	conn, err := websocket.Accept(w, r, nil) // don't require custom AcceptOptions
	if err != nil {
		log.Printf("upgrade HTTP conn to websocket: %v", err)
		return
	}
	defer conn.CloseNow()

	if ok, err := bs.verifyToken(conn); !ok || err != nil {
		log.Print(err)
		err = conn.Write(bs.ctx, websocket.MessageText, []byte("The connection token is invalid, please try again."))
		if err != nil {
			log.Printf("writing after failure to verify token: %v", err)
		}
		return
	} else {
		// send ack for the client to proceed
		log.Print("valid token, acknowledging to the client")
		err = conn.Write(bs.ctx, websocket.MessageText, []byte("token valid, ack"))
		if err != nil {
			log.Printf("failed the write of ack: %v", err)
			return
		}
	}

	client := &clientConn{conn: conn, nickname: clientNickname}
	bs.join <- client
	for {
		_, msg, err := conn.Read(bs.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				os.Exit(0)
			}
			log.Printf("read message from connection with ctx: %v", err)
			bs.leave <- client
			return
		}
		log.Printf("read message from connection: %s", msg)
		bs.message <- []byte(fmt.Sprintf("%s: %s", client.nickname, msg))
	}
}

func (bs *broadcastServer) ConnectionHub() {
	for {
		select {
		case joinConn := <-bs.join:
			bs.connections[joinConn] = struct{}{}
			log.Printf("new client connect -- number of connected clients: %d", len(bs.connections))
		case leaveConn := <-bs.leave:
			delete(bs.connections, leaveConn)
			log.Printf("client disconnected -- number of connected clients: %d", len(bs.connections))
		case msg := <-bs.message:
			for clientConn := range bs.connections {
				log.Printf("writing message to connection: %s", msg)
				err := clientConn.conn.Write(bs.ctx, websocket.MessageText, msg)
				if err != nil {
					log.Printf("write to client connection with ctx: %v", err)
					bs.leave <- clientConn // remove the client, if broadcasting to it errors
				}
			}
		}
	}
}
