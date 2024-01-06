package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/service"
)

// S3BucketCreator is a mock of the S3BucketCreator interface.
type S3BucketCreator func(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error)

// CreateS3Bucket calls the CreateS3BucketFunc.
func (m S3BucketCreator) CreateS3Bucket(ctx context.Context, input *service.S3BucketCreatorInput) (*service.S3BucketCreatorOutput, error) {
	return m(ctx, input)
}

// S3BucketLister is a mock of the S3BucketLister interface.
type S3BucketLister func(ctx context.Context, input *service.S3BucketListerInput) (*service.S3BucketListerOutput, error)

// ListS3Buckets calls the ListS3BucketsFunc.
func (m S3BucketLister) ListS3Buckets(ctx context.Context, input *service.S3BucketListerInput) (*service.S3BucketListerOutput, error) {
	return m(ctx, input)
}

// S3BucketLocationGetter is a mock of the S3BucketLocationGetter interface.
type S3BucketLocationGetter func(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error)

// GetS3BucketLocation calls the GetS3BucketLocationFunc.
func (m S3BucketLocationGetter) GetS3BucketLocation(ctx context.Context, input *service.S3BucketLocationGetterInput) (*service.S3BucketLocationGetterOutput, error) {
	return m(ctx, input)
}

// S3ObjectsLister is a mock of the S3ObjectLister interface.
type S3ObjectsLister func(ctx context.Context, input *service.S3ObjectsListerInput) (*service.S3ObjectsListerOutput, error)

// ListS3Objects calls the ListS3ObjectsFunc.
func (m S3ObjectsLister) ListS3Objects(ctx context.Context, input *service.S3ObjectsListerInput) (*service.S3ObjectsListerOutput, error) {
	return m(ctx, input)
}

// S3ObjectsDeleter is a mock of the S3ObjectsDeleter interface.
type S3ObjectsDeleter func(ctx context.Context, input *service.S3ObjectsDeleterInput) (*service.S3ObjectsDeleterOutput, error)

// DeleteS3Objects calls the DeleteS3ObjectsFunc.
func (m S3ObjectsDeleter) DeleteS3Objects(ctx context.Context, input *service.S3ObjectsDeleterInput) (*service.S3ObjectsDeleterOutput, error) {
	return m(ctx, input)
}

// S3BucketDeleter is a mock of the S3BucketDeleter interface.
type S3BucketDeleter func(ctx context.Context, input *service.S3BucketDeleterInput) (*service.S3BucketDeleterOutput, error)

// DeleteS3Bucket calls the DeleteS3BucketFunc.
func (m S3BucketDeleter) DeleteS3Bucket(ctx context.Context, input *service.S3BucketDeleterInput) (*service.S3BucketDeleterOutput, error) {
	return m(ctx, input)
}

// S3ObjectUploader is a mock of the S3ObjectUploader interface.
type S3ObjectUploader func(ctx context.Context, input *service.S3ObjectUploaderInput) (*service.S3ObjectUploaderOutput, error)

// UploadS3Object calls the UploadS3ObjectFunc.
func (m S3ObjectUploader) UploadS3Object(ctx context.Context, input *service.S3ObjectUploaderInput) (*service.S3ObjectUploaderOutput, error) {
	return m(ctx, input)
}

// S3BucketPublicAccessBlocker is a mock of the S3BucketPublicAccessBlocker interface.
type S3BucketPublicAccessBlocker func(ctx context.Context, input *service.S3BucketPublicAccessBlockerInput) (*service.S3BucketPublicAccessBlockerOutput, error)

// BlockS3BucketPublicAccess calls the BlockS3BucketPublicAccessFunc.
func (m S3BucketPublicAccessBlocker) BlockS3BucketPublicAccess(ctx context.Context, input *service.S3BucketPublicAccessBlockerInput) (*service.S3BucketPublicAccessBlockerOutput, error) {
	return m(ctx, input)
}

// S3BucketPolicySetter is a mock of the S3BucketPolicySetter interface.
type S3BucketPolicySetter func(ctx context.Context, input *service.S3BucketPolicySetterInput) (*service.S3BucketPolicySetterOutput, error)

// SetS3BucketPolicy calls the SetS3BucketPolicyFunc.
func (m S3BucketPolicySetter) SetS3BucketPolicy(ctx context.Context, input *service.S3BucketPolicySetterInput) (*service.S3BucketPolicySetterOutput, error) {
	return m(ctx, input)
}

// S3ObjectDownloader is a mock of the S3ObjectDownloader interface.
type S3ObjectDownloader func(ctx context.Context, input *service.S3ObjectDownloaderInput) (*service.S3ObjectDownloaderOutput, error)

// DownloadS3Object calls the DownloadS3ObjectFunc.
func (m S3ObjectDownloader) DownloadS3Object(ctx context.Context, input *service.S3ObjectDownloaderInput) (*service.S3ObjectDownloaderOutput, error) {
	return m(ctx, input)
}

// S3ObjectCopier is a mock of the S3ObjectCopier interface.
type S3ObjectCopier func(ctx context.Context, input *service.S3ObjectCopierInput) (*service.S3ObjectCopierOutput, error)

// CopyS3Object calls the CopyS3ObjectFunc.
func (m S3ObjectCopier) CopyS3Object(ctx context.Context, input *service.S3ObjectCopierInput) (*service.S3ObjectCopierOutput, error) {
	return m(ctx, input)
}
