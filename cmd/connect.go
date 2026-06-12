package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/coder/websocket"
	"github.com/inflame-ue/gobcast/internal/client"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect a client to the websocket server over a port",
	Long: `Connect a client to the websocket server over a port

gobcast connect [--port] allows you to connect a client to the websocket server
that will listen to the broadcasted messages over the specified port.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		portFlag := cmd.Flag("port")
		hostFlag := cmd.Flag("host")
		tokenFlag := cmd.Flag("token")
		nicknameFlag := cmd.Flag("nickname")

		connString := fmt.Sprintf("ws://%s:%s", hostFlag.Value.String(), portFlag.Value.String())
		ctx, cancel := context.WithCancel(context.Background())
		options := &websocket.DialOptions{
			HTTPHeader: http.Header{"Client-Nickname": []string{nicknameFlag.Value.String()},},
		}

		conn, _, err := websocket.Dial(ctx, connString, options)
		if err != nil {
			log.Fatalf("dialing websocket at port %s: %v", portFlag.Value.String(), err)
		}

		// connection token handshake
		err = conn.Write(ctx, websocket.MessageText, []byte(tokenFlag.Value.String()))
		if err != nil {
			log.Fatalf("failed to write token to the server: %v", err)
		}
		_, msg, err := conn.Read(ctx)
		if err != nil || strings.Contains(string(msg), "invalid") {
			fmt.Printf("token %s is invalid, please try again\n", tokenFlag.Value.String())
			cancel()
			return
		}

		wsClient := client.NewBroadcastClient(ctx, conn)
		go wsClient.PrintStdin()
		go wsClient.ReadStdin()
		go wsClient.Broadcast()
		go wsClient.Receive()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Println("\ninterrupt received, shutting down client")
			cancel()
			conn.CloseNow()
			os.Exit(0)
		}()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().Uint16P("port", "p", 8080, "Specify the client connection port")
	connectCmd.Flags().String("host", "localhost", "Specify the websocket host")
	connectCmd.Flags().StringP("token", "t", "", "Specify the connection token the server expects(must be passed in!)")
	connectCmd.Flags().String("nickname", "", "Specify the nickname identificator, which will be pased to all client's on a message")
}
