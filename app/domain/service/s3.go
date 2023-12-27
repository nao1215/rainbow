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
	CreateBucket(ctx context.Context, input *S3BucketCreatorInput) (*S3BucketCreatorOutput, error)
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
	ListBuckets(ctx context.Context, input *S3BucketListerInput) (*S3BucketListerOutput, error)
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
	GetBucketLocation(ctx context.Context, input *S3BucketLocationGetterInput) (*S3BucketLocationGetterOutput, error)
}
