// Package xregex provides a type that represents a regular expression pattern.
package xregex

import (
	"testing"
)

func TestRegexMatchString(t *testing.T) {
	t.Parallel()
	t.Run("MatchString returns true if the string s matches the pattern", func(t *testing.T) {
		t.Parallel()

		var r Regex
		for i := 0; i < 100; i++ {
			go func() {
				r.InitOnce(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
				if err := r.MatchString("test"); err != nil {
					t.Errorf("MatchString() error = %v", err)
				}
			}()
		}
	})
}
