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
	"strconv"

	"github.com/fox-one/pkg/text/columnize"
	"github.com/spf13/cobra"
)

// loadCourseCmd represents the loadCourse command
var loadCourseCmd = &cobra.Command{
	Use:   "load-course",
	Short: "test loading courses from files",
	Run: func(cmd *cobra.Command, args []string) {
		courses := provideCourseStore()
		list, err := courses.ListAll(context.Background())
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if len(list) == 0 {
			cmd.Println("no courses loaded")
			return
		}

		cmd.Println(len(list), "courses loaded")
		var form columnize.Form
		form.Append("id", "language", "title")
		for _, course := range list {
			form.Append(strconv.FormatInt(course.ID, 10), course.Language, course.Title)
		}

		form.Fprint(cmd.OutOrStdout())
		return
	},
}

func init() {
	rootCmd.AddCommand(loadCourseCmd)
}
