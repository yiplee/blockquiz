package parser

import (
	"strconv"
	"strings"
)

type Args []string

func newArgs(s string) Args {
	return Args(strings.Fields(s))
}

func (a Args) Get(idx int) (string, bool) {
	if idx >= 0 && idx < len(a) {
		return a[idx], true
	}

	return "", false
}

func (a Args) First() string {
	arg, _ := a.Get(0)
	return arg
}

func (a Args) GetInt64(idx int) (int64, bool) {
	arg, ok := a.Get(idx)
	if !ok {
		return 0, false
	}

	v, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}

func (a Args) GetInt(idx int) (int, bool) {
	arg, ok := a.Get(idx)
	if !ok {
		return 0, false
	}

	v, err := strconv.Atoi(arg)
	if err != nil {
		return 0, false
	}

	return v, true
}

func (a Args) Encode() string {
	return strings.Join(a, " ")
}
