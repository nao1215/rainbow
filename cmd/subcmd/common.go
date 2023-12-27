package subcmd

import "github.com/spf13/cobra"

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
