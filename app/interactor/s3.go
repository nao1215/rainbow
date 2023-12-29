// Package interactor contains the implementations of usecases.
package interactor

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/usecase"
)

// S3BucketCreatorSet is a provider set for S3BucketCreator.
//
//nolint:gochecknoglobals
var S3BucketCreatorSet = wire.NewSet(
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

// CreateS3Bucket creates a new S3 bucket.
func (s *S3BucketCreator) CreateS3Bucket(ctx context.Context, input *usecase.S3BucketCreatorInput) (*usecase.S3BucketCreatorOutput, error) {
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
	if _, err := s.S3BucketCreator.CreateS3Bucket(ctx, &in); err != nil {
		return nil, err
	}
	return &usecase.S3BucketCreatorOutput{}, nil
}

// S3BucketListerSet is a provider set for S3BucketLister.
//
//nolint:gochecknoglobals
var S3BucketListerSet = wire.NewSet(
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

// ListS3Buckets lists the buckets.
func (s *S3BucketLister) ListS3Buckets(ctx context.Context, _ *usecase.S3BucketListerInput) (*usecase.S3BucketListerOutput, error) {
	out, err := s.S3BucketLister.ListS3Buckets(ctx, &service.S3BucketListerInput{})
	if err != nil {
		return nil, err
	}

	for i, b := range out.Buckets {
		in := service.S3BucketLocationGetterInput{
			Bucket: b.Bucket,
		}
		o, err := s.S3BucketLocationGetter.GetS3BucketLocation(ctx, &in)
		if err != nil {
			return nil, err
		}
		out.Buckets[i].Region = o.Region
	}

	return &usecase.S3BucketListerOutput{
		Buckets: out.Buckets,
	}, nil
}

// S3BucketObjectsLister implements the S3BucketObjectsLister interface.
type S3BucketObjectsLister struct {
	service.S3BucketObjectsLister
}

// S3BucketObjectsListerSet is a provider set for S3BucketObjectsLister.
//
//nolint:gochecknoglobals
var S3BucketObjectsListerSet = wire.NewSet(
	NewS3BucketObjectsLister,
	wire.Bind(new(usecase.S3BucketObjectsLister), new(*S3BucketObjectsLister)),
)

var _ usecase.S3BucketObjectsLister = (*S3BucketObjectsLister)(nil)

// NewS3BucketObjectsLister creates a new S3BucketObjectsLister.
func NewS3BucketObjectsLister(l service.S3BucketObjectsLister) *S3BucketObjectsLister {
	return &S3BucketObjectsLister{
		S3BucketObjectsLister: l,
	}
}

// ListS3BucketObjects lists the objects in the S3 bucket.
func (s *S3BucketObjectsLister) ListS3BucketObjects(ctx context.Context, input *usecase.S3BucketObjectsListerInput) (*usecase.S3BucketObjectsListerOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	out, err := s.S3BucketObjectsLister.ListS3BucketObjects(ctx, &service.S3BucketObjectsListerInput{
		Bucket: input.Bucket,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.S3BucketObjectsListerOutput{
		Objects: out.Objects,
	}, nil
}

// S3BucketObjectsDeleter implements the S3BucketObjectsDeleter interface.
type S3BucketObjectsDeleter struct {
	service.S3BucketObjectsDeleter
}

// S3BucketObjectsDeleterSet is a provider set for S3BucketObjectsDeleter.
//
//nolint:gochecknoglobals
var S3BucketObjectsDeleterSet = wire.NewSet(
	NewS3BucketObjectsDeleter,
	wire.Bind(new(usecase.S3BucketObjectsDeleter), new(*S3BucketObjectsDeleter)),
)

var _ usecase.S3BucketObjectsDeleter = (*S3BucketObjectsDeleter)(nil)

// NewS3BucketObjectsDeleter creates a new S3BucketObjectsDeleter.
func NewS3BucketObjectsDeleter(d service.S3BucketObjectsDeleter) *S3BucketObjectsDeleter {
	return &S3BucketObjectsDeleter{
		S3BucketObjectsDeleter: d,
	}
}

// DeleteS3BucketObjects deletes the objects in the bucket.
func (s *S3BucketObjectsDeleter) DeleteS3BucketObjects(ctx context.Context, input *usecase.S3BucketObjectsDeleterInput) (*usecase.S3BucketObjectsDeleterOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}
	_, err := s.S3BucketObjectsDeleter.DeleteS3BucketObjects(ctx, &service.S3BucketObjectsDeleterInput{
		Bucket:       input.Bucket,
		S3ObjectSets: input.S3ObjectSets,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.S3BucketObjectsDeleterOutput{}, nil
}

// S3BucketDeleterSet is a provider set for S3BucketDeleter.
//
//nolint:gochecknoglobals
var S3BucketDeleterSet = wire.NewSet(
	NewS3BucketDeleter,
	wire.Bind(new(usecase.S3BucketDeleter), new(*S3BucketDeleter)),
)

var _ usecase.S3BucketDeleter = (*S3BucketDeleter)(nil)

// S3BucketDeleter implements the S3BucketDeleter interface.
type S3BucketDeleter struct {
	service.S3BucketLocationGetter
	service.S3BucketDeleter
}

// NewS3BucketDeleter creates a new S3BucketDeleter.
func NewS3BucketDeleter(
	s3BucketDeleter service.S3BucketDeleter,
	s3BucketLocationGetter service.S3BucketLocationGetter,
) *S3BucketDeleter {
	return &S3BucketDeleter{
		S3BucketDeleter:        s3BucketDeleter,
		S3BucketLocationGetter: s3BucketLocationGetter,
	}
}

// DeleteS3Bucket deletes the bucket.
func (s *S3BucketDeleter) DeleteS3Bucket(ctx context.Context, input *usecase.S3BucketDeleterInput) (*usecase.S3BucketDeleterOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	location, err := s.S3BucketLocationGetter.GetS3BucketLocation(ctx, &service.S3BucketLocationGetterInput{
		Bucket: input.Bucket,
	})
	if err != nil {
		return nil, err
	}

	if _, err := s.S3BucketDeleter.DeleteS3Bucket(ctx, &service.S3BucketDeleterInput{
		Bucket: input.Bucket,
		Region: location.Region,
	}); err != nil {
		return nil, err
	}
	return &usecase.S3BucketDeleterOutput{}, nil
}
