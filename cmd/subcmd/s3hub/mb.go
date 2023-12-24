package s3hub

import "github.com/spf13/cobra"

// newMbCmd return mb command. mb means make bucket.
func newMbCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mb",
		Short: "Make S3 bucket",
		RunE:  mb,
	}
}

// mb is the entrypoint of mb command.
func mb(cmd *cobra.Command, _ []string) error {
	cmd.Println("mb is not implemented yet")
	return nil
}
