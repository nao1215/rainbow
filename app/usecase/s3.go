package usecase

import "context"

// S3BucketCreatorInput is the input of the CreateBucket method.
type S3BucketCreatorInput struct {
	Bucket string
	Region string
}

// S3BucketCreatorOutput is the output of the CreateBucket method.
type S3BucketCreatorOutput struct{}

// S3BucketCreator is the interface that wraps the basic CreateBucket method.
type S3BucketCreator interface {
	CreateBucket(ctx context.Context, input *S3BucketCreatorInput) (*S3BucketCreatorOutput, error)
}
