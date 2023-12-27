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
	if err := input.Region.Validate(); err != nil {
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

// S3bucketListerSet is a provider set for S3BucketLister.
//
//nolint:gochecknoglobals
var S3bucketListerSet = wire.NewSet(
	NewS3BucketLister,
	wire.Bind(new(usecase.S3BucketLister), new(*S3BucketLister)),
)

var _ usecase.S3BucketLister = (*S3BucketLister)(nil)

// S3BucketLister implements the S3BucketLister interface.
type S3BucketLister struct {
	service.S3BucketLister
	service.S3BucketLocationGetter
}

// NewS3BucketLister creates a new S3BucketLister.
func NewS3BucketLister(l service.S3BucketLister, g service.S3BucketLocationGetter) *S3BucketLister {
	return &S3BucketLister{
		S3BucketLister:         l,
		S3BucketLocationGetter: g,
	}
}

// ListBuckets lists the buckets.
func (s *S3BucketLister) ListBuckets(ctx context.Context, _ *usecase.S3BucketListerInput) (*usecase.S3BucketListerOutput, error) {
	out, err := s.S3BucketLister.ListBuckets(ctx, &service.S3BucketListerInput{})
	if err != nil {
		return nil, err
	}

	for i, b := range out.Buckets {
		in := service.S3BucketLocationGetterInput{
			Bucket: b.Bucket,
		}
		o, err := s.S3BucketLocationGetter.GetBucketLocation(ctx, &in)
		if err != nil {
			return nil, err
		}
		out.Buckets[i].Region = o.Region
	}

	return &usecase.S3BucketListerOutput{
		Buckets: out.Buckets,
	}, nil
}
