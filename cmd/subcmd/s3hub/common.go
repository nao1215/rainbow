package s3hub

import (
	"context"

	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/spf13/cobra"
)

// s3hub have common fields and methods for s3hub commands.
type s3hub struct {
	// S3App is the application service for S3.
	*di.S3App
	// command is the cobra command.
	command *cobra.Command
	// ctx is the context of s3hub command.
	ctx context.Context
	// profile is the AWS profile name.
	profile model.AWSProfile
	// region is the AWS region name.
	region model.Region
}

// newS3hub returns a new s3hub.
func newS3hub() *s3hub {
	return &s3hub{}
}

// parse parses command line arguments.
func (s *s3hub) parse(cmd *cobra.Command) error {
	s.command = cmd
	s.ctx = context.Background()

	p, err := cmd.Flags().GetString("profile")
	if err != nil {
		return err
	}
	s.profile = model.NewAWSProfile(p)

	r, err := cmd.Flags().GetString("region")
	if err != nil {
		return err
	}
	s.region = model.Region(r)

	cfg, err := model.NewAWSConfig(s.ctx, s.profile, s.region)
	if err != nil {
		return errfmt.Wrap(err, "can not get aws config")
	}
	if s.region == "" {
		if cfg.Config.Region == "" {
			s.region = model.RegionUSEast1
		} else {
			s.region = model.Region(cfg.Config.Region)
		}
	}

	s.S3App, err = di.NewS3App(s.ctx, s.profile, s.region)
	if err != nil {
		return errfmt.Wrap(err, "can not create s3 application service")
	}
	return nil
}

// printf prints a formatted string.
func (s *s3hub) printf(format string, a ...interface{}) {
	s.command.Printf(format, a...)
}

// commandName returns the s3hub command name.
func commandName() string {
	return "s3hub"
}
