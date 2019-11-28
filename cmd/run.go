package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/operator/deliver"
	"github.com/yiplee/blockquiz/operator/hub"
	"github.com/yiplee/blockquiz/operator/messenger"
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
	defer db.Close()

	users := provideUserStore(db)
	commands := provideCommandStore(db)
	commandParser := provideParser()
	courseShuffler := provideShuffler()
	wallets := provideWalletStore(db)
	courses := provideCourseStore()
	tasks := provideTaskStore(db)
	localizer := provideLocalizer()
	messages := provideMessageStore(db)
	property := providePropertyStore(db)

	g.Go(func() error {
		h := hub.New(commands, commandParser, hub.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return h.Run(ctx)
	})

	g.Go(func() error {
		d := deliver.New(
			users,
			commands,
			commandParser,
			courseShuffler,
			courses,
			wallets,
			tasks,
			messages,
			property,
			localizer,
			deliver.Config{
				ClientID:      cfg.Bot.ClientID,
				SessionID:     cfg.Bot.SessionID,
				SessionKey:    cfg.Bot.SessionKey,
				ButtonColor:   cfg.Deliver.ButtonColor,
				BlockDuration: cfg.Deliver.BlockDuration,
				QuestionCount: cfg.Deliver.QuestionCount,
			},
		)

		return d.Run(ctx)
	})

	g.Go(func() error {
		m := messenger.New(messages, messenger.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return m.Run(ctx, 12*time.Millisecond)
	})

	return g.Wait()
}
