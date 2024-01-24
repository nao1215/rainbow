package cfn

import (
	"context"

	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/spf13/cobra"
)

// cfn have common fields and methods for cfn commands.
type cfn struct {
	// CFnApp is the application service for CloudFormation.
	*di.CFnApp
	// command is the cobra command.
	command *cobra.Command
	// ctx is the context of s3hub command.
	ctx context.Context
	// profile is the AWS profile name.
	profile model.AWSProfile
	// region is the AWS region name.
	region model.Region
}

// newCFn returns a new cfn.
func newCFn() *cfn {
	return &cfn{}
}

// parse parses command line arguments.
func (c *cfn) parse(cmd *cobra.Command) error {
	c.command = cmd
	c.ctx = context.Background()

	p, err := cmd.Flags().GetString("profile")
	if err != nil {
		return err
	}
	c.profile = model.NewAWSProfile(p)

	r, err := cmd.Flags().GetString("region")
	if err != nil {
		return err
	}
	c.region = model.Region(r)

	cfg, err := model.NewAWSConfig(c.ctx, c.profile, c.region)
	if err != nil {
		return errfmt.Wrap(err, "can not get aws config")
	}
	if c.region == "" {
		if cfg.Config.Region == "" {
			c.region = model.RegionUSEast1
		} else {
			c.region = model.Region(cfg.Config.Region)
		}
	}

	c.CFnApp, err = di.NewCFnApp(c.ctx, c.profile, c.region)
	if err != nil {
		return errfmt.Wrap(err, "can not create cloudformation application service")
	}
	return nil
}

// printf prints a formatted string.
func (c *cfn) printf(format string, a ...interface{}) {
	c.command.Printf(format, a...)
}

// commandName returns the cfn command name.
func commandName() string {
	return "cfn"
}
