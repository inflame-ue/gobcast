package cmd

import (
	"fmt"

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
		fmt.Println("connect called")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().Int64P("port", "p", 8080, "Specify the client connection port")
}
