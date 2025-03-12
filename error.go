package gokick

import "fmt"

type Error struct {
	code        int
	message     string
	description string
}

func NewError(code int, message string) Error {
	return Error{code: code, message: message}
}

func (e Error) WithDescription(description string) Error {
	e.description = description
	return e
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Description() string {
	return e.description
}

func (e Error) Error() string {
	if e.description == "" {
		return fmt.Sprintf("Error %d: %s", e.code, e.message)
	} else {
		return fmt.Sprintf("Error %d: %s (%s)", e.code, e.message, e.description)
	}
}
