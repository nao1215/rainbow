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

// S3BucketObjectsDeleterInput is the input of the DeleteBucketObjects method.
type S3BucketObjectsDeleterInput struct {
	// Bucket is the name of the bucket to delete.
	Bucket model.Bucket
	// Region is the region of the bucket that you want to delete.
	Region model.Region
	// S3ObjectSets is the list of the objects to delete.
	S3ObjectSets model.S3ObjectIdentifierSets
}

// S3BucketObjectsDeleterOutput is the output of the DeleteBucketObjects method.
type S3BucketObjectsDeleterOutput struct{}

// S3BucketObjectsDeleter is the interface that wraps the basic DeleteBucketObjects method.
type S3BucketObjectsDeleter interface {
	DeleteS3BucketObjects(ctx context.Context, input *S3BucketObjectsDeleterInput) (*S3BucketObjectsDeleterOutput, error)
}

// S3BucketObjectsListerInput is the input of the ListBucketObjects method.
type S3BucketObjectsListerInput struct {
	// Bucket is the name of the bucket to list.
	Bucket model.Bucket
}

// S3BucketObjectsListerOutput is the output of the ListBucketObjects method.
type S3BucketObjectsListerOutput struct {
	// Objects is the list of the objects.
	Objects model.S3ObjectIdentifierSets
}

// S3BucketObjectsLister is the interface that wraps the basic ListBucketObjects method.
type S3BucketObjectsLister interface {
	ListS3BucketObjects(ctx context.Context, input *S3BucketObjectsListerInput) (*S3BucketObjectsListerOutput, error)
}

// S3BucketObjectDownloaderInput is the input of the GetBucketObject method.
type S3BucketObjectDownloaderInput struct {
	// Bucket is the name of the bucket to get.
	Bucket model.Bucket
	// S3Key is the key of the object to get.
	S3Key model.S3Key
}

// S3BucketObjectDownloaderOutput is the output of the GetBucketObject method.
type S3BucketObjectDownloaderOutput struct {
	// S3Object is the object.
	S3Object *model.S3Object
}

// S3BucketObjectDownloader is the interface that wraps the basic GetBucketObject method.
type S3BucketObjectDownloader interface {
	DownloadS3BucketObject(ctx context.Context, input *S3BucketObjectDownloaderInput) (*S3BucketObjectDownloaderOutput, error)
}

// S3BucketObjectUploaderInput is the input of the PutBucketObject method.
type S3BucketObjectUploaderInput struct {
	// Bucket is the name of the bucket to put.
	Bucket model.Bucket
	// Region is the region of the bucket that you want to put.
	Region model.Region
	// S3Key is the key of the object to put.
	S3Key model.S3Key
	// S3Object is the content of the object to put.
	S3Object *model.S3Object
}

// S3BucketObjectUploaderOutput is the output of the PutBucketObject method.
type S3BucketObjectUploaderOutput struct {
	// ContentType is the content type of the object.
	ContentType string
	// ContentLength is the size of the object.
	ContentLength int64
}

// S3BucketObjectUploader is the interface that wraps the basic PutBucketObject method.
type S3BucketObjectUploader interface {
	UploadS3BucketObject(ctx context.Context, input *S3BucketObjectUploaderInput) (*S3BucketObjectUploaderOutput, error)
}
