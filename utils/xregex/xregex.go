// Package xregex provides a type that represents a regular expression pattern.
package xregex

import (
	"fmt"
	"regexp"
	"sync"
)

// Regex is a type that represents a regular expression pattern.
type Regex struct {
	// Pattern is the regular expression pattern.
	Pattern *regexp.Regexp
	// Mutex is a mutual exclusion lock.
	Mutex sync.Mutex
	// Once is an object that will perform exactly one action.
	Once sync.Once
}

// InitOnce initializes the Regex object.
func (r *Regex) InitOnce(pattern string) {
	r.Once.Do(func() {
		r.Pattern = regexp.MustCompile(pattern)
	})
}

// MatchString returns true if the string s matches the pattern.
func (r *Regex) MatchString(s string) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	if !r.Pattern.MatchString(s) {
		return fmt.Errorf("does not match the regular expression pattern: %s", r.Pattern.String())
	}
	return nil
}
