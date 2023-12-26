package ui

import (
	"fmt"

	"github.com/muesli/termenv"
)

// General stuff for styling the view
var (
	term   = termenv.EnvColorProfile()
	subtle = makeFgStyle("241")
	red    = makeFgStyle("196")
	green  = makeFgStyle("46")
	yellow = makeFgStyle("226")
)

type (
	errMsg error
)

// makeFgStyle returns a function that will colorize the foreground of a given.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// checkbox represent [ ] and [x] items in the view.
func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

// split splits a string into multiple lines.
// Each line has a maximum length of 80 characters.
func split(s string) []string {
	var result []string
	for i := 0; i < len(s); i += 80 {
		end := i + 80
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}
