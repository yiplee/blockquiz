package shuffler

import (
	"testing"

	"github.com/fox-one/pkg/uuid"
)

func TestSort(t *testing.T) {
	seed := uuid.New()
	t.Log(seed)
	for i := 0; i < 100; i++ {
		s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		Sort(seed, len(s), func(i, j int) {
			s[i], s[j] = s[j], s[i]
		})

		t.Log(s)
	}
}
