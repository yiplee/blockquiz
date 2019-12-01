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
	"github.com/yiplee/blockquiz/operator/acker"
)

// ackCmd represents the ack command
var ackCmd = &cobra.Command{
	Use:   "ack",
	Short: "run blockquiz ack service",
	RunE: func(cmd *cobra.Command, args []string) error {
		db := provideDB(false)
		defer db.Close()

		commands := provideCommandStore(db)
		property := providePropertyStore(db)

		ack := acker.New(commands, property, acker.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return ack.Run(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(ackCmd)
}
