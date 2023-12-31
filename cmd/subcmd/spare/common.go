package spare

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/config/spare"
	"github.com/nao1215/rainbow/utils/errfmt"
	"github.com/spf13/cobra"
)

type spareOption struct {
	// command is the cobra command.
	command *cobra.Command
	// ctx is a context.Context.
	ctx context.Context
	// spare is a struct that executes the sub command.
	spare *di.SpareApp
	// config is a struct that contains the settings for the spare CLI command.
	config *spare.Config
	// debug is a flag that indicates whether to run debug mode.
	debug bool
	// configFilePath is a path of the config file.
	configFilePath string
	// awsProfile is a profile name of AWS. If this is empty, use $AWS_PROFILE.
	awsProfile model.AWSProfile
}

// newSpareOption returns a new spareOption.
func newSpareOption() *spareOption {
	return &spareOption{}
}

// Parse parses the arguments and flags.
func (s *spareOption) parseCommon(cmd *cobra.Command, _ []string) error {
	s.ctx = context.Background()
	s.command = cmd

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return errfmt.Wrap(err, "can not parse command line argument (--debug)")
	}
	s.debug = debug

	configFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return errfmt.Wrap(err, "can not parse command line argument (--file)")
	}
	if configFilePath == "" {
		configFilePath = spare.ConfigFilePath
	}
	s.configFilePath = configFilePath

	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return errfmt.Wrap(err, "can not parse command line argument (--profile)")
	}
	s.awsProfile = model.NewAWSProfile(profile)

	if err := s.readConfig(configFilePath); err != nil {
		return err
	}

	spare, err := di.NewSpareApp(s.ctx, s.awsProfile, s.config.Region)
	if err != nil {
		return err
	}
	s.spare = spare

	return nil
}

// readConfig reads .spare.yml and returns config.Config.
func (s *spareOption) readConfig(configFilePath string) error {
	file, err := os.Open(filepath.Clean(configFilePath))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	cfg := spare.NewConfig()
	if err := cfg.Read(file); err != nil {
		return err
	}
	s.config = cfg

	return nil
}

// commandName returns the s3hub command name.
func commandName() string {
	return "spare"
}
