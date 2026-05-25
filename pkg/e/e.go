package e

import (
	"fmt"
)

type Error struct {
	Code  int
	Msg   string
	Cause error
}

func New(code int, msg string, cause error) *Error {
	return &Error{
		Code:  code,
		Msg:   msg,
		Cause: cause,
	}
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Msg, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}
