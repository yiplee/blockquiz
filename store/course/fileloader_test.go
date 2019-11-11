package course

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yiplee/blockquiz/core"
)

func TestSearchCourses(t *testing.T) {
	ctx := context.Background()

	s := &fileLoader{
		set: map[int64]bool{},
	}

	for idx := 100; idx > 0; idx-- {
		course := &core.Course{
			ID:       int64(idx),
			Language: "zh",
		}

		if idx%2 == 0 {
			course.Language = "en"
		}

		assert.Nil(t, s.insert(course))
	}

	s.sort()

	all, err := s.ListAll(ctx)
	assert.Nil(t, err)
	assert.Len(t, all, 100)
	t.Log("all", len(all))

	for idx := 1; idx <= 100; idx++ {
		course, err := s.Find(ctx, int64(idx))
		assert.Nil(t, err)
		assert.Equal(t, int64(idx), course.ID)
	}

	zh, err := s.ListLanguage(ctx, "zh")
	assert.Nil(t, err)
	t.Log("zh", len(zh))
	assert.Len(t, zh, 50)

	en, err := s.ListLanguage(ctx, "en")
	assert.Nil(t, err)
	t.Log("en", len(en))
	assert.Len(t, en, 50)
}
