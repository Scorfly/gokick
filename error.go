package gokick

import "fmt"

type Error struct {
	code    int
	message string
}

func NewError(code int, message string) Error {
	return Error{code: code, message: message}
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Error() string {
	return fmt.Sprintf("Error %d: %s", e.code, e.message)
}
