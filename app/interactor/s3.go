// Package interactor contains the implementations of usecases.
package interactor

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
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

// S3ObjectsLister implements the S3ObjectsLister interface.
type S3ObjectsLister struct {
	service.S3ObjectsLister
}

// S3ObjectsListerSet is a provider set for S3ObjectsLister.
//
//nolint:gochecknoglobals
var S3ObjectsListerSet = wire.NewSet(
	NewS3ObjectsLister,
	wire.Bind(new(usecase.S3ObjectsLister), new(*S3ObjectsLister)),
)

var _ usecase.S3ObjectsLister = (*S3ObjectsLister)(nil)

// NewS3ObjectsLister creates a new S3ObjectsLister.
func NewS3ObjectsLister(l service.S3ObjectsLister) *S3ObjectsLister {
	return &S3ObjectsLister{
		S3ObjectsLister: l,
	}
}

// ListS3Objects lists the objects in the S3 bucket.
func (s *S3ObjectsLister) ListS3Objects(ctx context.Context, input *usecase.S3ObjectsListerInput) (*usecase.S3ObjectsListerOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	out, err := s.S3ObjectsLister.ListS3Objects(ctx, &service.S3ObjectsListerInput{
		Bucket: input.Bucket,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.S3ObjectsListerOutput{
		Objects: out.Objects,
	}, nil
}

// S3ObjectsDeleter implements the S3ObjectsDeleter interface.
type S3ObjectsDeleter struct {
	service.S3ObjectsDeleter
	service.S3BucketLocationGetter
	service.S3ObjectVersionsLister
}

// S3ObjectsDeleterSet is a provider set for S3ObjectsDeleter.
//
//nolint:gochecknoglobals
var S3ObjectsDeleterSet = wire.NewSet(
	NewS3ObjectsDeleter,
	wire.Bind(new(usecase.S3ObjectsDeleter), new(*S3ObjectsDeleter)),
)

var _ usecase.S3ObjectsDeleter = (*S3ObjectsDeleter)(nil)

// NewS3ObjectsDeleter creates a new S3ObjectsDeleter.
func NewS3ObjectsDeleter(
	d service.S3ObjectsDeleter,
	g service.S3BucketLocationGetter,
	l service.S3ObjectVersionsLister,
) *S3ObjectsDeleter {
	return &S3ObjectsDeleter{
		S3ObjectsDeleter:       d,
		S3BucketLocationGetter: g,
		S3ObjectVersionsLister: l,
	}
}

// DeleteS3Objects deletes the objects in the bucket.
func (s *S3ObjectsDeleter) DeleteS3Objects(ctx context.Context, input *usecase.S3ObjectsDeleterInput) (*usecase.S3ObjectsDeleterOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	location, err := s.S3BucketLocationGetter.GetS3BucketLocation(ctx, &service.S3BucketLocationGetterInput{
		Bucket: input.Bucket,
	})
	if err != nil {
		return nil, err
	}

	versions, err := s.S3ObjectVersionsLister.ListS3ObjectVersions(ctx, &service.S3ObjectVersionsListerInput{
		Bucket: input.Bucket,
	})
	if err != nil {
		return nil, err
	}
	if len(versions.Objects) == 0 {
		return &usecase.S3ObjectsDeleterOutput{}, nil // no objects to delete
	}

	targets := make(model.S3ObjectIdentifiers, 0, len(versions.Objects))
	versionMap := make(map[model.S3Key][]model.VersionID)
	for _, version := range versions.Objects {
		versionMap[version.S3Key] = append(versionMap[version.S3Key], version.VersionID)
	}

	for _, inputIdentifier := range input.S3ObjectIdentifiers {
		if versionIDs, ok := versionMap[inputIdentifier.S3Key]; ok {
			for _, versionID := range versionIDs {
				targets = append(targets, model.S3ObjectIdentifier{
					S3Key:     inputIdentifier.S3Key,
					VersionID: versionID,
				})
			}
		}
	}

	if _, err = s.S3ObjectsDeleter.DeleteS3Objects(ctx, &service.S3ObjectsDeleterInput{
		Bucket:       input.Bucket,
		Region:       location.Region,
		S3ObjectSets: targets,
	}); err != nil {
		return nil, err
	}

	return &usecase.S3ObjectsDeleterOutput{}, nil
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

// FileUploaderSet is a provider set for FileUploader.
//
//nolint:gochecknoglobals
var FileUploaderSet = wire.NewSet(
	NewFileUploader,
	wire.Bind(new(usecase.FileUploader), new(*FileUploader)),
)

var _ usecase.FileUploader = (*FileUploader)(nil)

// FileUploader is an implementation for FileUploader.
type FileUploader struct {
	service.S3ObjectUploader
}

// NewFileUploader returns a new FileUploader struct.
func NewFileUploader(uploader service.S3ObjectUploader) *FileUploader {
	return &FileUploader{
		S3ObjectUploader: uploader,
	}
}

// UploadFile uploads a file to external storage.
func (u *FileUploader) UploadFile(ctx context.Context, input *usecase.FileUploaderInput) (*usecase.FileUploaderOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}
	if err := input.Region.Validate(); err != nil {
		return nil, err
	}

	output, err := u.S3ObjectUploader.UploadS3Object(ctx, &service.S3ObjectUploaderInput{
		Bucket:   input.Bucket,
		Region:   input.Region,
		S3Key:    input.Key,
		S3Object: model.NewS3Object(input.Data),
	})
	if err != nil {
		return nil, err
	}

	return &usecase.FileUploaderOutput{
		ContentType:   output.ContentType,
		ContentLength: output.ContentLength,
	}, nil
}

// S3BucketPublicAccessBlockerSet is a provider set for BucketPublicAccessBlocker.
//
//nolint:gochecknoglobals
var S3BucketPublicAccessBlockerSet = wire.NewSet(
	NewS3BucketPublicAccessBlocker,
	wire.Bind(new(usecase.S3BucketPublicAccessBlocker), new(*S3BucketPublicAccessBlocker)),
)

// S3BucketPublicAccessBlocker is an implementation for BucketPublicAccessBlocker.
type S3BucketPublicAccessBlocker struct {
	service.S3BucketPublicAccessBlocker
}

var _ usecase.S3BucketPublicAccessBlocker = (*S3BucketPublicAccessBlocker)(nil)

// NewS3BucketPublicAccessBlocker returns a new S3BucketPublicAccessBlocker struct.
func NewS3BucketPublicAccessBlocker(b service.S3BucketPublicAccessBlocker) *S3BucketPublicAccessBlocker {
	return &S3BucketPublicAccessBlocker{
		S3BucketPublicAccessBlocker: b,
	}
}

// BlockS3BucketPublicAccess blocks public access to a bucket on S3.
func (s *S3BucketPublicAccessBlocker) BlockS3BucketPublicAccess(ctx context.Context, input *usecase.S3BucketPublicAccessBlockerInput) (*usecase.S3BucketPublicAccessBlockerOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}
	if err := input.Region.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.S3BucketPublicAccessBlocker.BlockS3BucketPublicAccess(ctx, &service.S3BucketPublicAccessBlockerInput{
		Bucket: input.Bucket,
		Region: input.Region,
	}); err != nil {
		return nil, err
	}
	return &usecase.S3BucketPublicAccessBlockerOutput{}, nil
}

// S3BucketPolicySetterSet is a provider set for BucketPolicySetter.
//
//nolint:gochecknoglobals
var S3BucketPolicySetterSet = wire.NewSet(
	NewS3BucketPolicySetter,
	wire.Bind(new(usecase.S3BucketPolicySetter), new(*S3BucketPolicySetter)),
)

// S3BucketPolicySetter is an implementation for BucketPolicySetter.
type S3BucketPolicySetter struct {
	service.S3BucketPolicySetter
}

var _ usecase.S3BucketPolicySetter = (*S3BucketPolicySetter)(nil)

// NewS3BucketPolicySetter returns a new S3BucketPolicySetter struct.
func NewS3BucketPolicySetter(s service.S3BucketPolicySetter) *S3BucketPolicySetter {
	return &S3BucketPolicySetter{
		S3BucketPolicySetter: s,
	}
}

// SetS3BucketPolicy sets a bucket policy on S3.
func (s *S3BucketPolicySetter) SetS3BucketPolicy(ctx context.Context, input *usecase.S3BucketPolicySetterInput) (*usecase.S3BucketPolicySetterOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.S3BucketPolicySetter.SetS3BucketPolicy(ctx, &service.S3BucketPolicySetterInput{
		Bucket: input.Bucket,
		Policy: input.Policy,
	}); err != nil {
		return nil, err
	}
	return &usecase.S3BucketPolicySetterOutput{}, nil
}

// S3ObjectDownloaderSet is a provider set for S3ObjectDownloader.
//
//nolint:gochecknoglobals
var S3ObjectDownloaderSet = wire.NewSet(
	NewS3ObjectDownloader,
	wire.Bind(new(usecase.S3ObjectDownloader), new(*S3ObjectDownloader)),
)

// S3ObjectDownloader is an implementation for S3ObjectDownloader.
type S3ObjectDownloader struct {
	service.S3ObjectDownloader
}

var _ usecase.S3ObjectDownloader = (*S3ObjectDownloader)(nil)

// NewS3ObjectDownloader returns a new S3ObjectDownloader struct.
func NewS3ObjectDownloader(d service.S3ObjectDownloader) *S3ObjectDownloader {
	return &S3ObjectDownloader{
		S3ObjectDownloader: d,
	}
}

// DownloadS3Object downloads an object from S3.
func (s *S3ObjectDownloader) DownloadS3Object(ctx context.Context, input *usecase.S3ObjectDownloaderInput) (*usecase.S3ObjectDownloaderOutput, error) {
	if err := input.Bucket.Validate(); err != nil {
		return nil, err
	}

	out, err := s.S3ObjectDownloader.DownloadS3Object(ctx, &service.S3ObjectDownloaderInput{
		Bucket: input.Bucket,
		Key:    input.Key,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.S3ObjectDownloaderOutput{
		Bucket:        out.Bucket,
		Key:           out.Key,
		ContentType:   out.ContentType,
		ContentLength: out.ContentLength,
		S3Object:      out.S3Object,
	}, nil
}

// S3ObjectCopierSet is a provider set for S3ObjectCopier.
//
//nolint:gochecknoglobals
var S3ObjectCopierSet = wire.NewSet(
	NewS3ObjectCopier,
	wire.Bind(new(usecase.S3ObjectCopier), new(*S3ObjectCopier)),
)

// S3ObjectCopier is an implementation for S3ObjectCopier.
type S3ObjectCopier struct {
	service.S3ObjectCopier
}

var _ usecase.S3ObjectCopier = (*S3ObjectCopier)(nil)

// NewS3ObjectCopier returns a new S3ObjectCopier struct.
func NewS3ObjectCopier(c service.S3ObjectCopier) *S3ObjectCopier {
	return &S3ObjectCopier{
		S3ObjectCopier: c,
	}
}

// CopyS3Object copies an object from S3 to S3.
func (s *S3ObjectCopier) CopyS3Object(ctx context.Context, input *usecase.S3ObjectCopierInput) (*usecase.S3ObjectCopierOutput, error) {
	if err := input.SourceBucket.Validate(); err != nil {
		return nil, err
	}
	if err := input.DestinationBucket.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.S3ObjectCopier.CopyS3Object(ctx, &service.S3ObjectCopierInput{
		SourceBucket:      input.SourceBucket,
		SourceKey:         input.SourceKey,
		DestinationBucket: input.DestinationBucket,
		DestinationKey:    input.DestinationKey,
	}); err != nil {
		return nil, err
	}
	return &usecase.S3ObjectCopierOutput{}, nil
}
