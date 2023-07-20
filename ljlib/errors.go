package ljlib

import "fmt"

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string, a ...interface{}) NotFoundError {
	return NotFoundError{
		message: fmt.Sprintf(message, a),
	}
}

func (n NotFoundError) Error() string {
	return n.message
}

func (n NotFoundError) Is(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

type IllegalArgumentError struct {
	message string
}

func NewIllegalArgumentError(message string, a ...interface{}) IllegalArgumentError {
	return IllegalArgumentError{
		message: fmt.Sprintf(message, a),
	}
}

func (i IllegalArgumentError) Error() string {
	return i.message
}

func (i IllegalArgumentError) Is(err error) bool {
	_, ok := err.(IllegalArgumentError)
	return ok
}
