// Package usecase has interfaces that wrap the basic business logic.
package usecase

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// S3BucketCreatorInput is the input of the CreateBucket method.
type S3BucketCreatorInput struct {
	// Bucket is the name of the bucket that you want to create.
	Bucket model.Bucket
	// Region is the region of the bucket that you want to create.
	Region model.Region
}

// S3BucketCreatorOutput is the output of the CreateBucket method.
type S3BucketCreatorOutput struct{}

// S3BucketCreator is the interface that wraps the basic CreateBucket method.
type S3BucketCreator interface {
	CreateS3Bucket(ctx context.Context, input *S3BucketCreatorInput) (*S3BucketCreatorOutput, error)
}

// S3BucketListerInput is the input of the ListBuckets method.
type S3BucketListerInput struct{}

// S3BucketListerOutput is the output of the ListBuckets method.
type S3BucketListerOutput struct {
	// Buckets is the list of the buckets.
	Buckets model.BucketSets
}

// S3BucketLister is the interface that wraps the basic ListBuckets method.
type S3BucketLister interface {
	ListS3Buckets(ctx context.Context, input *S3BucketListerInput) (*S3BucketListerOutput, error)
}

// S3ObjectsListerInput is the input of the ListObjects method.
type S3ObjectsListerInput struct {
	// Bucket is the name of the bucket that you want to list objects.
	Bucket model.Bucket
}

// S3ObjectsListerOutput is the output of the ListObjects method.
type S3ObjectsListerOutput struct {
	// Objects is the list of the objects.
	Objects model.S3ObjectIdentifierSets
}

// S3ObjectsLister is the interface that wraps the basic ListObjects method.
type S3ObjectsLister interface {
	ListS3Objects(ctx context.Context, input *S3ObjectsListerInput) (*S3ObjectsListerOutput, error)
}

// S3BucketDeleterInput is the input of the DeleteBucket method.
type S3BucketDeleterInput struct {
	// Bucket is the name of the bucket that you want to delete.
	Bucket model.Bucket
}

// S3BucketDeleterOutput is the output of the DeleteBucket method.
type S3BucketDeleterOutput struct{}

// S3BucketDeleter is the interface that wraps the basic DeleteBucket method.
type S3BucketDeleter interface {
	DeleteS3Bucket(ctx context.Context, input *S3BucketDeleterInput) (*S3BucketDeleterOutput, error)
}

// S3ObjectsDeleterInput is the input of the DeleteObjects method.
type S3ObjectsDeleterInput struct {
	// Bucket is the name of the bucket that you want to delete.
	Bucket model.Bucket
	// S3ObjectSets is the list of the objects to delete.
	S3ObjectSets model.S3ObjectIdentifierSets
}

// S3ObjectsDeleterOutput is the output of the DeleteObjects method.
type S3ObjectsDeleterOutput struct{}

// S3ObjectsDeleter is the interface that wraps the basic DeleteObjects method.
type S3ObjectsDeleter interface {
	DeleteS3Objects(ctx context.Context, input *S3ObjectsDeleterInput) (*S3ObjectsDeleterOutput, error)
}

// S3ObjectDownLoaderInput is the input of the DownloadObject method.
type S3ObjectDownloaderInput struct {
	// Bucket is the name of the bucket that you want to download.
	Bucket model.Bucket
	// Key is the S3 key.
	Key model.S3Key
}

// S3ObjectDownLoaderOutput is the output of the DownloadObject method.
type S3ObjectDownloaderOutput struct {
	// Bucket is the name of the bucket that you want to download.
	Bucket model.Bucket
	// Key is the S3 key.
	Key model.S3Key
	// ContentType is the content type of the downloaded file.
	ContentType string
	// ContentLength is the content length of the downloaded file.
	ContentLength int64
	// S3Object is the downloaded object.
	S3Object *model.S3Object
}

// S3ObjectDownloader is the interface that wraps the basic DownloadObject method.
type S3ObjectDownloader interface {
	DownloadS3Object(ctx context.Context, input *S3ObjectDownloaderInput) (*S3ObjectDownloaderOutput, error)
}

// FileUploader is an interface for uploading files to external storage.
type FileUploader interface {
	// UploadFile uploads a file from external storage.
	UploadFile(ctx context.Context, input *UploadFileInput) (*UploadFileOutput, error)
}

// UploadFileInput is an input struct for FileUploader.
type UploadFileInput struct {
	// Bucket is the name of the bucket.
	Bucket model.Bucket
	// Region is the name of the region where the bucket is located.
	Region model.Region
	// Key is the S3 key
	Key model.S3Key
	// Data is the data to upload.
	Data []byte
}

// UploadFileOutput is an output struct for FileUploader.
type UploadFileOutput struct {
	// ContentType is the content type of the uploaded file.
	ContentType string
	// ContentLength is the content length of the uploaded file.
	ContentLength int64
}
