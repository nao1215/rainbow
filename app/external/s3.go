// Package external implements the external service.
package external

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

// NewS3Client creates a new S3 service client.
// If profile is empty, the default profile is used.
func NewS3Client(cfg *model.AWSConfig) *s3.Client {
	return s3.NewFromConfig(*cfg.Config)
}

// S3BucketCreator implements the S3BucketCreator interface.
type S3BucketCreator struct {
	*s3.Client
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
	return &S3BucketCreator{Client: client}
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

	_, err := c.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket:                    aws.String(input.Bucket.String()),
		CreateBucketConfiguration: locationContstraint,
	})
	if err != nil {
		var alreadyExistsErr *types.BucketAlreadyExists
		var alreadyOwnedByYouErr *types.BucketAlreadyOwnedByYou
		if errors.As(err, &alreadyExistsErr) {
			return nil, fmt.Errorf("%w: region=%s, bucket name=%s", domain.ErrBucketAlreadyExistsOwnedByOther, input.Region.String(), input.Bucket.String())
		}
		if errors.As(err, &alreadyOwnedByYouErr) {
			return nil, fmt.Errorf("%w: region=%s, bucket name=%s", domain.ErrBucketAlreadyOwnedByYou, input.Region.String(), input.Bucket.String())
		}
		return nil, fmt.Errorf("%w: region=%s, bucket name=%s", err, input.Region.String(), input.Bucket.String())
	}
	return &service.S3BucketCreatorOutput{}, nil
}

// S3BucketLister implements the S3BucketLister interface.
type S3BucketLister struct {
	*s3.Client
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
	return &S3BucketLister{Client: client}
}

// ListS3Buckets lists the buckets.
func (c *S3BucketLister) ListS3Buckets(ctx context.Context, _ *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
	out, err := c.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	buckets := make(model.BucketSets, 0, len(out.Buckets))
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
	*s3.Client
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
	return &S3BucketLocationGetter{Client: client}
}

// GetS3BucketLocation gets the location of the bucket.
func (c *S3BucketLocationGetter) GetS3BucketLocation(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
	out, err := c.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(input.Bucket.String()),
	})
	if err != nil {
		var noSuchBucket *types.NoSuchBucket
		if errors.As(err, &noSuchBucket) {
			return nil, fmt.Errorf("%w: bucket name=%s", domain.ErrNoSuchBucket, input.Bucket.String())
		}
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
	*s3.Client
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
	return &S3BucketDeleter{Client: client}
}

// DeleteS3Bucket deletes the bucket.
func (c *S3BucketDeleter) DeleteS3Bucket(ctx context.Context, input *service.S3BucketDeleterInput) (*service.S3BucketDeleterOutput, error) {
	_, err := c.DeleteBucket(ctx,
		&s3.DeleteBucketInput{
			Bucket: aws.String(input.Bucket.String()),
		},
		func(o *s3.Options) {
			o.Region = input.Region.String()
		})
	if err != nil {
		var noSuchBucket *types.NoSuchBucket
		if errors.As(err, &noSuchBucket) {
			return nil, fmt.Errorf("%w: bucket name=%s", domain.ErrNoSuchBucket, input.Bucket.String())
		}
		return nil, err
	}
	return &service.S3BucketDeleterOutput{}, nil
}

// S3ObjectsDeleter implements the S3ObjectsDeleter interface.
type S3ObjectsDeleter struct {
	*s3.Client
}

// S3ObjectsDeleterSet is a provider set for S3ObjectsDeleter.
//
//nolint:gochecknoglobals
var S3ObjectsDeleterSet = wire.NewSet(
	NewS3ObjectsDeleter,
	wire.Bind(new(service.S3ObjectsDeleter), new(*S3ObjectsDeleter)),
)

var _ service.S3ObjectsDeleter = (*S3ObjectsDeleter)(nil)

// NewS3ObjectsDeleter creates a new S3ObjectsDeleter.
func NewS3ObjectsDeleter(client *s3.Client) *S3ObjectsDeleter {
	return &S3ObjectsDeleter{Client: client}
}

