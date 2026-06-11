package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inflame-ue/gobcast/internal/server"
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
		portFlag := cmd.Flag("port")
		ctx, cancel := context.WithCancel(context.Background())
		wsServ := server.NewBroadcastServer(ctx)
		go wsServ.ConnectionHub()

		httpServ := http.Server{
			Addr:    ":" + portFlag.Value.String(),
			Handler: wsServ,
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigs
			log.Println("interrupt received, shutting down the server")
			cancel()

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			httpServ.Shutdown(shutdownCtx)
		}()

		log.Printf("starting server on port %s", portFlag.Value.String())
		err := httpServ.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			os.Exit(0)
		}
		log.Fatal(err)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().Uint16P("port", "p", 8080, "Specify the server listening port")
}
