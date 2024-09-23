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
			Use:   "serve",
			Short: "Run worker",
			Run: func(cmd *cobra.Command, _ []string) {
				runScheduler(ctx)
			},
			PreRun: func(cmd *cobra.Command, _ []string) {
				go func() {
					runCompleteOrderConsumer(ctx)
				}()
				go func() {
					runCreateOrderConsumer(ctx)
				}()
				go func() {
					runListTicketUpdater(ctx)
				}()
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
