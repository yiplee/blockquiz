package shuffler

import (
	"github.com/yiplee/blockquiz/core"
)

type rand struct{}

func Rand() core.CourseShuffler {
	return &rand{}
}

func (r *rand) Shuffle(course *core.Course, userID string, questionCount int) {
	seed := userID + course.Title

	questions := course.Questions
	Sort(seed, len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	if len(questions) > questionCount {
		questions = questions[:questionCount]
	}

	for _, question := range questions {
		Sort(seed, len(question.Choices), func(i, j int) {
			switch question.Answer {
			case i:
				question.Answer = j
			case j:
				question.Answer = i
			}

			question.Choices[i], question.Choices[j] = question.Choices[j], question.Choices[i]
		})
	}

	course.Questions = questions
}
