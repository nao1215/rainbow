package model

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// AWSProfile is the name of the AWS profile.
type AWSProfile string

// NewAWSProfile returns a new AWSProfile.
// If p is empty, read $AWS_PROFILE and return it.
func NewAWSProfile(p string) AWSProfile {
	if p == "" {
		profile := os.Getenv("AWS_PROFILE")
		if profile == "" {
			return AWSProfile("default")
		}
		return AWSProfile(profile)
	}
	return AWSProfile(p)
}

// String returns the string representation of the AWSProfile.
func (p AWSProfile) String() string {
	return string(p)
}

// AWSConfig is the AWS config.
type AWSConfig struct {
	*aws.Config
}

// NewAWSConfig creates a new AWS config.
func NewAWSConfig(ctx context.Context, profile AWSProfile, region Region) (*AWSConfig, error) {
	opts := []func(*config.LoadOptions) error{}
	if profile.String() != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile.String()))
	}
	if region.String() != "" {
		opts = append(opts, config.WithRegion(string(region)))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &AWSConfig{
		Config: &cfg,
	}, nil
}

// Region returns the AWS region.
func (c *AWSConfig) Region() Region {
	if Region(c.Config.Region) == "" {
		return RegionUSEast1
	}
	return Region(c.Config.Region)
}
