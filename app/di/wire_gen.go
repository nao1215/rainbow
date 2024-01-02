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
	s3ObjectsLister := external.NewS3ObjectsLister(client)
	interactorS3ObjectsLister := interactor.NewS3ObjectsLister(s3ObjectsLister)
	s3ObjectsDeleter := external.NewS3ObjectsDeleter(client)
	interactorS3ObjectsDeleter := interactor.NewS3ObjectsDeleter(s3ObjectsDeleter, s3BucketLocationGetter)
	s3ObjectDownloader := external.NewS3ObjectDownloader(client)
	interactorS3ObjectDownloader := interactor.NewS3ObjectDownloader(s3ObjectDownloader)
	s3ObjectUploader := external.NewS3ObjectUploader(client)
	fileUploaderOptions := &interactor.FileUploaderOptions{
		S3ObjectUploader: s3ObjectUploader,
	}
	fileUploader := interactor.NewFileUploader(fileUploaderOptions)
	s3App := newS3App(interactorS3BucketCreator, interactorS3BucketLister, interactorS3BucketDeleter, interactorS3ObjectsLister, interactorS3ObjectsDeleter, interactorS3ObjectDownloader, fileUploader)
	return s3App, nil
}

// NewSpareApp creates a new SpareApp.
func NewSpareApp(ctx context.Context, profile model.AWSProfile, region model.Region) (*SpareApp, error) {
	awsConfig, err := model.NewAWSConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}
	client, err := external.NewCloudFrontClient(awsConfig)
	if err != nil {
		return nil, err
	}
	cloudFrontCreator := external.NewCloudFrontCreator(client)
	cloudFrontOAICreator := external.NewCloudFrontOAICreator(client)
	cloudFrontCreatorOptions := &interactor.CloudFrontCreatorOptions{
		CloudFrontCreator: cloudFrontCreator,
		OAICreator:        cloudFrontOAICreator,
	}
	interactorCloudFrontCreator := interactor.NewCloudFrontCreator(cloudFrontCreatorOptions)
	s3Client, err := external.NewS3Client(awsConfig)
	if err != nil {
		return nil, err
	}
	s3ObjectUploader := external.NewS3ObjectUploader(s3Client)
	fileUploaderOptions := &interactor.FileUploaderOptions{
		S3ObjectUploader: s3ObjectUploader,
	}
	fileUploader := interactor.NewFileUploader(fileUploaderOptions)
	s3BucketCreator := external.NewS3BucketCreator(s3Client)
	interactorS3BucketCreator := interactor.NewS3BucketCreator(s3BucketCreator)
	s3BucketPublicAccessBlocker := external.NewS3BucketPublicAccessBlocker(s3Client)
	interactorS3BucketPublicAccessBlocker := interactor.NewS3BucketPublicAccessBlocker(s3BucketPublicAccessBlocker)
	s3BucketPolicySetter := external.NewS3BucketPolicySetter(s3Client)
	interactorS3BucketPolicySetter := interactor.NewS3BucketPolicySetter(s3BucketPolicySetter)
	spareApp := newSpareApp(interactorCloudFrontCreator, fileUploader, interactorS3BucketCreator, interactorS3BucketPublicAccessBlocker, interactorS3BucketPolicySetter)
	return spareApp, nil
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
	usecase.S3ObjectsLister
	usecase.
		// S3ObjectsLister is the usecase for listing S3 bucket objects.
		S3ObjectsDeleter
	usecase.S3ObjectDownloader

	// S3ObjectsDeleter is the usecase for deleting S3 bucket objects.

	// S3ObjectUploader is the usecase for uploading a file to S3 bucket.
	usecase.FileUploader
	// FileUploader is the usecase for uploading a file.

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
	usecase.
		// CloudFrontCreator is the usecase for creating CloudFront.
		CloudFrontCreator
	usecase.FileUploader
	usecase.S3BucketCreator

	// FileUploader is the usecase for uploading a file.

	// S3BucketCreator is the usecase for creating a new S3 bucket.
	usecase.S3BucketPublicAccessBlocker
	// S3BucketPublicAccessBlocker is the usecase for blocking public access to a S3 bucket.
	usecase.S3BucketPolicySetter

	// BucketPolicySetter is the usecase for setting a bucket policy.

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
