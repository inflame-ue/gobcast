package cmd

import (
	"context"
	"log"

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
		ctx := context.Background()

		conn, _, err := websocket.Dial(ctx, "ws://localhost:"+portFlag.Value.String(), nil)
		if err != nil {
			log.Fatalf("dialing websocket at port %s: %v", portFlag.Value.String(), err)
		}
		defer conn.CloseNow()

		wsClient := client.NewBroadcastClient(ctx, conn)
		go wsClient.ReadStdin()
		go wsClient.Broadcast()
		go wsClient.Receive()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().Uint16P("port", "p", 8080, "Specify the client connection port")
}
