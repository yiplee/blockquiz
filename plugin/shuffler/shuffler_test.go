package shuffler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yiplee/blockquiz/core"
)

func TestShuffleCourse(t *testing.T) {
	c := &core.Course{
		Questions: []*core.Question{
			{
				Answer: 2,
				Choices: []string{
					"1",
					"2",
					"3",
					"4",
				},
			},
		},
	}

	r := Rand()
	r.Shuffle(c, "seed", 10)

	assert.Len(t, c.Questions, 1)
	question := c.Questions[0]
	t.Log(question)
}
