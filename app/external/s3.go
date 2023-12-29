// Package external implements the external service.
package external

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

// NewS3Client creates a new S3 service client.
// If profile is empty, the default profile is used.
func NewS3Client(cfg *model.AWSConfig) (*s3.Client, error) {
	return s3.NewFromConfig(*cfg.Config), nil
}

// S3BucketCreator implements the S3BucketCreator interface.
type S3BucketCreator struct {
	client *s3.Client
}

// S3BucketCreatorSet is a provider set for S3BucketCreator.
//
//nolint:gochecknoglobals
var S3BucketCreatorSet = wire.NewSet(
	NewS3BucketCreator,
	wire.Bind(new(service.S3BucketCreator), new(*S3BucketCreator)),
)

var _ service.S3BucketCreator = (*S3BucketCreator)(nil)

// NewS3BucketCreator creates a new S3BucketCreator.
func NewS3BucketCreator(client *s3.Client) *S3BucketCreator {
	return &S3BucketCreator{client: client}
}

// CreateS3Bucket creates a new S3 bucket.
func (c *S3BucketCreator) CreateS3Bucket(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
	// If region is us-east-1, you must not specify the location constraint.
	// If you specify the location constraint in this case, the following error will occur.
	// [api error InvalidLocationConstraint: The specified location-constraint is not valid]
	locationContstraint := &types.CreateBucketConfiguration{
		LocationConstraint: types.BucketLocationConstraint(input.Region.String()),
	}
	if input.Region == model.RegionUSEast1 {
		locationContstraint = nil
	}

	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket:                    aws.String(input.Bucket.String()),
		CreateBucketConfiguration: locationContstraint,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: region=%s, bucket name=%s", err, input.Region.String(), input.Bucket.String())
	}
	return &service.S3BucketCreatorOutput{}, nil
}

// S3BucketLister implements the S3BucketLister interface.
type S3BucketLister struct {
	client *s3.Client
}

// S3BucketListerSet is a provider set for S3BucketLister.
//
//nolint:gochecknoglobals
var S3BucketListerSet = wire.NewSet(
	NewS3BucketLister,
	wire.Bind(new(service.S3BucketLister), new(*S3BucketLister)),
)

var _ service.S3BucketLister = (*S3BucketLister)(nil)

// NewS3BucketLister creates a new S3BucketLister.
func NewS3BucketLister(client *s3.Client) *S3BucketLister {
	return &S3BucketLister{client: client}
}

// ListS3Buckets lists the buckets.
func (c *S3BucketLister) ListS3Buckets(ctx context.Context, _ *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
	out, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var buckets model.BucketSets
	for _, b := range out.Buckets {
		buckets = append(buckets, model.BucketSet{
			Bucket:       model.Bucket(*b.Name),
			CreationDate: *b.CreationDate,
		})
	}
	return &service.S3BucketListerOutput{Buckets: buckets}, nil
}

// S3BucketLocationGetter implements the S3BucketLocationGetter interface.
type S3BucketLocationGetter struct {
	client *s3.Client
}

// S3BucketLocationGetterSet is a provider set for S3BucketLocationGetter.
//
//nolint:gochecknoglobals
var S3BucketLocationGetterSet = wire.NewSet(
	NewS3BucketLocationGetter,
	wire.Bind(new(service.S3BucketLocationGetter), new(*S3BucketLocationGetter)),
)

var _ service.S3BucketLocationGetter = (*S3BucketLocationGetter)(nil)

// NewS3BucketLocationGetter creates a new S3BucketLocationGetter.
func NewS3BucketLocationGetter(client *s3.Client) *S3BucketLocationGetter {
	return &S3BucketLocationGetter{client: client}
}

// GetS3BucketLocation gets the location of the bucket.
func (c *S3BucketLocationGetter) GetS3BucketLocation(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
	out, err := c.client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(input.Bucket.String()),
	})
	if err != nil {
		return nil, err
	}

	region := model.Region(out.LocationConstraint)
	if region == "" {
		region = model.RegionUSEast1
	}

	return &service.S3BucketLocationGetterOutput{
		Region: region,
	}, nil
}

// S3BucketDeleter implements the S3BucketDeleter interface.
type S3BucketDeleter struct {
	client *s3.Client
}

