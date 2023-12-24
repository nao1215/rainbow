package s3hub

import (
	"fmt"

	ver "github.com/nao1215/rainbow/version"
	"github.com/spf13/cobra"
)

// newVersionCmd return version command.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print %s version", commandName()),
		Run:   version,
	}
}

// version return s3hub command version.
func version(cmd *cobra.Command, _ []string) {
	cmd.Printf("%s %s (under MIT LICENSE)\n", commandName(), ver.GetVersion())
}
