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
	ctx  context.Context
	conn *websocket.Conn
}

func NewBroadcastClient(ctx context.Context, conn *websocket.Conn) *broadcastClient {
	return &broadcastClient{
		ctx:  ctx,
		conn: conn,
	}
}

func readStdin() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Message to Broadcast: ")
	msg, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read from stdin: %w", err)
	}
	return msg, nil
}

func (bc *broadcastClient) Broadcast() {
	for {
		msg, err := readStdin()
		if err != nil {
			log.Print(err)
		}

		err = bc.conn.Write(bc.ctx, websocket.MessageText, msg)
		if err != nil {
			log.Fatalf("writing to websocket connection: %v", err)
		}
	}
}

func (bc *broadcastClient) Receive() {
	for {
		_, msg, err := bc.conn.Read(bc.ctx)
		if err != nil {
			log.Fatalf("reading from websocket connection: %v", err)
		}
		fmt.Print(msg)
	}
}
