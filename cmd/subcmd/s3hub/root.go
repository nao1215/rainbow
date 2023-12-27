// Package s3hub is the root command of s3hub.
package s3hub

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute starts the root command of s3hub.
func Execute() error {
	if len(os.Args) == 1 {
		return interactive()
	}
	if err := newRootCmd().Execute(); err != nil {
		return err
	}
	return nil
}

// newRootCmd returns a root command for s3hub.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "s3hub",
		Long: `s3hub is user-friendly S3 buckets management tool.
If you want to use interactive mode, run s3hub without any arguments.`,
	}
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.DisableFlagParsing = true

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newMbCmd())
	cmd.AddCommand(newLsCmd())
	cmd.AddCommand(newRmCmd())
	cmd.AddCommand(newCpCmd())
	return cmd
}
