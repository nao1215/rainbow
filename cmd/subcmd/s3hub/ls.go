package s3hub

import (
	"errors"
	"fmt"
	"sort"

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
		Aliases: []string{"list"},
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
	// lsMode is the mode for listing.
	mode lsMode
}

// lsMode is the mode for listing.
type lsMode int

const (
	// lsModeBucket is the mode for listing buckets.
	lsModeBucket lsMode = 0
	// lsModeObject is the mode for listing objects.
	lsModeObject lsMode = 1
)

// Parse parses command line arguments.
func (l *lsCmd) Parse(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		l.bucket = model.NewBucketWithoutProtocol(args[0])
	}

	if !l.bucket.Empty() {
		l.mode = lsModeObject
		l.bucket, _ = l.bucket.Split()
	}

	l.s3hub = newS3hub()
	return l.s3hub.parse(cmd)
}

func (l *lsCmd) Do() error {
	switch l.mode {
	case lsModeBucket:
		return l.printBucket()
	case lsModeObject:
		return l.printObject()
	default:
		return errors.New("invalid mode: please report this bug: https://github.com/nao1215/rainbow")
	}
}

// printBucket prints buckets.
func (l *lsCmd) printBucket() error {
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

// printObject prints objects.
func (l *lsCmd) printObject() error {
	listBuckets, err := l.s3hub.S3BucketLister.ListS3Buckets(l.ctx, &usecase.S3BucketListerInput{})
	if err != nil {
		return err
	}
	if !listBuckets.Buckets.Contains(l.bucket) {
		return fmt.Errorf("bucket not found: %s", color.YellowString("%s", l.bucket))
	}

	listS3Objects, err := l.s3hub.S3ObjectsLister.ListS3Objects(l.ctx, &usecase.S3ObjectsListerInput{
		Bucket: l.bucket,
	})
	if err != nil {
		return err
	}

	l.printf("[S3Objects (profile=%s)]\n", l.profile.String())
	if len(listS3Objects.Objects) == 0 {
		l.printf("  No S3 Objects\n")
		return nil
	}

	sort.Sort(listS3Objects.Objects)
	for _, o := range listS3Objects.Objects {
		if o.VersionID == "" {
			l.printf("  %s/%s\n", l.bucket, o.S3Key)
			continue
		}
		l.printf("  %s/%s (version id=%s)\n", l.bucket, o.S3Key, o.VersionID)
	}
	return nil
}
