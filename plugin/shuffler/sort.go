package shuffler

import (
	"hash/fnv"
	"math/rand"
)

func Sort(seed string, n int, swap func(i, j int)) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	src := rand.NewSource(int64(h.Sum32()))
	r := rand.New(src)
	r.Shuffle(n, swap)
}
