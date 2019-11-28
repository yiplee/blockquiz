package property

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	now := time.Now()

	var testData = []struct {
		value interface{}
		raw   string
		date  time.Time
		num   int64
	}{
		{
			value: 123,
			raw:   "123",
			date:  time.Time{},
			num:   123,
		},
		{
			value: "haha",
			raw:   "haha",
			date:  time.Time{},
			num:   0,
		},
		{
			value: "123",
			raw:   "123",
			date:  time.Time{},
			num:   123,
		},
		{
			value: now,
			raw:   now.Format(timeLayout),
			date:  now,
			num:   0,
		},
	}

	for _, data := range testData {
		v := Parse(data.value)
		assert.Equal(t, data.raw, v.String())
		assert.Equal(t, data.date.UnixNano(), v.Time().UnixNano())
		assert.Equal(t, data.num, v.Int64())
	}
}
