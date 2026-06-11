package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coder/websocket"
)

type broadcastClient struct {
	ctx     context.Context
	conn    *websocket.Conn
	message chan []byte
	errors  chan error
}

func NewBroadcastClient(ctx context.Context, conn *websocket.Conn) *broadcastClient {
	return &broadcastClient{
		ctx:     ctx,
		conn:    conn,
		message: make(chan []byte),
		errors:  make(chan error),
	}
}

func (bc *broadcastClient) ReadStdin() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to Broadcast: ")
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("read from stdin: %v", err)
		}
		bc.message <- msg
	}
}

func (bc *broadcastClient) Broadcast() {
	for {
		select {
		case msg := <-bc.message:
			err := bc.conn.Write(bc.ctx, websocket.MessageText, msg)
			if err != nil {
				log.Fatalf("writing to websocket connection: %v", err)
			}
		case err := <-bc.errors:
			log.Fatalf("broadcast interrupted: %v", err)
		}

	}
}

func (bc *broadcastClient) Receive() {
	for {
		_, msg, err := bc.conn.Read(bc.ctx)
		if err != nil {
			bc.errors <- err
			log.Fatalf("reading from websocket connection: %v", err)
		}
		fmt.Printf("%s", msg)
	}
}
