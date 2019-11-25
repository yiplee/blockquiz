package view

import (
	"github.com/yiplee/blockquiz/core"
)

type Course struct {
	ID       int64  `json:"id,omitempty"`
	Language string `json:"language,omitempty"`
	Title    string `json:"title,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

func CourseView(course *core.Course) *Course {
	return &Course{
		Language: course.Language,
		Title:    course.Title,
		Summary:  course.Summary,
	}
}
