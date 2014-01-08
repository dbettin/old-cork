package cork

import (
	"fmt"
	"runtime/debug"
)

type Error struct {
	error
	Status int
	//stacktrace
}

func NewError(err error, status int) *Error {
	fmt.Printf("stacktrace: %s", debug.Stack())
	return &Error{err, status}
}

var ErrorHandler = func(message *Message) {
	// log error
	message.ReturnStatus(message.Error.Status)
}
