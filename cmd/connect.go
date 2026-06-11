package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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
		ctx, cancel := context.WithCancel(context.Background())

		conn, _, err := websocket.Dial(ctx, "ws://localhost:"+portFlag.Value.String(), nil)
		if err != nil {
			log.Fatalf("dialing websocket at port %s: %v", portFlag.Value.String(), err)
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
}
