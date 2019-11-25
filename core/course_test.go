package core

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/go-yaml/yaml"
)

func TestGenerateCourses(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	c := Course{
		Title: CourseTitleByDate(time.Now()),
	}

	operators := []string{"+", "-", "x", "%"}

	for idx := 0; idx < 100; idx++ {
		a := rand.Intn(100) + 1
		b := rand.Intn(100) + 1

		question := Question{
			Answer:  rand.Intn(len(operators)),
			Choices: make([]string, len(operators)),
		}

		question.Content = fmt.Sprintf("%d %s %d = ?", a, operators[question.Answer], b)

		for idx := range operators {
			var result int

			switch idx {
			case 0:
				result = a + b
			case 1:
				result = a - b
			case 2:
				result = a * b
			case 3:
				result = a % b
			}

			question.Choices[idx] = strconv.Itoa(result)
		}

		c.Questions = append(c.Questions, &question)
	}

	data, _ := yaml.Marshal(c)
	t.Log(string(data))
}
