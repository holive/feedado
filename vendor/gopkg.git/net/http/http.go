package http

import (
	"fmt"
	"net/http"
)

// Runner is a interface that executes HTTP requests.
type Runner interface {
	Do(*http.Request) (*http.Response, error)
}

// Error encapsulates error details as returned from http Runner.
type Error struct {
	Status int
}

func (e *Error) Error() string {
	return fmt.Sprintf("status code %d", e.Status)
}

// NewError returns a  initialized error.
func NewError(status int) *Error {
	return &Error{status}
}
