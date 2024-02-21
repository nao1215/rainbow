package external

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/utils/errfmt"
)

// S3BucketPublicAccessBlockerSet is a provider set for BucketPublicAccessBlocker.
//
//nolint:gochecknoglobals
var S3BucketPublicAccessBlockerSet = wire.NewSet(
	NewS3BucketPublicAccessBlocker,
	wire.Bind(new(service.S3BucketPublicAccessBlocker), new(*S3BucketPublicAccessBlocker)),
)

// S3BucketPublicAccessBlocker is an implementation for BucketPublicAccessBlocker.
type S3BucketPublicAccessBlocker struct {
	*s3.Client
}

var _ service.S3BucketPublicAccessBlocker = &S3BucketPublicAccessBlocker{}

// NewS3BucketPublicAccessBlocker returns a new S3BucketPublicAccessBlocker struct.
func NewS3BucketPublicAccessBlocker(client *s3.Client) *S3BucketPublicAccessBlocker {
	return &S3BucketPublicAccessBlocker{client}
}

// BlockS3BucketPublicAccess blocks public access to a bucket on S3.
func (s *S3BucketPublicAccessBlocker) BlockS3BucketPublicAccess(ctx context.Context, input *service.S3BucketPublicAccessBlockerInput) (*service.S3BucketPublicAccessBlockerOutput, error) {
	if _, err := s.PutPublicAccessBlock(ctx, &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(input.Bucket.String()),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(true),
			BlockPublicPolicy:     aws.Bool(true),
			IgnorePublicAcls:      aws.Bool(true),
			RestrictPublicBuckets: aws.Bool(true),
		},
	}); err != nil {
		return nil, errfmt.Wrap(domain.ErrBucketPublicAccessBlock, err.Error())
	}
	return &service.S3BucketPublicAccessBlockerOutput{}, nil
}

// S3BucketPolicySetterSet is a provider set for BucketPolicySetter.
//
//nolint:gochecknoglobals
var S3BucketPolicySetterSet = wire.NewSet(
	NewS3BucketPolicySetter,
	wire.Bind(new(service.S3BucketPolicySetter), new(*S3BucketPolicySetter)),
)

// S3BucketPolicySetter is an implementation for BucketPolicySetter.
type S3BucketPolicySetter struct {
	*s3.Client
}

var _ service.S3BucketPolicySetter = &S3BucketPolicySetter{}

// NewS3BucketPolicySetter returns a new S3BucketPolicySetter struct.
func NewS3BucketPolicySetter(client *s3.Client) *S3BucketPolicySetter {
	return &S3BucketPolicySetter{Client: client}
}

// SetS3BucketPolicy sets a bucket policy on S3.
func (s *S3BucketPolicySetter) SetS3BucketPolicy(ctx context.Context, input *service.S3BucketPolicySetterInput) (*service.S3BucketPolicySetterOutput, error) {
	policy, err := input.Policy.String()
	if err != nil {
		return nil, err
	}

	if _, err = s.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(input.Bucket.String()),
		Policy: aws.String(policy),
	}); err != nil {
		return nil, errfmt.Wrap(domain.ErrBucketPolicySet, err.Error())
	}
	return &service.S3BucketPolicySetterOutput{}, nil
}
