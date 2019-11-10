package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/operator/cashier"
	"github.com/yiplee/blockquiz/operator/deliver"
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

	db := provideDB()
	users := provideUserStore(db)
	commands := provideCommandStore(db)
	commandParser := provideParser()
	wallets := provideWalletStore(db)
	walletz := provideWalletService()
	courses := provideCourseStore(db)
	localizer := provideLocalizer()

	g.Go(func() error {
		h := hub.New(commands, commandParser, hub.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return h.Run(ctx)
	})

	g.Go(func() error {
		c := cashier.New(wallets, walletz)
		return c.Work(ctx, 300*time.Millisecond)
	})

	g.Go(func() error {
		d := deliver.New(
			users,
			commands,
			commandParser,
			courses,
			wallets,
			localizer,
			deliver.Config{
				ClientID:    cfg.Bot.ClientID,
				SessionID:   cfg.Bot.SessionID,
				SessionKey:  cfg.Bot.SessionKey,
				CoinAsset:   cfg.Course.CoinAsset,
				CoinAmount:  cfg.Course.CoinAmount,
				ButtonColor: cfg.Deliver.ButtonColor,
			},
		)

		return d.Run(ctx, 200*time.Millisecond)
	})

	return g.Wait()
}
