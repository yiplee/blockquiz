/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/cmd/pprof"
	"github.com/yiplee/blockquiz/operator/deliver"
	taskcache "github.com/yiplee/blockquiz/store/task/cache"
	usercache "github.com/yiplee/blockquiz/store/user/cache"
)

// deliverCmd represents the deliver command
var deliverCmd = &cobra.Command{
	Use:   "deliver",
	Short: "run deliver service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

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

		go pprof.Listen(8001)

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
	},
}

func init() {
	rootCmd.AddCommand(deliverCmd)
}
