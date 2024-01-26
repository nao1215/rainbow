package s3hub

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/gogf/gf/os/gfile"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/nao1215/rainbow/utils/file"
	"github.com/spf13/cobra"
)

// newCpCmd return cp command.
func newCpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cp [flags] SOURCE_PATH DESTINATION_PATH",
		Aliases: []string{"copy"},
		Short:   "Copy file from local(S3 bucket) to S3 bucket(local)",
		Example: `  [S3 bucket to local]
    s3hub cp -p myprofile -r us-east-1 s3://mybucket/path/to/file.txt /path/to/file.txt

  [local to S3 bucket]
    s3hub cp -p myprofile -r us-east-1 /path/to/file.txt s3://mybucket/path/to/file.txt

  [S3 bucket to S3 bucket]
    s3hub cp -p myprofile -r us-east-1 s3://mybucket1/path/to/file.txt s3://mybucket2/path/to/file.txt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &cpCmd{})
		},
	}

	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	cmd.Flags().StringP("region", "r", "", "AWS region name, default is us-east-1")
	return cmd
}

type cpCmd struct {
	// s3hub have common fields and methods for s3hub commands.
	*s3hub
	// pair is a slice of CopyPathPair.
	pair *copyPathPair
}

// copyType is a type of copy.
type copyType int

const (
	// copyTypeUnknown is a type of copy that is unknown.
	copyTypeUnknown copyType = -1
	// copyTypeLocalToS3 is a type of copy from local to S3.
	copyTypeLocalToS3 copyType = 0
	// copyTypeS3ToLocal is a type of copy from S3 to local.
	copyTypeS3ToLocal copyType = 1
	// copyTypeS3ToS3 is a type of copy from S3 to S3.
	copyTypeS3ToS3 copyType = 2
)

// copyPathPair is a pair of paths.
type copyPathPair struct {
	// From is a path of source.
	From string
	// To is a path of destination.
	To string
	// Type is represents a type of copy.
	// Type: "from local to S3", "from S3 to local", "from S3 to S3".
	Type copyType
}

// newCopyPathPair returns a new copyPathPair.
func newCopyPathPair(from, to string) *copyPathPair {
	pair := &copyPathPair{
		From: from,
		To:   to,
	}
	pair.Type = pair.copyType()
	return pair
}

// copyType returns a type of copy.
func (c *copyPathPair) copyType() copyType {
	if c.From == "" {
		return copyTypeUnknown
	}
	if c.To == "" {
		return copyTypeUnknown
	}
	if strings.HasPrefix(c.From, model.S3Protocol) && !strings.HasPrefix(c.To, model.S3Protocol) {
		return copyTypeS3ToLocal
	}
	if !strings.HasPrefix(c.From, model.S3Protocol) && strings.HasPrefix(c.To, model.S3Protocol) {
		return copyTypeLocalToS3
	}
	if strings.HasPrefix(c.From, model.S3Protocol) && strings.HasPrefix(c.To, model.S3Protocol) {
		return copyTypeS3ToS3
	}
	return copyTypeUnknown
}

// Parse parses command line arguments.
func (c *cpCmd) Parse(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("you must specify copy %s and %s",
			color.YellowString("source path(arg1)"), color.YellowString("destination path(arg2)"))
	}

	c.pair = newCopyPathPair(args[0], args[1])
	c.s3hub = newS3hub()
	return c.s3hub.parse(cmd)
}

// Do executes cp command.
func (c *cpCmd) Do() error {
	switch c.pair.Type {
	case copyTypeLocalToS3:
		return c.localToS3()
	case copyTypeS3ToLocal:
		return c.s3ToLocal()
	case copyTypeS3ToS3:
		return c.s3ToS3()
	case copyTypeUnknown:
		fallthrough
	default:
		return fmt.Errorf("unsupported copy type. from=%s, to=%s",
			color.YellowString(c.pair.From), color.YellowString(c.pair.To))
	}
}

// copyTargetsInLocal returns a slice of target files in local.
func (c *cpCmd) copyTargetsInLocal() ([]string, error) {
	if gfile.IsFile(c.pair.From) {
		return []string{c.pair.From}, nil
	}
	targets, err := file.WalkDir(c.pair.From)
	if err != nil {
		return nil, err
	}
	return targets, nil
}

// localToS3 copies from local to S3.
func (c *cpCmd) localToS3() error {
	targets, err := c.copyTargetsInLocal()
	if err != nil {
		return err
	}

	toBucket, toKey := model.NewBucketWithoutProtocol(c.pair.To).Split()
	fileNum := len(targets)

	for i, v := range targets {
		data, err := os.ReadFile(filepath.Clean(v))
		if err != nil {
			return err
		}
		if _, err := c.s3hub.FileUploader.UploadFile(c.ctx, &usecase.FileUploaderInput{
			Bucket: toBucket,
			Region: c.s3hub.region,
			Key:    model.S3Key(filepath.Join(toKey.String(), filepath.Base(v))),
			Data:   data,
		}); err != nil {
			return err
		}
		c.printf("[%d/%d] copy %s to %s\n",
			i+1,
			fileNum,
			color.YellowString(v),
			color.YellowString(toBucket.Join(toKey).WithProtocol().String()),
		)
	}
	return nil
}

// s3ToLocal copies from S3 to local.
func (c *cpCmd) s3ToLocal() error {
	fromBucket, fromKey := model.NewBucketWithoutProtocol(c.pair.From).Split()
	_, toKey := model.NewBucketWithoutProtocol(c.pair.To).Split()

	listOutput, err := c.s3hub.ListS3Objects(c.ctx, &usecase.S3ObjectsListerInput{
		Bucket: fromBucket,
	})
	if err != nil {
		return err
	}

	targets := make([]model.S3Key, 0, len(listOutput.Objects))
	for _, v := range listOutput.Objects {
		if strings.Contains(v.S3Key.String(), fromKey.String()) {
			targets = append(targets, v.S3Key)
		}
	}

	if len(targets) == 0 {
		return fmt.Errorf("no objects found. bucket=%s, key=%s",
			color.YellowString(fromBucket.String()), color.YellowString(fromKey.String()))
	}

	fileNum := len(targets)
	for i, v := range targets {
		downloadOutput, err := c.s3hub.S3ObjectDownloader.DownloadS3Object(c.ctx, &usecase.S3ObjectDownloaderInput{
			Bucket: fromBucket,
			Key:    v,
		})
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(fromKey.String(), v.String())
		if err != nil {
			return err
		}
		destinationPath := filepath.Join(toKey.String(), relativePath)
		if err := downloadOutput.S3Object.ToFile(destinationPath, 0644); err != nil {
			return err
		}

		c.printf("[%d/%d] copy %s to %s\n",
			i+1,
			fileNum,
			color.YellowString(fromBucket.Join(v).WithProtocol().String()),
			color.YellowString(destinationPath),
		)
	}
	return nil
}

// s3ToS3 copies from S3 to S3.
func (c *cpCmd) s3ToS3() error {
	fromBucket, fromKey := model.NewBucketWithoutProtocol(c.pair.From).Split()
	toBucket, toKey := model.NewBucketWithoutProtocol(c.pair.To).Split()

	listOutput, err := c.s3hub.ListS3Objects(c.ctx, &usecase.S3ObjectsListerInput{
		Bucket: fromBucket,
	})
	if err != nil {
		return err
	}

	targets := make([]model.S3Key, 0, len(listOutput.Objects))
	for _, v := range listOutput.Objects {
		if strings.Contains(v.S3Key.String(), fromKey.String()) {
			targets = append(targets, v.S3Key)
		}
	}

	if len(targets) == 0 {
		return fmt.Errorf("no objects found. bucket=%s, key=%s", color.YellowString(fromBucket.String()), color.YellowString(fromKey.String()))
	}

	fileNum := len(targets)
	for i, v := range targets {
		relativePath, err := filepath.Rel(fromKey.String(), v.String())
		if err != nil {
			return err
		}
		destinationKey := model.S3Key(filepath.Join(toKey.String(), relativePath))

		if _, err := c.s3hub.S3ObjectCopier.CopyS3Object(c.ctx, &usecase.S3ObjectCopierInput{
			SourceBucket:      fromBucket,
			SourceKey:         v, // from key
			DestinationBucket: toBucket,
			DestinationKey:    destinationKey,
		}); err != nil {
			return err
		}
		c.printf("[%d/%d] copy %s to %s\n",
			i+1,
			fileNum,
			color.YellowString(fromBucket.Join(v).WithProtocol().String()),
			color.YellowString(toBucket.Join(destinationKey).WithProtocol().String()),
		)
	}
	return nil
}
