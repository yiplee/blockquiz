package store

import (
	"errors"

	"github.com/fox-one/pkg/store/db"
)

var ErrNotFound = errors.New("store: not found")

func IsErrNotFound(err error) bool {
	for err != nil {
		if err == ErrNotFound || db.IsErrorNotFound(err) {
			return true
		}

		err = errors.Unwrap(err)
	}

	return false
}

// ErrOptimisticLock is returned by if the struct being
// modified has a Version field and the value is not equal
// to the current value in the database
var ErrOptimisticLock = db.ErrOptimisticLock

func IsErrOptimisticLock(err error) bool {
	return errors.Is(err, ErrOptimisticLock)
}
