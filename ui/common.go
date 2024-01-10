package ui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/muesli/termenv"
)

// General stuff for styling the view
var (
	Term   = termenv.EnvColorProfile()
	Subtle = MakeFgStyle("241")
	Red    = MakeFgStyle("196")
	Green  = MakeFgStyle("46")
	Yellow = MakeFgStyle("226")
)

type (
	// ErrMsg is an error message.
	ErrMsg error
)

// MakeFgStyle returns a function that will colorize the foreground of a given.
func MakeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(Term.Color(color)).Styled
}

// ColorFg a string's foreground with the given value.
func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(Term.Color(color)).String()
}

// Checkbox represent [ ] and [x] items in the view.
func Checkbox(label string, checked bool) string {
	if checked {
		return ColorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Split splits a string into multiple lines.
// Each line has a maximum length of 80 characters.
func Split(s string) []string {
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

// GoodByeMessage returns a goodbye message.
func GoodByeMessage() string {
	s := fmt.Sprintf("\n  See you later 🌈\n  %s\n  %s\n\n",
		"Following URL for bug reports and encouragement (e.g. GitHub Star ⭐️ )",
		color.GreenString("https://github.com/nao1215/rainbow"))
	return s
}

// ErrorMessage returns an error message.
func ErrorMessage(err error) string {
	message := fmt.Sprintf("%s\n", Red("[Error]"))
	for _, line := range Split(err.Error()) {
		message += fmt.Sprintf("  %s\n", Red(line))
	}
	return message
}

// Choice represents a choice.
type Choice struct {
	Choice int
	Max    int
	Min    int
}

// NewChoice returns a new choice.
func NewChoice(min, max int) *Choice {
	return &Choice{
		Choice: min,
		Max:    max,
		Min:    min,
	}
}

// Increment increments the choice.
// If the choice is greater than the maximum, the choice is set to the minimum.
func (c *Choice) Increment() {
	c.Choice++
	if c.Choice > c.Max {
		c.Choice = c.Min
	}
}

// Decrement decrements the choice.
// If the choice is less than the minimum, the choice is set to the maximum.
func (c *Choice) Decrement() {
	c.Choice--
	if c.Choice < c.Min {
		c.Choice = c.Max
	}
}