// DeleteS3Objects deletes the objects in the bucket.
// If the bucket has versioning enabled, all versions of the specified objects are deleted.
func (c *S3ObjectsDeleter) DeleteS3Objects(ctx context.Context, input *service.S3ObjectsDeleterInput) (*service.S3ObjectsDeleterOutput, error) {
	optFn := func(o *s3.Options) {
		o.Retryer = NewRetryer(func(err error) bool {
			return strings.Contains(err.Error(), "api error SlowDown")
		}, model.S3DeleteObjectsDelayTimeSec)
		o.Region = input.Region.String()
	}

	if _, err := c.DeleteObjects(
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
	return &service.S3ObjectsDeleterOutput{}, nil
}

// S3ObjectsLister implements the S3ObjectsLister interface.
type S3ObjectsLister struct {
	*s3.Client
}

// S3ObjectsListerSet is a provider set for S3ObjectsLister.
//
//nolint:gochecknoglobals
var S3ObjectsListerSet = wire.NewSet(
	NewS3ObjectsLister,
	wire.Bind(new(service.S3ObjectsLister), new(*S3ObjectsLister)),
)

var _ service.S3ObjectsLister = (*S3ObjectsLister)(nil)

// NewS3ObjectsLister creates a new S3ObjectsLister.
func NewS3ObjectsLister(client *s3.Client) *S3ObjectsLister {
	return &S3ObjectsLister{Client: client}
}

// ListS3Objects lists the objects in the bucket.
func (c *S3ObjectsLister) ListS3Objects(ctx context.Context, input *service.S3ObjectsListerInput) (*service.S3ObjectsListerOutput, error) {
	var objects model.S3ObjectIdentifiers
	in := &s3.ListObjectsV2Input{
		Bucket:  aws.String(input.Bucket.String()),
		MaxKeys: aws.Int32(model.MaxS3Keys),
	}
	for {
		output, err := c.ListObjectsV2(ctx, in)
		if err != nil {
			return nil, err
		}

		for _, o := range output.Contents {
			objects = append(objects, model.S3ObjectIdentifier{
				S3Key: model.S3Key(*o.Key),
			})
		}

		if !*output.IsTruncated {
			break
		}
		in.ContinuationToken = output.NextContinuationToken
	}
	return &service.S3ObjectsListerOutput{Objects: objects}, nil
}

// S3ObjectDownloader implements the S3ObjectDownloader interface.
type S3ObjectDownloader struct {
	*s3.Client
}

// S3ObjectDownloaderSet is a provider set for S3ObjectGetter.
//
//nolint:gochecknoglobals
var S3ObjectDownloaderSet = wire.NewSet(
	NewS3ObjectDownloader,
	wire.Bind(new(service.S3ObjectDownloader), new(*S3ObjectDownloader)),
)

var _ service.S3ObjectDownloader = (*S3ObjectDownloader)(nil)

// NewS3ObjectDownloader creates a new S3ObjectGetter.
func NewS3ObjectDownloader(client *s3.Client) *S3ObjectDownloader {
	return &S3ObjectDownloader{Client: client}
}

// DownloadS3Object gets the object in the bucket.
func (c *S3ObjectDownloader) DownloadS3Object(ctx context.Context, input *service.S3ObjectDownloaderInput) (*service.S3ObjectDownloaderOutput, error) {
	out, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(input.Bucket.String()),
		Key:    aws.String(input.Key.String()),
	})
	if err != nil {
		return nil, err
	}

	body := out.Body
	defer func() {
		e := body.Close()
		if e != nil {
			err = errors.Join(err, e)
		}
	}()

	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return &service.S3ObjectDownloaderOutput{
		Bucket:        input.Bucket,
		Key:           input.Key,
		ContentType:   aws.ToString(out.ContentType),
		ContentLength: aws.ToInt64(out.ContentLength),
		S3Object:      model.NewS3Object(b),
	}, nil
}

// S3ObjectUploader implements the S3ObjectUploader interface.
type S3ObjectUploader struct {
	*s3.Client
}

// S3ObjectUploaderSet is a provider set for S3ObjectUploader.
//
//nolint:gochecknoglobals
var S3ObjectUploaderSet = wire.NewSet(
	NewS3ObjectUploader,
	wire.Bind(new(service.S3ObjectUploader), new(*S3ObjectUploader)),
)

var _ service.S3ObjectUploader = (*S3ObjectUploader)(nil)

// NewS3ObjectUploader creates a new S3ObjectUploader.
func NewS3ObjectUploader(client *s3.Client) *S3ObjectUploader {
	return &S3ObjectUploader{Client: client}
}

