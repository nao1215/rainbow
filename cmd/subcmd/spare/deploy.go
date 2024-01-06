package spare

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/nao1215/rainbow/config/spare"
	"github.com/nao1215/rainbow/utils/file"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// newDeployCmd return deploy sub command.
func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "deploy SPA to AWS",
		Example: "   spare deploy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &deployCmd{})
		},
	}
	cmd.Flags().BoolP("debug", "d", false, "run debug mode. you must run localstack before using this flag")
	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	cmd.Flags().StringP("file", "f", spare.ConfigFilePath, "config file path")
	return cmd
}

type deployCmd struct {
	// ctx is a context.Context.
	ctx context.Context
	// spare is a struct that executes the deploy command.
	spare *di.SpareApp
	// config is a struct that contains the settings for the spare CLI command.
	config *spare.Config
	// debug is a flag that indicates whether to run debug mode.
	debug bool
	// awsProfile is a profile name of AWS. If this is empty, use $AWS_PROFILE.
	awsProfile model.AWSProfile
}

// Parse parses the arguments and flags.
func (d *deployCmd) Parse(cmd *cobra.Command, _ []string) (err error) {
	spareOption := newSpareOption()
	if err := spareOption.parseCommon(cmd, nil); err != nil {
		return err
	}

	d.ctx = spareOption.ctx
	d.spare = spareOption.spare
	d.config = spareOption.config
	d.debug = spareOption.debug
	d.awsProfile = spareOption.awsProfile

	return nil
}

// Do deploy SPA to AWS.
func (d *deployCmd) Do() error {
	log.Info("[  MODE  ]", "debug", d.debug)
	log.Info("[ CONFIG ]", "profile", d.awsProfile)
	log.Info("[ DEPLOY ]", "target path", d.config.DeployTarget, "bucket name", d.config.S3Bucket)

	files, err := file.WalkDir(d.config.DeployTarget.String())
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(d.ctx)
	weighted := semaphore.NewWeighted(int64(runtime.NumCPU()))
	for _, file := range files {
		file := file
		eg.Go(func() error {
			if err := weighted.Acquire(ctx, 1); err != nil {
				return err
			}
			defer weighted.Release(1)

			return d.uploadFile(ctx, file)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// uploadFile uploads a file to S3.
func (d *deployCmd) uploadFile(ctx context.Context, file string) error {
	data, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return err
	}

	key := strings.Replace(file, d.config.DeployTarget.String()+string(filepath.Separator), "", 1)
	output, err := d.spare.FileUploader.UploadFile(ctx, &usecase.FileUploaderInput{
		Bucket: d.config.S3Bucket,
		Region: d.config.Region,
		// e.g. src/index.html -> index.html
		Key:  model.S3Key(key),
		Data: data, // TODO: change io.Reader?
	})
	if err != nil {
		return err
	}
	log.Info("[ DEPLOY ]", "file name", key, "mimetype", output.ContentType, "content length", output.ContentLength)
	return nil
}
