package ui

import (
	"fmt"

	"github.com/muesli/termenv"
)

// General stuff for styling the view
var (
	term   = termenv.EnvColorProfile()
	subtle = makeFgStyle("241")
	dot    = colorFg(" • ", "236")
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
