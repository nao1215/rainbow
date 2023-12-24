package s3hub

import "github.com/spf13/cobra"

// newCpCmd return cp command.
func newCpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cp",
		Short: "Copy file from local(S3 bucket) to S3 bucket(local)",
		RunE:  cp,
	}
}

// cp is the entrypoint of cp command.
func cp(cmd *cobra.Command, _ []string) error {
	cmd.Println("cp is not implemented yet")
	return nil
}