// UploadS3Object puts the object in the bucket.
func (c *S3ObjectUploader) UploadS3Object(ctx context.Context, input *service.S3ObjectUploaderInput) (*service.S3ObjectUploaderOutput, error) {
	if _, err := c.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:        aws.String(input.Bucket.String()),
			Key:           aws.String(input.S3Key.String()),
			Body:          input.S3Object,
			ContentType:   aws.String(input.S3Object.ContentType()),
			ContentLength: aws.Int64(input.S3Object.ContentLength()),
		},
		s3.WithAPIOptions(
			v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
		)); err != nil {
		return nil, err
	}

	return &service.S3ObjectUploaderOutput{
		ContentType:   input.S3Object.ContentType(),
		ContentLength: input.S3Object.ContentLength(),
	}, nil
}

// S3ObjectCopier implements the S3ObjectCopier interface.
type S3ObjectCopier struct {
	*s3.Client
}

// S3ObjectCopierSet is a provider set for S3ObjectCopier.
//
//nolint:gochecknoglobals
var S3ObjectCopierSet = wire.NewSet(
	NewS3ObjectCopier,
	wire.Bind(new(service.S3ObjectCopier), new(*S3ObjectCopier)),
)

var _ service.S3ObjectCopier = (*S3ObjectCopier)(nil)

// NewS3ObjectCopier creates a new S3ObjectCopier.
func NewS3ObjectCopier(client *s3.Client) *S3ObjectCopier {
	return &S3ObjectCopier{Client: client}
}

// CopyS3Object copies the object in the bucket.
func (c *S3ObjectCopier) CopyS3Object(ctx context.Context, input *service.S3ObjectCopierInput) (*service.S3ObjectCopierOutput, error) {
	_, err := c.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(input.DestinationBucket.String()),
		CopySource: aws.String(input.SourceBucket.Join(input.SourceKey).String()),
		Key:        aws.String(input.DestinationKey.String()),
	})
	if err != nil {
		return nil, err
	}
	return &service.S3ObjectCopierOutput{}, nil
}

// S3ObjectVersionsLister implements the S3ObjectVersionsLister interface.
type S3ObjectVersionsLister struct {
	*s3.Client
}

// S3ObjectVersionsListerSet is a provider set for S3ObjectVersionsLister.
//
//nolint:gochecknoglobals
var S3ObjectVersionsListerSet = wire.NewSet(
	NewS3ObjectVersionsLister,
	wire.Bind(new(service.S3ObjectVersionsLister), new(*S3ObjectVersionsLister)),
)

var _ service.S3ObjectVersionsLister = (*S3ObjectVersionsLister)(nil)

// NewS3ObjectVersionsLister creates a new S3ObjectVersionsLister.
func NewS3ObjectVersionsLister(client *s3.Client) *S3ObjectVersionsLister {
	return &S3ObjectVersionsLister{Client: client}
}

// ListS3ObjectVersions lists the object versions in the bucket.
func (c *S3ObjectVersionsLister) ListS3ObjectVersions(ctx context.Context, input *service.S3ObjectVersionsListerInput) (*service.S3ObjectVersionsListerOutput, error) {
	var objects model.S3ObjectIdentifiers
	listObjectVersionsInput := &s3.ListObjectVersionsInput{
		Bucket:  aws.String(input.Bucket.String()),
		MaxKeys: aws.Int32(model.MaxS3Keys),
	}

	for {
		listObjectVersionsOutput, err := c.ListObjectVersions(ctx, listObjectVersionsInput)
		if err != nil {
			return nil, err
		}
		for _, version := range listObjectVersionsOutput.Versions {
			objects = append(objects, model.S3ObjectIdentifier{
				S3Key:     model.S3Key(*version.Key),
				VersionID: model.VersionID(*version.VersionId),
			})
		}
		for _, deleteMarker := range listObjectVersionsOutput.DeleteMarkers {
			objects = append(objects, model.S3ObjectIdentifier{
				S3Key:     model.S3Key(*deleteMarker.Key),
				VersionID: model.VersionID(*deleteMarker.VersionId),
			})
		}

		if !*listObjectVersionsOutput.IsTruncated {
			break
		}

		listObjectVersionsInput.KeyMarker = listObjectVersionsOutput.NextKeyMarker
		listObjectVersionsInput.VersionIdMarker = listObjectVersionsOutput.NextVersionIdMarker
	}
	return &service.S3ObjectVersionsListerOutput{Objects: objects}, nil
}
