package usecase

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// S3BucketPublicAccessBlockerInput is the input of the BlockBucketPublicAccess method.
type S3BucketPublicAccessBlockerInput struct {
	// Bucket is the name of the  bucket.
	Bucket model.Bucket
	// Region is the name of the region.
	Region model.Region
}

// S3BucketPublicAccessBlockerOutput is the output of the BlockBucketPublicAccess method.
type S3BucketPublicAccessBlockerOutput struct{}

// S3BucketPublicAccessBlocker is the interface that wraps the basic BlockBucketPublicAccess method.
type S3BucketPublicAccessBlocker interface {
	BlockS3BucketPublicAccess(ctx context.Context, input *S3BucketPublicAccessBlockerInput) (*S3BucketPublicAccessBlockerOutput, error)
}

// S3BucketPolicySetterInput is the input of the SetBucketPolicy method.
type S3BucketPolicySetterInput struct {
	// Bucket is the name of the  bucket.
	Bucket model.Bucket
	// Policy is the policy to set.
	Policy *model.BucketPolicy
}

// S3BucketPolicySetterOutput is an output struct for BucketPolicySetter.
type S3BucketPolicySetterOutput struct{}

// S3BucketPolicySetter is an interface for setting a bucket policy.
type S3BucketPolicySetter interface {
	SetS3BucketPolicy(context.Context, *S3BucketPolicySetterInput) (*S3BucketPolicySetterOutput, error)
}
