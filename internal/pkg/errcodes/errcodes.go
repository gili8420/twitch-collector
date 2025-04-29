package errcodes

import (
	"strings"
)

type ErrorCode struct {
	StatusCode int
	Message    string
	Details    string
}

func New(statusCode int, message string, details ...string) *ErrorCode {
	return &ErrorCode{
		StatusCode: statusCode,
		Message:    message,
		Details:    strings.Join(details, ": "),
	}
}

func (e *ErrorCode) Error() string {
	return e.Message
}
