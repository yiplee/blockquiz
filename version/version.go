package version

import (
	"fmt"
)

var (
	Major int64 = 0
	Minor int64 = 1
	Patch int64 = 1

	Commit string
)

func String() string {
	v := fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	if Commit != "" {
		v = fmt.Sprintf("%s-%s", v, Commit)
	}

	return v
}
