// Package errfmt format the error message.
package errfmt

import (
	"errors"
	"fmt"
)

// Wrap return wrapping error with message.
// If e is nil, return new error with msg. If msg is empty string, return e.
// For example: Wrap(errors.New("original error"), "add message") returns "add message: original error".
func Wrap(e error, msg string) error {
	if e == nil {
		return errors.New(msg)
	}
	if msg == "" {
		return e
	}
	return fmt.Errorf("%s: %w", msg, e)
}
