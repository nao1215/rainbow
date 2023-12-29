package subcmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Doer is an interface that represents the behavior of a command.
type Doer interface {
	Do() error
}

// SubCommand is an interface that represents the behavior of a command.
type SubCommand interface {
	Parse(cmd *cobra.Command, args []string) error
	Doer
}

// Run runs the subcommand.
func Run(cmd *cobra.Command, args []string, subCmd SubCommand) error {
	if err := subCmd.Parse(cmd, args); err != nil {
		return err
	}
	return subCmd.Do()
}

// FmtScanln is wrapper for fmt.Scanln(). It's for unit test.
var FmtScanln = fmt.Scanln

// Question displays the question in the terminal and receives an answer from the user.
func Question(w io.Writer, ask string) bool {
	var response string

	fmt.Fprintf(w, "%s: %s", color.GreenString("CHECK"), ask+" [Y/n] ")
	_, err := FmtScanln(&response)
	if err != nil {
		// If user input only enter.
		if strings.Contains(err.Error(), "expected newline") {
			return Question(w, ask)
		}
		fmt.Fprint(os.Stderr, err.Error())
		return false
	}
	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return Question(w, ask)
	}
}
