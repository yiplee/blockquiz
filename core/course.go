package core

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type (
	Course struct {
		ID        int64       `gorm:"PRIMARY_KEY" json:"id,omitempty"`
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
		Answer  int      `json:"answer,omitempty"` // >= 0
	}

	CourseStore interface {
		Add(ctx context.Context, course *Course) error
		ListAll(ctx context.Context) ([]*Course, error)
		ListLanguage(ctx context.Context, language string) ([]*Course, error)
		Find(ctx context.Context, id int64) (*Course, error)
		FindNext(ctx context.Context, course *Course) (*Course, error)
	}
)

func AnswerToString(answer int) string {
	return string([]byte{'A' + byte(answer)})
}

func (lesson *Course) Question(idx int) (*Question, bool) {
	if questions := lesson.Questions[:]; idx >= 0 && idx < len(questions) {
		return questions[idx], true
	}

	return nil, false
}

func ValidateCourse(course *Course) error {
	if !govalidator.IsIn(course.Language, ActionSwitchEnglish, ActionSwitchChinese) {
		return fmt.Errorf("language must be zh or en")
	}

	if len(course.Questions) == 0 {
		return fmt.Errorf("questions is empty")
	}

	for idx, question := range course.Questions {
		if len(question.Choices) < 4 {
			return fmt.Errorf("question %d lack of choice", idx)
		}

		if question.Answer >= len(question.Choices) {
			return fmt.Errorf("question %d answer out of range", idx)
		}
	}

	return nil
}
