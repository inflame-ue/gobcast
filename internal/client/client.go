package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/coder/websocket"
)

type broadcastClient struct {
	ctx     context.Context
	conn    *websocket.Conn
	message chan []byte
	errors  chan error
	print   chan string
}

func NewBroadcastClient(ctx context.Context, conn *websocket.Conn) *broadcastClient {
	return &broadcastClient{
		ctx:     ctx,
		conn:    conn,
		message: make(chan []byte),
		errors:  make(chan error),
		print:   make(chan string),
	}
}

func (bc *broadcastClient) PrintStdin() {
	for msg := range bc.print {
		fmt.Print(msg)
	}
}

func (bc *broadcastClient) ReadStdin() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Message to Broadcast: ")
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
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
		case _ = <-bc.ctx.Done():
			os.Exit(0)
		}

	}
}

func (bc *broadcastClient) Receive() {
	for {
		_, msg, err := bc.conn.Read(bc.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				os.Exit(0)
			}
			bc.errors <- err
			log.Printf("reading from websocket connection: %v", err)
		}
		bc.print <- string(msg)
		bc.print <- "Message to Broadcast: \n"
	}
}
