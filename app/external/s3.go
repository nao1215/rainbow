// Package external implements the external service.
package external

import (
	"context"
	"fmt"

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

// CreateBucket creates a new S3 bucket.
func (c *S3BucketCreator) CreateBucket(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
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

// ListBuckets lists the buckets.
func (c *S3BucketLister) ListBuckets(ctx context.Context, _ *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
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

// GetBucketLocation gets the location of the bucket.
func (c *S3BucketLocationGetter) GetBucketLocation(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
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
