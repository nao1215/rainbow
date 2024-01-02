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
	usecase.S3BucketCreator
	// S3BucketLister is the usecase for listing S3 buckets.
	usecase.S3BucketLister
	// S3BucketDeleter is the usecase for deleting a S3 bucket.
	usecase.S3BucketDeleter
	// S3ObjectsLister is the usecase for listing S3 bucket objects.
	usecase.S3ObjectsLister
	// S3ObjectsDeleter is the usecase for deleting S3 bucket objects.
	usecase.S3ObjectsDeleter
	// S3ObjectUploader is the usecase for uploading a file to S3 bucket.
	usecase.S3ObjectDownloader
	// FileUploader is the usecase for uploading a file.
	usecase.FileUploader
}

// NewS3App creates a new S3App.
func NewS3App(ctx context.Context, profile model.AWSProfile, region model.Region) (*S3App, error) {
	wire.Build(
		model.NewAWSConfig,
		external.NewS3Client,
		external.S3BucketCreatorSet,
		external.S3BucketListerSet,
		external.S3BucketLocationGetterSet,
		external.S3BucketDeleterSet,
		external.S3ObjectsListerSet,
		external.S3ObjectsDeleterSet,
		external.S3ObjectDownloaderSet,
		external.S3ObjectUploaderSet,
		interactor.S3BucketCreatorSet,
		interactor.S3BucketListerSet,
		interactor.S3BucketDeleterSet,
		interactor.S3ObjectsListerSet,
		interactor.S3ObjectsDeleterSet,
		interactor.S3ObjectDownloaderSet,
		interactor.FileUploaderSet,
		newS3App,
	)
	return nil, nil
}

func newS3App(
	s3BucketCreator usecase.S3BucketCreator,
	s3BucketLister usecase.S3BucketLister,
	s3BucketDeleter usecase.S3BucketDeleter,
	S3ObjectsLister usecase.S3ObjectsLister,
	S3ObjectsDeleter usecase.S3ObjectsDeleter,
	s3ObjectDownloader usecase.S3ObjectDownloader,
	fileUploader usecase.FileUploader,
) *S3App {
	return &S3App{
		S3BucketCreator:    s3BucketCreator,
		S3BucketLister:     s3BucketLister,
		S3BucketDeleter:    s3BucketDeleter,
		S3ObjectsLister:    S3ObjectsLister,
		S3ObjectsDeleter:   S3ObjectsDeleter,
		S3ObjectDownloader: s3ObjectDownloader,
		FileUploader:       fileUploader,
	}
}

// SpareApp is the application service for spare command.
type SpareApp struct {
	// CloudFrontCreator is the usecase for creating CloudFront.
	usecase.CloudFrontCreator
	// FileUploader is the usecase for uploading a file.
	usecase.FileUploader
	// S3BucketCreator is the usecase for creating a new S3 bucket.
	usecase.S3BucketCreator
	// S3BucketPublicAccessBlocker is the usecase for blocking public access to a S3 bucket.
	usecase.S3BucketPublicAccessBlocker
	// BucketPolicySetter is the usecase for setting a bucket policy.
	usecase.S3BucketPolicySetter
}

// NewSpareApp creates a new SpareApp.
func NewSpareApp(ctx context.Context, profile model.AWSProfile, region model.Region) (*SpareApp, error) {
	wire.Build(
		model.NewAWSConfig,
		external.NewCloudFrontClient,
		external.CloudFrontCreatorSet,
		external.OAICreatorSet,
		external.NewS3Client,
		external.S3BucketCreatorSet,
		external.S3ObjectUploaderSet,
		external.S3BucketPublicAccessBlockerSet,
		external.S3BucketPolicySetterSet,
		interactor.CloudFrontCreatorSet,
		interactor.FileUploaderSet,
		interactor.S3BucketCreatorSet,
		interactor.S3BucketPublicAccessBlockerSet,
		interactor.S3BucketPolicySetterSet,
		newSpareApp,
	)
	return nil, nil
}

// newSpareApp creates a new SpareApp.
func newSpareApp(
	cloudFrontCreator usecase.CloudFrontCreator,
	fileUploader usecase.FileUploader,
	s3BucketCreator usecase.S3BucketCreator,
	s3BucketPublicAccessBlocker usecase.S3BucketPublicAccessBlocker,
	s3BucketPolicySetter usecase.S3BucketPolicySetter,
) *SpareApp {
	return &SpareApp{
		CloudFrontCreator:           cloudFrontCreator,
		FileUploader:                fileUploader,
		S3BucketCreator:             s3BucketCreator,
		S3BucketPublicAccessBlocker: s3BucketPublicAccessBlocker,
		S3BucketPolicySetter:        s3BucketPolicySetter,
	}
}
