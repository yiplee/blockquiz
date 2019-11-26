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
	awsSession := provideAwsSession()
	pubsub := providePubSub(awsSession)

	g.Go(func() error {
		h := hub.New(commands, commandParser, pubsub, hub.Config{
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
			pubsub,
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

		return d.Run(ctx, cfg.Deliver.Capacity)
	})

	g.Go(func() error {
		m := messenger.New(messages, messenger.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return m.Run(ctx, 10*time.Millisecond)
	})

	return g.Wait()
}
