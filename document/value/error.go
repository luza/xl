package value

import (
	"errors"
	"fmt"
)

const (
	ErrorKindFormula = iota
	ErrorKindName
	ErrorKindRef
	ErrorKindCasting
	ErrorKindDiv0
)

type Error struct {
	error
	kind int
}

func NewError(kind int, msg string, a ...interface{}) *Error {
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a)
	}
	return &Error{
		error: errors.New(msg),
		kind:  kind,
	}
}

func (e *Error) Kind() int {
	return e.kind
}
