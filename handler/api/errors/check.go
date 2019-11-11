package errors

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var DisableDuplicateErrorCheck bool

func init() {
	DisableDuplicateErrorCheck, _ = strconv.ParseBool(
		os.Getenv("DISABLE_DUPLICATE_ERROR_CHECK"),
	)
}

type duplicatedCodeCheck struct {
	codes map[int]string
	mux   sync.Mutex
}

var globalCheck = &duplicatedCodeCheck{
	codes: make(map[int]string),
}

func check(code int, msg string) {
	if !DisableDuplicateErrorCheck {
		globalCheck.check(code, msg)
	}
}

func (c *duplicatedCodeCheck) check(code int, msg string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if m, ok := c.codes[code]; ok {
		if m != msg {
			panic(fmt.Errorf("code %d has bind with %s", code, m))
		}
	} else {
		c.codes[code] = msg
	}
}
