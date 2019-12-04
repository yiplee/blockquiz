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
	"github.com/yiplee/blockquiz/operator/messenger"
)

// messengerCmd represents the messenger command
var messengerCmd = &cobra.Command{
	Use:   "messenger",
	Short: "run message service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		db := provideDB(false)
		defer db.Close()

		messages := provideMessageStore(db)
		m := messenger.New(messages, messenger.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return m.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(messengerCmd)
}