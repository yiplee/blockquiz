package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/operator/hub"
	"golang.org/x/sync/errgroup"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run blockquiz engine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEngine(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runEngine(ctx context.Context) error {
	var g errgroup.Group

	g.Go(func() error {
		h := hub.New(nil, nil, hub.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return h.Run(ctx)
	})

	return g.Wait()
}
