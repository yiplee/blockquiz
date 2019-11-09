package core

import (
	"context"
)

type (
	Lesson struct {
		Chapter   int64       `json:"chapter,omitempty"`
		Language  string      `gorm:"size:36" json:"language,omitempty"`
		Title     string      `gorm:"size:128" json:"title,omitempty"`
		Summary   string      `gorm:"size:1024" json:"summary,omitempty"`
		Content   string      `gorm:"type:LONGTEXT" json:"content,omitempty"`
		URL       string      `gorm:"size:256" json:"url,omitempty"`
		Questions []*Question `gorm:"-" json:"question,omitempty"`
	}

	Question struct {
		Content string   `json:"content,omitempty"`
		Choices []string `json:"choice,omitempty"`
		Answer  Answer   `json:"answer,omitempty"`
	}
)

type Answer int

func (a Answer) String() string {
	return string([]byte{'A' + byte(a)})
}

func (lesson *Lesson) Question(idx int) (*Question, bool) {
	if questions := lesson.Questions[:]; idx >= 0 && idx < len(questions) {
		return questions[idx], true
	}

	return nil, false
}

type LessonStore interface {
	Add(ctx context.Context, lesson *Lesson) error
	ListAll(ctx context.Context, language string) ([]*Lesson, error)
	Find(ctx context.Context, language string, chapter int64) (*Lesson, error)
	FindNext(ctx context.Context, language string, chapter int64) (*Lesson, error)
}
