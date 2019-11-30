package cmd

import (
	"context"
	// _ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/operator/deliver"
	"github.com/yiplee/blockquiz/operator/hub"
	"github.com/yiplee/blockquiz/operator/messenger"
	taskcache "github.com/yiplee/blockquiz/store/task/cache"
	usercache "github.com/yiplee/blockquiz/store/user/cache"
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

	db := provideDB(false)
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

	if !cfg.Hub.Disable {
		g.Go(func() error {
			h := hub.New(commands, commandParser, hub.Config{
				ClientID:   cfg.Bot.ClientID,
				SessionID:  cfg.Bot.SessionID,
				SessionKey: cfg.Bot.SessionKey,
			})

			return h.Run(ctx)
		})
	}

	g.Go(func() error {
		d := deliver.New(
			usercache.Cache(users),
			commands,
			commandParser,
			courseShuffler,
			courses,
			wallets,
			taskcache.Cache(tasks),
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

		return m.Run(ctx)
	})

	// g.Go(func() error {
	// 	mux := http.NewServeMux()
	// 	mux.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
	// 		num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
	// 		_, _ = w.Write([]byte(num))
	// 	})
	//
	// 	return http.ListenAndServe("127.0.0.1:6061", mux)
	// })
	//
	// g.Go(func() error {
	// 	return http.ListenAndServe("127.0.0.1:6060", nil)
	// })

	return g.Wait()
}
