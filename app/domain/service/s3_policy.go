package service

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// S3BucketPublicAccessBlockerInput is an input struct for BucketAccessBlocker.
type S3BucketPublicAccessBlockerInput struct {
	// Bucket is the name of the  bucket.
	Bucket model.Bucket
	// Region is the name of the region.
	Region model.Region
}

// S3BucketPublicAccessBlockerOutput is an output struct for BucketAccessBlocker.
type S3BucketPublicAccessBlockerOutput struct{}

// S3BucketPublicAccessBlocker is an interface for blocking access to a bucket.
type S3BucketPublicAccessBlocker interface {
	BlockS3BucketPublicAccess(context.Context, *S3BucketPublicAccessBlockerInput) (*S3BucketPublicAccessBlockerOutput, error)
}

// S3BucketPolicySetterInput is an input struct for BucketPolicySetter.
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
