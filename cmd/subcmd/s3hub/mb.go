package s3hub

import (
	"errors"

	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/spf13/cobra"
)

// newMbCmd return mb command. mb means make bucket.
func newMbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mb [flags] BUCKET_NAME",
		Short:   "Make S3 bucket",
		Example: "  s3hub mb -p myprofile -r us-east-1 BUCKET_NAME",
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &mbCmd{})
		},
	}
	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	cmd.Flags().StringP("region", "r", "", "AWS region name, default is us-east-1")
	return cmd
}

// mbCmd is the command for mb.
type mbCmd struct {
	// s3hub have common fields and methods for s3hub commands.
	*s3hub
	// bucket is the name of the bucket to create.
	bucket model.Bucket
}

// Parse parses command line arguments.
func (m *mbCmd) Parse(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a bucket name")
	}
	m.bucket = model.Bucket(args[0])

	m.s3hub = newS3hub()
	return m.s3hub.parse(cmd)
}

// Do executes mb command.
func (m *mbCmd) Do() error {
	_, err := m.S3BucketCreator.CreateS3Bucket(m.ctx, &usecase.S3BucketCreatorInput{
		Bucket: m.bucket,
		Region: m.s3hub.region,
	})
	if err != nil {
		return errfmt.Wrap(err, "can not create bucket")
	}

	m.printf("[Success]\n")
	m.printf("  profile: %s\n", m.profile.String())
	m.printf("  region : %s\n", m.s3hub.region)
	m.printf("  bucket : %s\n", color.YellowString("%s", m.bucket))
	return nil
}