// S3BucketDeleterSet is a provider set for S3BucketDeleter.
//
//nolint:gochecknoglobals
var S3BucketDeleterSet = wire.NewSet(
	NewS3BucketDeleter,
	wire.Bind(new(service.S3BucketDeleter), new(*S3BucketDeleter)),
)

var _ service.S3BucketDeleter = (*S3BucketDeleter)(nil)

// NewS3BucketDeleter creates a new S3BucketDeleter.
func NewS3BucketDeleter(client *s3.Client) *S3BucketDeleter {
	return &S3BucketDeleter{client: client}
}

// DeleteS3Bucket deletes the bucket.
func (c *S3BucketDeleter) DeleteS3Bucket(ctx context.Context, input *service.S3BucketDeleterInput) (*service.S3BucketDeleterOutput, error) {
	_, err := c.client.DeleteBucket(ctx,
		&s3.DeleteBucketInput{
			Bucket: aws.String(input.Bucket.String()),
		},
		func(o *s3.Options) {
			o.Region = input.Region.String()
		})
	if err != nil {
		return nil, err
	}
	return &service.S3BucketDeleterOutput{}, nil
}

// S3BucketObjectsDeleter implements the S3BucketObjectsDeleter interface.
type S3BucketObjectsDeleter struct {
	client *s3.Client
}

// S3BucketObjectsDeleterSet is a provider set for S3BucketObjectsDeleter.
//
//nolint:gochecknoglobals
var S3BucketObjectsDeleterSet = wire.NewSet(
	NewS3BucketObjectsDeleter,
	wire.Bind(new(service.S3BucketObjectsDeleter), new(*S3BucketObjectsDeleter)),
)

var _ service.S3BucketObjectsDeleter = (*S3BucketObjectsDeleter)(nil)

// NewS3BucketObjectsDeleter creates a new S3BucketObjectsDeleter.
func NewS3BucketObjectsDeleter(client *s3.Client) *S3BucketObjectsDeleter {
	return &S3BucketObjectsDeleter{client: client}
}

// DeleteS3BucketObjects deletes the objects in the bucket.
func (c *S3BucketObjectsDeleter) DeleteS3BucketObjects(ctx context.Context, input *service.S3BucketObjectsDeleterInput) (*service.S3BucketObjectsDeleterOutput, error) {
	optFn := func(o *s3.Options) {
		o.Retryer = NewRetryer(func(err error) bool {
			return strings.Contains(err.Error(), "api error SlowDown")
		}, model.S3DeleteObjectsDelayTimeSec)
		o.Region = input.Region.String()
	}

	if _, err := c.client.DeleteObjects(
		ctx,
		&s3.DeleteObjectsInput{
			Bucket: aws.String(input.Bucket.String()),
			Delete: &types.Delete{
				Objects: input.S3ObjectSets.ToS3ObjectIdentifiers(),
				Quiet:   aws.Bool(true),
			},
		},
		optFn,
	); err != nil {
		return nil, err
	}
	return &service.S3BucketObjectsDeleterOutput{}, nil
}

// S3BucketObjectsLister implements the S3BucketObjectsLister interface.
type S3BucketObjectsLister struct {
	client *s3.Client
}

// S3BucketObjectsListerSet is a provider set for S3BucketObjectsLister.
//
//nolint:gochecknoglobals
var S3BucketObjectsListerSet = wire.NewSet(
	NewS3BucketObjectsLister,
	wire.Bind(new(service.S3BucketObjectsLister), new(*S3BucketObjectsLister)),
)

var _ service.S3BucketObjectsLister = (*S3BucketObjectsLister)(nil)

// NewS3BucketObjectsLister creates a new S3BucketObjectsLister.
func NewS3BucketObjectsLister(client *s3.Client) *S3BucketObjectsLister {
	return &S3BucketObjectsLister{client: client}
}

// ListS3BucketObjects lists the objects in the bucket.
func (c *S3BucketObjectsLister) ListS3BucketObjects(ctx context.Context, input *service.S3BucketObjectsListerInput) (*service.S3BucketObjectsListerOutput, error) {
	out, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(input.Bucket.String()),
	})
	if err != nil {
		return nil, err
	}

	var objects model.S3ObjectSets
	for _, o := range out.Contents {
		objects = append(objects, model.S3Object{
			S3Key: model.S3Key(*o.Key),
		})
	}
	return &service.S3BucketObjectsListerOutput{Objects: objects}, nil
}
