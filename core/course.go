package core

import (
	"context"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
)

type (
	Course struct {
		Title     string      `gorm:"size:128" json:"title,omitempty" yaml:"title"`
		Language  string      `gorm:"size:36" json:"language,omitempty" yaml:"language"`
		Summary   string      `gorm:"size:1024" json:"summary,omitempty" yaml:"summary"`
		Content   string      `gorm:"type:LONGTEXT" json:"content,omitempty" yaml:"content"`
		URL       string      `gorm:"size:256" json:"url,omitempty" yaml:"url"`
		Questions []*Question `gorm:"-" json:"questions,omitempty" yaml:"questions"`
	}

	Question struct {
		Content string   `json:"content,omitempty" yaml:"content"`
		Choices []string `json:"choices,omitempty" yaml:"choices"`
		Answer  int      `json:"answer,omitempty" yaml:"answer"` // >= 0
	}

	CourseStore interface {
		Add(ctx context.Context, course *Course) error
		ListAll(ctx context.Context) ([]*Course, error)
		ListLanguage(ctx context.Context, language string) ([]*Course, error)
		Find(ctx context.Context, title, language string) (*Course, error)
	}

	CourseShuffler interface {
		Shuffle(course *Course, userID string, questionCount int)
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

func CourseTitleByDate(at time.Time) string {
	return at.UTC().Format("2006-01-02")
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
