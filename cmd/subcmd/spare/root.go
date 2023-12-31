// Package spare is a package that contains subcommands for the spare CLI command.
package spare

import (
	"github.com/spf13/cobra"
)

// Execute starts the root command of s3hub.
func Execute() error {
	if err := newRootCmd().Execute(); err != nil {
		return err
	}
	return nil
}

// newRootCmd creates a new root command. This command is the entry point of the CLI.
// It is responsible for parsing the command line arguments and flags, and then
// executing the appropriate subcommand. It also sets up logging and error handling.
// The root command does not have any functionality of its own.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spare",
		Short: "spare release single page application and aws infrastructure",
	}
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newBuildCmd())
	cmd.AddCommand(newDeployCmd())
	return cmd
}
