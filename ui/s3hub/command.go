package s3hub

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/ui"
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
func deleteS3BucketCmd(ctx context.Context, app *di.S3App, bucket model.Bucket) tea.Cmd {
	d, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		// エラーのハンドリング
		return tea.Quit
	}
	delay := time.Millisecond * time.Duration(d.Int64())

	return tea.Tick(delay, func(t time.Time) tea.Msg {
		_, err := app.S3BucketDeleter.DeleteS3Bucket(ctx, &usecase.S3BucketDeleterInput{
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
