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
	"strconv"

	"github.com/fox-one/pkg/text/columnize"
	"github.com/spf13/cobra"
	"github.com/yiplee/blockquiz/plugin/shuffler"
)

// loadCourseCmd represents the loadCourse command
var loadCourseCmd = &cobra.Command{
	Use:   "load-course",
	Short: "test loading courses from files",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		courses := provideCourseStore()
		list, err := courses.ListAll(ctx)
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
		form.Append("num", "language", "title", "questions")
		for idx, course := range list {
			form.Append(strconv.Itoa(idx+1), course.Language, course.Title, strconv.Itoa(len(course.Questions)))
		}

		form.Fprint(cmd.OutOrStdout())

		userID := "8017d200-7870-4b82-b53f-74bae1d2dad7"
		title := list[0].Title
		language := list[0].Language

		for idx := 0; idx < 2; idx++ {
			course, _ := courses.Find(ctx, title, language)
			fmt.Fprintln(cmd.OutOrStdout(), course.Questions[0].Content)
			shuffler.Rand().Shuffle(course, userID, 5)
			for idx, question := range course.Questions {
				fmt.Fprintf(cmd.OutOrStdout(), "%d %v\n", idx, question)
			}
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(loadCourseCmd)
}
