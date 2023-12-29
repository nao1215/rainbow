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

// S3BucketObjectsListerInput is the input of the ListObjects method.
type S3BucketObjectsListerInput struct {
	// Bucket is the name of the bucket that you want to list objects.
	Bucket model.Bucket
}

// S3BucketObjectsListerOutput is the output of the ListObjects method.
type S3BucketObjectsListerOutput struct {
	// Objects is the list of the objects.
	Objects model.S3ObjectSets
}

// S3BucketObjectsLister is the interface that wraps the basic ListObjects method.
type S3BucketObjectsLister interface {
	ListS3BucketObjects(ctx context.Context, input *S3BucketObjectsListerInput) (*S3BucketObjectsListerOutput, error)
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

// S3BucketObjectsDeleterInput is the input of the DeleteObjects method.
type S3BucketObjectsDeleterInput struct {
	// Bucket is the name of the bucket that you want to delete.
	Bucket model.Bucket
	// S3ObjectSets is the list of the objects to delete.
	S3ObjectSets model.S3ObjectSets
}

// S3BucketObjectsDeleterOutput is the output of the DeleteObjects method.
type S3BucketObjectsDeleterOutput struct{}

// S3BucketObjectsDeleter is the interface that wraps the basic DeleteObjects method.
type S3BucketObjectsDeleter interface {
	DeleteS3BucketObjects(ctx context.Context, input *S3BucketObjectsDeleterInput) (*S3BucketObjectsDeleterOutput, error)
}
