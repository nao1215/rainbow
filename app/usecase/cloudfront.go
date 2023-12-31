package usecase

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// CloudFrontCreator is an interface for creating CloudFront.
type CloudFrontCreator interface {
	CreateCloudFront(ctx context.Context, input *CreateCloudFrontInput) (*CreateCloudFrontOutput, error)
}

// CreateCloudFrontInput is an input struct for CloudFrontCreator.
type CreateCloudFrontInput struct {
	// Bucket is the name of the bucket.
	Bucket model.Bucket
}

// CreateCloudFrontOutput is an output struct for CloudFrontCreator.
type CreateCloudFrontOutput struct {
	// Domain is the domain of the CDN.
	Domain model.Domain
}
