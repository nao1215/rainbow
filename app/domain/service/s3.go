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
