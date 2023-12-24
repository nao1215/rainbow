package s3hub

import "github.com/spf13/cobra"

// newRmCmd return rm command.
func newRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm",
		Short: "Remove contents from S3 bucket (or remove S3 bucket)",
		RunE:  rm,
	}
}

// rm is the entrypoint of rm command.
func rm(cmd *cobra.Command, _ []string) error {
	cmd.Println("rm is not implemented yet")
	return nil
}
