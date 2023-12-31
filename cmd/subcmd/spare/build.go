package spare

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/log"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/nao1215/rainbow/config/spare"
	"github.com/nao1215/spare/config"

	"github.com/spf13/cobra"
)

// newBuildCmd return build sub command.
func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Short:   "build AWS infrastructure for SPA",
		Example: "   spare build",
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &buildCmd{})
		},
	}
	cmd.Flags().BoolP("debug", "d", false, "run debug mode. you must run localstack before using this flag")
	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	cmd.Flags().StringP("file", "f", config.ConfigFilePath, "config file path")
	return cmd
}

type buildCmd struct {
	// ctx is a context.Context.
	ctx context.Context
	// spare is a struct that executes the build command.
	spare *di.SpareApp
	// config is a struct that contains the settings for the spare CLI command.
	config *spare.Config
	// configFilePath is a path of the config file.
	configFilePath string
	// debug is a flag that indicates whether to run debug mode.
	debug bool
	// awsProfile is a profile name of AWS. If this is empty, use $AWS_PROFILE.
	awsProfile model.AWSProfile
}

// Parse parses the arguments and flags.
func (b *buildCmd) Parse(cmd *cobra.Command, _ []string) (err error) {
	spareOption := newSpareOption()
	if err := spareOption.parseCommon(cmd, nil); err != nil {
		return err
	}

	b.ctx = spareOption.ctx
	b.spare = spareOption.spare
	b.config = spareOption.config
	b.configFilePath = spareOption.configFilePath
	b.debug = spareOption.debug
	b.awsProfile = spareOption.awsProfile

	return nil
}

// Do generate .spare.yml at current directory.
// If .spare.yml already exists, return error.
func (b *buildCmd) Do() error {
	log.Info(fmt.Sprintf("[VALIDATE] check %s", b.configFilePath))
	if err := b.config.Validate(b.debug); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("[VALIDATE] ok %s", b.configFilePath))

	if err := b.confirm(); err != nil {
		return err
	}

	log.Info("[ CREATE ] start building AWS infrastructure")
	log.Info("[ CREATE ] s3 bucket with public access block policy", "name", b.config.S3Bucket.String())
	if _, err := b.spare.S3BucketCreator.CreateS3Bucket(b.ctx, &usecase.S3BucketCreatorInput{
		Bucket: b.config.S3Bucket,
		Region: b.config.Region,
	}); err != nil {
		return err
	}
	if _, err := b.spare.S3BucketPublicAccessBlocker.BlockS3BucketPublicAccess(b.ctx, &usecase.S3BucketPublicAccessBlockerInput{
		Bucket: b.config.S3Bucket,
		Region: b.config.Region,
	}); err != nil {
		return err
	}
	if _, err := b.spare.S3BucketPolicySetter.SetS3BucketPolicy(b.ctx, &usecase.S3BucketPolicySetterInput{
		Bucket: b.config.S3Bucket,
		Policy: model.NewAllowCloudFrontS3BucketPolicy(b.config.S3Bucket),
	}); err != nil {
		return err
	}

	log.Info("[ CREATE ] cloudfront distribution")
	createCDNOutput, err := b.spare.CloudFrontCreator.CreateCloudFront(b.ctx, &usecase.CreateCloudFrontInput{
		Bucket: b.config.S3Bucket,
	})
	if err != nil {
		return err
	}
	log.Info("[ CREATE ] cloudfront distribution", "domain", createCDNOutput.Domain.String())

	return nil
}

// confirm shows the settings and asks if you want to build AWS infrastructure.
func (b *buildCmd) confirm() error {
	log.Info("[CONFIRM ] check the settings")
	fmt.Println("")
	fmt.Println("[debug mode]")
	fmt.Printf(" %t\n", b.debug)
	fmt.Println("[aws profile]")
	fmt.Printf(" %s\n", b.awsProfile.String())
	fmt.Printf("[%s]\n", b.configFilePath)
	fmt.Printf(" spareTemplateVersion: %s\n", b.config.SpareTemplateVersion)
	fmt.Printf(" deployTarget: %s\n", b.config.DeployTarget)
	fmt.Printf(" region: %s\n", b.config.Region)
	fmt.Printf(" customDomain: %s\n", b.config.CustomDomain)
	fmt.Printf(" s3BucketName: %s\n", b.config.S3Bucket)
	fmt.Printf(" allowOrigins: %s\n", b.config.AllowOrigins.String())
	if b.debug {
		fmt.Printf(" debugLocalstackEndpoint: %s\n", b.config.DebugLocalstackEndpoint)
	}
	fmt.Println("")

	var result bool
	if err := survey.AskOne(
		&survey.Confirm{
			Message: "want to build AWS infrastructure with the above settings?",
		},
		&result,
	); err != nil {
		return err
	}

	if !result {
		return errors.New("canceled")
	}
	return nil
}
