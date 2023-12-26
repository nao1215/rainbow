package s3hub

import (
	"context"
	"errors"

	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/spf13/cobra"
)

// newMbCmd return mb command. mb means make bucket.
func newMbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mb [flags] BUCKET_NAME",
		Short:   "Make S3 bucket",
		Example: "  s3hub mb -p myprofile -r us-east-1 BUCKET_NAME",
		RunE:    mb,
	}
	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	cmd.Flags().StringP("region", "r", "", "AWS region name. if this is empty, use us-east-1")
	return cmd
}

// mbCmd is the command for mb.
type mbCmd struct {
	*cobra.Command
	// ctx is the context of mb command.
	ctx context.Context
	// app is the application service for S3.
	app *di.S3App
	// profile is the AWS profile name.
	profile model.AWSProfile
	// region is the AWS region name.
	region model.Region
	// bucket is the name of the bucket to create.
	bucket model.Bucket
}

// mb is the entrypoint of mb command.
func mb(cmd *cobra.Command, args []string) error {
	mb, err := parse(cmd, args)
	if err != nil {
		return err
	}
	return mb.do()
}

// parse parses command line arguments.
func parse(cmd *cobra.Command, args []string) (*mbCmd, error) {
	if len(args) != 1 {
		return nil, errors.New("you must specify a bucket name")
	}

	ctx := context.Background()
	p, err := cmd.Flags().GetString("profile")
	if err != nil {
		return nil, errfmt.Wrap(err, "can not parse command line argument (--profile)")
	}
	profile := model.NewAWSProfile(p)

	r, err := cmd.Flags().GetString("region")
	if err != nil {
		return nil, errfmt.Wrap(err, "can not parse command line argument (--region)")
	}
	region := model.Region(r)

	cfg, err := model.NewAWSConfig(ctx, profile, region)
	if err != nil {
		return nil, errfmt.Wrap(err, "can not get aws config")
	}
	if region == "" {
		if cfg.Config.Region == "" {
			region = model.RegionUSEast1
		} else {
			region = model.Region(cfg.Config.Region)
		}
	}

	app, err := di.NewS3App(ctx, profile, region)
	if err != nil {
		return nil, errfmt.Wrap(err, "can not create s3 application service")
	}

	return &mbCmd{
		Command: cmd,
		ctx:     ctx,
		app:     app,
		profile: profile,
		region:  region,
		bucket:  model.Bucket(args[0]),
	}, nil
}

// do executes mb command.
func (mb *mbCmd) do() error {
	_, err := mb.app.S3BucketCreator.CreateBucket(mb.ctx, &usecase.S3BucketCreatorInput{
		Bucket: mb.bucket,
		Region: mb.region,
	})
	if err != nil {
		return errfmt.Wrap(err, "can not create bucket")
	}

	mb.Printf("[Success]\n")
	mb.Printf("  profile: %s\n", mb.profile.String())
	mb.Printf("  region : %s\n", mb.region)
	mb.Printf("  bucket : %s\n", color.YellowString("%s", mb.bucket))
	return nil
}
