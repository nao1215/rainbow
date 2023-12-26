// Package interactor contains the implementations of usecases.
package interactor

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/usecase"
)

// S3bucketCreatorSet is a provider set for S3BucketCreator.
//
//nolint:gochecknoglobals
var S3bucketCreatorSet = wire.NewSet(
	NewS3BucketCreator,
	wire.Bind(new(usecase.S3BucketCreator), new(*S3BucketCreator)),
)

var _ usecase.S3BucketCreator = (*S3BucketCreator)(nil)

// S3BucketCreator implements the S3BucketCreator interface.
type S3BucketCreator struct {
	service.S3BucketCreator
}

// NewS3BucketCreator creates a new S3BucketCreator.
func NewS3BucketCreator(c service.S3BucketCreator) *S3BucketCreator {
	return &S3BucketCreator{
		S3BucketCreator: c,
	}
}

// CreateBucket creates a new S3 bucket.
func (s *S3BucketCreator) CreateBucket(ctx context.Context, input *usecase.S3BucketCreatorInput) (*usecase.S3BucketCreatorOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	in := service.S3BucketCreatorInput{
		Bucket: input.Bucket,
		Region: input.Region,
	}
	if _, err := s.S3BucketCreator.CreateBucket(ctx, &in); err != nil {
		return nil, err
	}
	return &usecase.S3BucketCreatorOutput{}, nil
}
