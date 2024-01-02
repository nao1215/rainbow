// Package service
package service

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// S3BucketCreatorInput is the input of the CreateBucket method.
type S3BucketCreatorInput struct {
	// Bucket is the name of the bucket to create.
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

// S3BucketLocationGetterInput is the input of the GetBucketLocation method.
type S3BucketLocationGetterInput struct {
	Bucket model.Bucket
}

// S3BucketLocationGetterOutput is the output of the GetBucketLocation method.
type S3BucketLocationGetterOutput struct {
	// Region is the region of the bucket.
	Region model.Region
}

// S3BucketLocationGetter is the interface that wraps the basic GetBucketLocation method.
type S3BucketLocationGetter interface {
	GetS3BucketLocation(ctx context.Context, input *S3BucketLocationGetterInput) (*S3BucketLocationGetterOutput, error)
}

// S3BucketDeleterInput is the input of the DeleteBucket method.
type S3BucketDeleterInput struct {
	// Bucket is the name of the bucket to delete.
	Bucket model.Bucket
	// Region is the region of the bucket that you want to delete.
	Region model.Region
}

// S3BucketDeleterOutput is the output of the DeleteBucket method.
type S3BucketDeleterOutput struct{}

// S3BucketDeleter is the interface that wraps the basic DeleteBucket method.
type S3BucketDeleter interface {
	DeleteS3Bucket(ctx context.Context, input *S3BucketDeleterInput) (*S3BucketDeleterOutput, error)
}

// S3ObjectsDeleterInput is the input of the DeleteBucketObjects method.
type S3ObjectsDeleterInput struct {
	// Bucket is the name of the bucket to delete.
	Bucket model.Bucket
	// Region is the region of the bucket that you want to delete.
	Region model.Region
	// S3ObjectSets is the list of the objects to delete.
	S3ObjectSets model.S3ObjectIdentifierSets
}

// S3ObjectsDeleterOutput is the output of the DeleteBucketObjects method.
type S3ObjectsDeleterOutput struct{}

// S3ObjectsDeleter is the interface that wraps the basic DeleteBucketObjects method.
type S3ObjectsDeleter interface {
	DeleteS3Objects(ctx context.Context, input *S3ObjectsDeleterInput) (*S3ObjectsDeleterOutput, error)
}

// S3ObjectsListerInput is the input of the ListBucketObjects method.
type S3ObjectsListerInput struct {
	// Bucket is the name of the bucket to list.
	Bucket model.Bucket
}

// S3ObjectsListerOutput is the output of the ListBucketObjects method.
type S3ObjectsListerOutput struct {
	// Objects is the list of the objects.
	Objects model.S3ObjectIdentifierSets
}

// S3ObjectsLister is the interface that wraps the basic ListBucketObjects method.
type S3ObjectsLister interface {
	ListS3Objects(ctx context.Context, input *S3ObjectsListerInput) (*S3ObjectsListerOutput, error)
}

// S3ObjectDownloaderInput is the input of the GetBucketObject method.
type S3ObjectDownloaderInput struct {
	// Bucket is the name of the bucket to get.
	Bucket model.Bucket
	// Key is the key of the object to get.
	Key model.S3Key
}

// S3ObjectDownloaderOutput is the output of the GetBucketObject method.
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

// S3ObjectDownloader is the interface that wraps the basic GetBucketObject method.
type S3ObjectDownloader interface {
	DownloadS3Object(ctx context.Context, input *S3ObjectDownloaderInput) (*S3ObjectDownloaderOutput, error)
}

// S3ObjectUploaderInput is the input of the PutBucketObject method.
type S3ObjectUploaderInput struct {
	// Bucket is the name of the bucket to put.
	Bucket model.Bucket
	// Region is the region of the bucket that you want to put.
	Region model.Region
	// S3Key is the key of the object to put.
	S3Key model.S3Key
	// S3Object is the content of the object to put.
	S3Object *model.S3Object
}

// S3ObjectUploaderOutput is the output of the PutBucketObject method.
type S3ObjectUploaderOutput struct {
	// ContentType is the content type of the object.
	ContentType string
	// ContentLength is the size of the object.
	ContentLength int64
}

// S3ObjectUploader is the interface that wraps the basic PutBucketObject method.
type S3ObjectUploader interface {
	UploadS3Object(ctx context.Context, input *S3ObjectUploaderInput) (*S3ObjectUploaderOutput, error)
}

// S3ObjectCopierInput is the input of the CopyBucketObject method.
type S3ObjectCopierInput struct {
	// SourceBucket is the name of the source bucket.
	SourceBucket model.Bucket
	// SourceKey is the key of the source object.
	SourceKey model.S3Key
	// DestinationBucket is the name of the destination bucket.
	DestinationBucket model.Bucket
	// DestinationKey is the key of the destination object.
	DestinationKey model.S3Key
}

// S3ObjectCopierOutput is the output of the CopyBucketObject method.
type S3ObjectCopierOutput struct{}

// S3ObjectCopier is the interface that wraps the basic CopyBucketObject method.
type S3ObjectCopier interface {
	CopyS3Object(ctx context.Context, input *S3ObjectCopierInput) (*S3ObjectCopierOutput, error)
}
