package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a websocket server over the specified port",
	Long: `Start a websocket server over the specified port.

gobcast start [--port] allows you to start a server that will listen on --port
and accept client connections into a broadcast pool, which will be notified in full,
if any single client sends a message.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().Uint16P("port", "p", 8080, "Specify the server listening port")
}
