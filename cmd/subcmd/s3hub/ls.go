package s3hub

import "github.com/spf13/cobra"

// newLsCmd return ls command.
func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List S3 buckets",
		RunE:  ls,
	}
}

// ls is the entrypoint of ls command.
func ls(cmd *cobra.Command, _ []string) error {
	cmd.Println("ls is not implemented yet")
	return nil
}
