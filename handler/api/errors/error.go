package errors

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized  = New(401, "Unauthorized")
	ErrForbidden     = New(403, "Forbidden")
	ErrInternalError = New(500, "Internal Error")
)

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"msg,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func New(code int, msg string) error {
	check(code, msg)

	return &Error{
		Code:    code,
		Message: msg,
	}
}

func Unwrap(err error) (code int, msg string) {
	var Err *Error
	if errors.As(err, &Err) {
		code, msg = Err.Code, Err.Message
	} else {
		msg = err.Error()
	}

	return
}
