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
		l.printf("  %s (status=%s, updated_at=%s)\n",
			color.GreenString(*stack.StackName),
			stackStatusString(stack.StackStatus),
			stack.LastUpdatedTime.Format("2006-01-02 15:04:05"))
	}
	return nil
}

// stackStatusString returns a string representation of the stack status.
func stackStatusString(status model.StackStatus) string {
	switch status {
	case model.StackStatusCreateComplete:
		return color.GreenString("CREATE_COMPLETE")
	case model.StackStatusCreateFailed:
		return color.RedString("CREATE_FAILED")
	case model.StackStatusCreateInProgress:
		return color.YellowString("CREATE_IN_PROGRESS")
	case model.StackStatusDeleteComplete:
		return color.GreenString("DELETE_COMPLETE")
	case model.StackStatusDeleteFailed:
		return color.RedString("DELETE_FAILED")
	case model.StackStatusDeleteInProgress:
		return color.YellowString("DELETE_IN_PROGRESS")
	case model.StackStatusRollbackComplete:
		return color.GreenString("ROLLBACK_COMPLETE")
	case model.StackStatusRollbackFailed:
		return color.RedString("ROLLBACK_FAILED")
	case model.StackStatusRollbackInProgress:
		return color.YellowString("ROLLBACK_IN_PROGRESS")
	case model.StackStatusUpdateComplete:
		return color.GreenString("UPDATE_COMPLETE")
	case model.StackStatusUpdateCompleteCleanupInProgress:
		return color.YellowString("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS")
	case model.StackStatusUpdateInProgress:
		return color.YellowString("UPDATE_IN_PROGRESS")
	case model.StackStatusUpdateRollbackComplete:
		return color.GreenString("UPDATE_ROLLBACK_COMPLETE")
	case model.StackStatusUpdateRollbackCompleteCleanupInProgress:
		return color.YellowString("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS")
	case model.StackStatusUpdateRollbackFailed:
		return color.RedString("UPDATE_ROLLBACK_FAILED")
	case model.StackStatusUpdateFailed:
		return color.RedString("UPDATE_FAILED")
	case model.StackStatusUpdateRollbackInProgress:
		return color.YellowString("UPDATE_ROLLBACK_IN_PROGRESS")
	case model.StackStatusReviewInProgress:
		return color.YellowString("REVIEW_IN_PROGRESS")
	case model.StackStatusImportInProgress:
		return color.YellowString("IMPORT_IN_PROGRESS")
	case model.StackStatusImportComplete:
		return color.GreenString("IMPORT_COMPLETE")
	case model.StackStatusImportRollbackInProgress:
		return color.YellowString("IMPORT_ROLLBACK_IN_PROGRESS")
	case model.StackStatusImportRollbackFailed:
		return color.RedString("IMPORT_ROLLBACK_FAILED")
	case model.StackStatusImportRollbackComplete:
		return color.GreenString("IMPORT_ROLLBACK_COMPLETE")
	default:
		return color.RedString("UNKNOWN")
	}
}
