package s3hub

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
		Use:     "ls [flags] [BUCKET_NAME]",
		Short:   "List S3 buckets or contents of a bucket",
		Example: `  s3hub ls -p myprofile -r us-east-1`,
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
	// s3hub have common fields and methods for s3hub commands.
	*s3hub
	// bucket is the name of the bucket.
	bucket model.Bucket
}

// Parse parses command line arguments.
func (l *lsCmd) Parse(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		l.bucket = model.Bucket(args[0])
	}

	l.s3hub = newS3hub()
	return l.s3hub.parse(cmd)
}

func (l *lsCmd) Do() error {
	out, err := l.s3hub.S3BucketLister.ListS3Buckets(l.ctx, &usecase.S3BucketListerInput{})
	if err != nil {
		return err
	}

	l.printf("[Buckets (profile=%s)]\n", l.profile.String())
	if len(out.Buckets) == 0 {
		l.printf("  No Buckets\n")
		return nil
	}
	for _, b := range out.Buckets {
		l.printf("  %s (region=%s, updated_at=%s)\n",
			color.GreenString("%s", b.Bucket),
			color.YellowString("%s", b.Region),
			b.CreationDate.Format("2006-01-02 15:04:05 MST"))
	}
	return nil
}
