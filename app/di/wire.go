//go:build wireinject
// +build wireinject

// Package di Inject dependence by wire command.
package di

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/external"
	"github.com/nao1215/rainbow/app/interactor"
	"github.com/nao1215/rainbow/app/usecase"
)

// S3App is the application service for S3.
type S3App struct {
	// S3BucketCreator is the usecase for creating a new S3 bucket.
	S3BucketCreator usecase.S3BucketCreator
	// S3BucketLister is the usecase for listing S3 buckets.
	S3BucketLister usecase.S3BucketLister
}

// NewS3App creates a new S3App.
func NewS3App(ctx context.Context, profile model.AWSProfile, region model.Region) (*S3App, error) {
	wire.Build(
		model.NewAWSConfig,
		external.NewS3Client,
		external.S3BucketCreatorSet,
		external.S3BucketListerSet,
		external.S3BucketLocationGetterSet,
		interactor.S3bucketCreatorSet,
		interactor.S3bucketListerSet,
		newS3App,
	)
	return nil, nil
}

func newS3App(s3bucketCreator usecase.S3BucketCreator, s3bucketLister usecase.S3BucketLister) *S3App {
	return &S3App{
		S3BucketCreator: s3bucketCreator,
		S3BucketLister:  s3bucketLister,
	}
}
