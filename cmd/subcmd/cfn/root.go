// Package cfn is the root command of cfn.
package cfn

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute starts the root command of cfn command.
func Execute() error {
	if len(os.Args) == 1 {
		return interactive()
	}
	if err := newRootCmd().Execute(); err != nil {
		return err
	}
	return nil
}

// newRootCmd returns a root command for cfn.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cfn",
		Short: `cfn is a command line tool for AWS CloudFormation.`,
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.DisableFlagParsing = true

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newLsCmd())
	return cmd
}
