package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func Start() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rootCmd := &cobra.Command{}
	cmd := []*cobra.Command{
		{
			Use:   "serve-http",
			Short: "Run HTTP server",
			Run: func(cmd *cobra.Command, _ []string) {
				runHTTPServer(ctx)
			},
			PreRun: func(cmd *cobra.Command, _ []string) {
				go func() {
					runCreateOrderConsumer(ctx)
				}()
				go func() {
					runCancelOrderConsumer(ctx)
				}()
				go func() {
					runNotifyUpdateListTicketConsumer(ctx)
				}()
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
