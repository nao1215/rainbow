package s3hub

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/ui"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// fetchS3BucketMsg is the message that is sent when the user wants to fetch the list of the S3 buckets.
type fetchS3BucketMsg struct {
	buckets model.BucketSets
}

// fetchS3BucketListCmd fetches the list of the S3 buckets.
func fetchS3BucketListCmd(ctx context.Context, app *di.S3App) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		output, err := app.S3BucketLister.ListS3Buckets(ctx, &usecase.S3BucketListerInput{})
		if err != nil {
			return ui.ErrMsg(err)
		}
		return fetchS3BucketMsg{
			buckets: output.Buckets,
		}
	})
}

type deleteS3BucketMsg struct {
	deletedBucket model.Bucket
}

// deleteS3BucketCmd deletes the S3 bucket.
// TODO: refactor
func deleteS3BucketCmd(ctx context.Context, app *di.S3App, bucket model.Bucket) tea.Cmd {
	d, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		return func() tea.Msg {
			return ui.ErrMsg(fmt.Errorf("failed to start deleting s3 bucket: %w", err))
		}
	}
	delay := time.Millisecond * time.Duration(d.Int64())

	return tea.Tick(delay, func(t time.Time) tea.Msg {
		output, err := app.S3ObjectsLister.ListS3Objects(ctx, &usecase.S3ObjectsListerInput{
			Bucket: bucket,
		})
		if err != nil {
			return err
		}

		if len(output.Objects) != 0 {
			eg, ctx := errgroup.WithContext(ctx)
			sem := semaphore.NewWeighted(model.MaxS3DeleteObjectsParallelsCount)
			chunks := divideIntoChunks(output.Objects, model.S3DeleteObjectChunksSize)

			for _, chunk := range chunks {
				chunk := chunk // Create a new variable to avoid concurrency issues
				// Acquire semaphore to control the number of concurrent goroutines
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}

				eg.Go(func() error {
					defer sem.Release(1)
					if _, err := app.S3ObjectsDeleter.DeleteS3Objects(ctx, &usecase.S3ObjectsDeleterInput{
						Bucket:       bucket,
						S3ObjectSets: chunk,
					}); err != nil {
						return err
					}
					return nil
				})
			}

			if err := eg.Wait(); err != nil {
				return err
			}
		}

		_, err = app.S3BucketDeleter.DeleteS3Bucket(ctx, &usecase.S3BucketDeleterInput{
			Bucket: bucket,
		})
		if err != nil {
			return ui.ErrMsg(err)
		}
		return deleteS3BucketMsg{
			deletedBucket: bucket,
		}
	})
}

// divideIntoChunks divides a slice into chunks of the specified size.
func divideIntoChunks(slice []model.S3ObjectIdentifier, chunkSize int) [][]model.S3ObjectIdentifier {
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
