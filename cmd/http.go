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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/handler/api"
	"golang.org/x/sync/errgroup"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "run blockquiz http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		return runServer(port)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().Int("port", 9999, "port server serve on")
}

func runServer(port int) error {
	db := provideDB()
	defer db.Close()

	tasks := provideTaskStore(db)
	commands := provideCommandStore(db)
	courses := provideCourseStore()
	handler := api.New(tasks, commands, courses)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler.Handle(),
	}

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Infof("starting api server at %s", server.Addr)

	ctx := context.Background()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(c)

	var g errgroup.Group
	g.Go(func() error {
		select {
		case <-c:
			logrus.Debug("server: shutdown")
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		return server.ListenAndServe()
	})

	return g.Wait()
}
