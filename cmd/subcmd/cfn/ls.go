package cfn

import (
	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/spf13/cobra"
)

// newLsCmd return ls command.
func newLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [flags]",
		Aliases: []string{"list"},
		Short:   "List CloudFormation stacks",
		Example: `  cfn ls -p myprofile -r us-east-1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &lsCmd{})
		},
	}
	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	// not used. however, this is common flag.
	cmd.Flags().StringP("region", "r", "", "AWS region name, default is us-east-1")
	return cmd
}

// lsCmd is the command for ls.
type lsCmd struct {
	// cfn have common fields and methods for cfn commands.
	*cfn
}

// Parse parses command line arguments.
func (l *lsCmd) Parse(cmd *cobra.Command, _ []string) error {
	l.cfn = newCFn()
	return l.cfn.parse(cmd)
}

func (l *lsCmd) Do() error {
	out, err := l.CFnStackLister.ListCFnStack(l.ctx, &usecase.CFnStackListerInput{
		Region: l.cfn.region,
	})
	if err != nil {
		return err
	}

	l.printf("[CloudFormation Stack (profile=%s, region=%s)]\n", l.profile.String(), l.region)
	if len(out.Stacks) == 0 {
		l.printf("  No stacks\n")
		return nil
	}

	for _, stack := range out.Stacks {
		if stack.StackName == nil || stack.StackStatus == model.StackStatusDeleteComplete {
			continue
		}

		l.printf("  %s (status=%s, updated_at=%s)\n",
			color.GreenString(*stack.StackName),
			stack.StackStatus.StringWithColor(),
			stack.LastUpdatedTime.Format("2006-01-02 15:04:05"))
	}
	return nil
}
