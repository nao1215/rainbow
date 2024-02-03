package s3hub

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/cmd/subcmd"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// newRmCmd return rm command.
func newRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Aliases: []string{"remove"},
		Short:   "Remove objects in S3 bucket or remove S3 bucket.",
		Example: `  [Delete a object in S3 bucket]
    s3hub rm BUCKET_NAME/S3_KEY

  [Delete all objects in S3 bucket (retain S3 bucket)]		
    s3hub rm BUCKET_NAME/*

  [Delete S3 bucket and all objects]
    s3hub rm BUCKET_NAME
     or
    s3hub rm BUCKET_NAME/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return subcmd.Run(cmd, args, &rmCmd{})
		},
	}
	cmd.Flags().StringP("profile", "p", "", "AWS profile name. if this is empty, use $AWS_PROFILE")
	// not used. however, this is common flag.
	cmd.Flags().StringP("region", "r", "", "AWS region name, default is us-east-1")
	cmd.Flags().BoolP("force", "f", false, "Force delete")
	return cmd
}

type rmCmd struct {
	// s3hub have common fields and methods for s3hub commands.
	*s3hub
	// buckets is the name of the bucket to delete.
	buckets []model.Bucket
	// force is the flag to force delete.
	force bool
}

// Parse parses command line arguments.
func (r *rmCmd) Parse(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("you must specify a bucket name")
	}

	for _, arg := range args {
		r.buckets = append(r.buckets, model.Bucket(arg))
	}
	r.s3hub = newS3hub()

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}
	r.force = force

	return r.s3hub.parse(cmd)
}

// Do executes rm command.
func (r *rmCmd) Do() error {
	if err := r.existBuckets(); err != nil {
		return err
	}
	for _, b := range r.buckets {
		bucket, key := b.Split()
		if err := r.remove(bucket, key); err != nil {
			return err
		}
	}
	return nil
}

// remove removes a bucket or a object in bucket.
func (r *rmCmd) remove(bucket model.Bucket, key model.S3Key) error {
	// delete bucket and all objects
	if key.Empty() {
		if !r.force {
			if !subcmd.Question(r.command.OutOrStdout(), fmt.Sprintf("delete %s with objects?", color.YellowString("%s", bucket))) {
				return nil
			}
		}
		if err := r.removeObjects(bucket); err != nil {
			return err
		}
		if err := r.removeBucket(bucket); err != nil {
			return err
		}
		r.printf("deleted %s\n", color.YellowString("%s", bucket))
		return nil
	}

	// delete all objects in bucket
	if key.IsAll() {
		if !r.force {
			if !subcmd.Question(r.command.OutOrStdout(), fmt.Sprintf("delete all objects in %s? (retains bucket)", color.YellowString("%s", bucket))) {
				return nil
			}
		}
		if err := r.removeObjects(bucket); err != nil {
			return err
		}
		r.printf("deleted %s with objects\n", color.YellowString("%s", bucket))
		return nil
	}

	// delete a object in bucket
	if !r.force {
		if !subcmd.Question(r.command.OutOrStdout(), fmt.Sprintf("delete %s", color.YellowString(filepath.Join(bucket.String(), key.String())))) {
			return nil
		}
	}
	if err := r.removeObject(bucket, key); err != nil {
		return err
	}
	r.printf("deleted %s\n", color.YellowString(filepath.Join(bucket.String(), key.String())))
	return nil
}

// removeObject removes a object in bucket.
func (r *rmCmd) removeObject(bucket model.Bucket, key model.S3Key) error {
	if _, err := r.S3App.S3ObjectsDeleter.DeleteS3Objects(r.ctx, &usecase.S3ObjectsDeleterInput{
		Bucket: bucket,
		S3ObjectIdentifiers: model.S3ObjectIdentifiers{
			model.S3ObjectIdentifier{
				S3Key: key,
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

// removeObjects removes all objects in bucket.
func (r *rmCmd) removeObjects(bucket model.Bucket) error {
	output, err := r.S3App.S3ObjectsLister.ListS3Objects(r.ctx, &usecase.S3ObjectsListerInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	if len(output.Objects) == 0 {
		return nil
	}

	eg, ctx := errgroup.WithContext(r.ctx)
	sem := semaphore.NewWeighted(model.MaxS3DeleteObjectsParallelsCount)
	chunks := r.divideIntoChunks(output.Objects, model.S3DeleteObjectChunksSize)

	bar := progressbar.Default(int64(output.Objects.Len()))
	for _, chunk := range chunks {
		chunk := chunk // Create a new variable to avoid concurrency issues
		// Acquire semaphore to control the number of concurrent goroutines
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		eg.Go(func() error {
			defer sem.Release(1)
			if _, err := r.S3App.S3ObjectsDeleter.DeleteS3Objects(ctx, &usecase.S3ObjectsDeleterInput{
				Bucket:              bucket,
				S3ObjectIdentifiers: chunk,
			}); err != nil {
				return err
			}
			return bar.Add(len(chunk))
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	r.printf("delete %s objects in %s\n", color.YellowString("%d", output.Objects.Len()), color.YellowString("%s", bucket))
	return nil
}

// divideIntoChunks divides a slice into chunks of the specified size.
func (r *rmCmd) divideIntoChunks(slice []model.S3ObjectIdentifier, chunkSize int) [][]model.S3ObjectIdentifier {
	var chunks [][]model.S3ObjectIdentifier

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// removeBucket removes a bucket.
// If the bucket is not empty, return error.
func (r *rmCmd) removeBucket(bucket model.Bucket) error {
	if _, err := r.S3App.S3BucketDeleter.DeleteS3Bucket(r.ctx, &usecase.S3BucketDeleterInput{
		Bucket: bucket,
	}); err != nil {
		return err
	}
	return nil
}

// existBuckets checks if the buckets exist.
// If the buckets do not exist, return error.
func (r *rmCmd) existBuckets() error {
	output, err := r.S3App.S3BucketLister.ListS3Buckets(r.ctx, &usecase.S3BucketListerInput{})
	if err != nil {
		return err
	}

	notExistBuckets := make([]string, 0, len(r.buckets))
	for _, bucket := range r.buckets {
		if output.Buckets.Contains(bucket.TrimKey()) {
			continue
		}
		b, _ := bucket.Split()
		notExistBuckets = append(notExistBuckets, b.String())
	}
	if len(notExistBuckets) == 0 {
		return nil
	}
	return fmt.Errorf("s3 bucket does not exist: %s", color.YellowString(strings.Join(notExistBuckets, ", ")))
}
