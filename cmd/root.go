package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gobcast",
	Short: "A fully-fledged broadcast server over websockets written in Go",
	Long: `The broadcast server allows for multiple client to establish connection
	with a server and listen on the messages being broadcasted.

To start the websocket server use: gobcast start   [--port]
To connect a client use:           gobcast connect [--port]`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
