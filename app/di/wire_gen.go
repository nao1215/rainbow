// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"context"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/external"
	"github.com/nao1215/rainbow/app/interactor"
	"github.com/nao1215/rainbow/app/usecase"
)

// Injectors from wire.go:

// NewS3App creates a new S3App.
func NewS3App(ctx context.Context, profile model.AWSProfile, region model.Region) (*S3App, error) {
	awsConfig, err := model.NewAWSConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}
	client, err := external.NewS3Client(awsConfig)
	if err != nil {
		return nil, err
	}
	s3BucketCreator := external.NewS3BucketCreator(client)
	interactorS3BucketCreator := interactor.NewS3BucketCreator(s3BucketCreator)
	s3BucketLister := external.NewS3BucketLister(client)
	s3BucketLocationGetter := external.NewS3BucketLocationGetter(client)
	interactorS3BucketLister := interactor.NewS3BucketLister(s3BucketLister, s3BucketLocationGetter)
	s3BucketDeleter := external.NewS3BucketDeleter(client)
	interactorS3BucketDeleter := interactor.NewS3BucketDeleter(s3BucketDeleter, s3BucketLocationGetter)
	s3BucketObjectsLister := external.NewS3BucketObjectsLister(client)
	interactorS3BucketObjectsLister := interactor.NewS3BucketObjectsLister(s3BucketObjectsLister)
	s3BucketObjectsDeleter := external.NewS3BucketObjectsDeleter(client)
	interactorS3BucketObjectsDeleter := interactor.NewS3BucketObjectsDeleter(s3BucketObjectsDeleter)
	s3App := newS3App(interactorS3BucketCreator, interactorS3BucketLister, interactorS3BucketDeleter, interactorS3BucketObjectsLister, interactorS3BucketObjectsDeleter)
	return s3App, nil
}

// wire.go:

// S3App is the application service for S3.
type S3App struct {
	usecase.
		// S3BucketCreator is the usecase for creating a new S3 bucket.
		S3BucketCreator
	usecase.S3BucketLister
	usecase.S3BucketDeleter

	// S3BucketLister is the usecase for listing S3 buckets.

	// S3BucketDeleter is the usecase for deleting a S3 bucket.
	usecase.S3BucketObjectsLister
	// S3BucketObjectsLister is the usecase for listing S3 bucket objects.
	usecase.S3BucketObjectsDeleter

	// S3BucketObjectsDeleter is the usecase for deleting S3 bucket objects.

}

func newS3App(
	s3BucketCreator usecase.S3BucketCreator,
	s3BucketLister usecase.S3BucketLister,
	s3BucketDeleter usecase.S3BucketDeleter,
	s3BucketObjectsLister usecase.S3BucketObjectsLister,
	s3BucketObjectsDeleter usecase.S3BucketObjectsDeleter,
) *S3App {
	return &S3App{
		S3BucketCreator:        s3BucketCreator,
		S3BucketLister:         s3BucketLister,
		S3BucketDeleter:        s3BucketDeleter,
		S3BucketObjectsLister:  s3BucketObjectsLister,
		S3BucketObjectsDeleter: s3BucketObjectsDeleter,
	}
}
