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
	"github.com/yiplee/blockquiz/operator/hub"
)

// blazeCmd represents the blaze command
var blazeCmd = &cobra.Command{
	Use:   "blaze",
	Short: "run blaze service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		db := provideDB(false)
		defer db.Close()
		commands := provideCommandStore(db)
		commandParser := provideParser()

		go pprof.Listen(8000)

		h := hub.New(commands, commandParser, hub.Config{
			ClientID:   cfg.Bot.ClientID,
			SessionID:  cfg.Bot.SessionID,
			SessionKey: cfg.Bot.SessionKey,
		})

		return h.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(blazeCmd)
}
