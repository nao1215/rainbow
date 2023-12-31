package model

import (
	"encoding/json"
	"fmt"

	"github.com/nao1215/rainbow/utils/errfmt"
)

// Statement is a type that represents a statement.
type Statement struct {
	// Sid is an identifier for the statement.
	Sid string `json:"Sid"` //nolint
	// Effect is whether the statement allows or denies access.
	Effect string `json:"Effect"` //nolint
	// Principal is the AWS account, IAM user, IAM role, federated user, or assumed-role user that the statement applies to.
	Principal Principal `json:"Principal"` //nolint
	// Action is the specific action or actions that will be allowed or denied.
	Action []string `json:"Action"` //nolint
	// Resource is the specific Amazon S3 resources that the statement covers.
	Resource []string `json:"Resource"` //nolint
	// The Condition element (or Condition block) lets you specify conditions for when a policy is in effect.
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition.html
	Condition map[string]map[string]string `json:"Condition,omitempty"` //nolint
}

// Principal is a type that represents a principal.
type Principal struct {
	// Service is the AWS service to which the principal belongs.
	Service string `json:"Service"` //nolint
}

// BucketPolicy is a type that represents a bucket policy.
type BucketPolicy struct {
	// Version is the policy language version.
	Version string `json:"Version"` //nolint
	// Statement is the policy statement.
	Statement []Statement `json:"Statement"` //nolint
}

// NewAllowCloudFrontS3BucketPolicy returns a new BucketPolicy that allows CloudFront to access the S3 bucket.
func NewAllowCloudFrontS3BucketPolicy(bucket Bucket) *BucketPolicy {
	return &BucketPolicy{
		Version: "2012-10-17",
		Statement: []Statement{
			{
				Sid:       "Allow CloudFront to GetObject",
				Effect:    "Allow",
				Principal: Principal{Service: "cloudfront.amazonaws.com"},
				Action: []string{
					"s3:GetObject",
					"s3:ListBucket",
				},
				Resource: []string{
					fmt.Sprintf("arn::aws:s3:::%s", bucket.String()),
					fmt.Sprintf("arn::aws:s3:::%s/*", bucket.String()),
				},
			},
			{
				Sid:       "Secure Access",
				Effect:    "Deny",
				Principal: Principal{Service: "*"},
				Action: []string{
					"s3:*",
				},
				Resource: []string{
					fmt.Sprintf("arn::aws:s3:::%s", bucket.String()),
					fmt.Sprintf("arn::aws:s3:::%s/*", bucket.String()),
				},
				Condition: map[string]map[string]string{
					"Bool": {
						"aws:SecureTransport": "false",
					},
				},
			},
		},
	}
}

// String returns the string representation of the BucketPolicy.
func (b *BucketPolicy) String() (string, error) {
	policy, err := json.Marshal(b)
	if err != nil {
		return "", errfmt.Wrap(err, "failed to marshal bucket policy")
	}
	return string(policy), nil
}
