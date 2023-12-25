// Package errfmt format the error message.
package errfmt

import (
	"errors"
	"fmt"
)

// Wrap return wrapping error with message.
// If e is nil, return new error with msg. If msg is empty string, return e.
// For example: Wrap(errors.New("original error"), "add message") returns "original error: add message".
func Wrap(e error, msg string) error {
	if e == nil {
		return errors.New(msg)
	}
	if msg == "" {
		return e
	}
	return fmt.Errorf("%w: %s", e, msg)
}
