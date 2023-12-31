package service

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// CloudFrontCreatorInput is an input struct for CDNCreator.
type CloudFrontCreatorInput struct {
	// BucketName is the name of the  bucket.
	Bucket model.Bucket
	// OAIID is the ID of the OAI.
	OAIID *string
}

// CloudFrontCreatorOutput is an output struct for CDNCreator.
type CloudFrontCreatorOutput struct {
	// Domain is the domain of the CDN.
	Domain model.Domain
}

// CloudFrontCreator is an interface for creating CDN.
type CloudFrontCreator interface {
	CreateCloudFront(context.Context, *CloudFrontCreatorInput) (*CloudFrontCreatorOutput, error)
}

// OAICreatorInput is an input struct for OAICreator.
type OAICreatorInput struct{}

// OAICreatorOutput is an output struct for OAICreator.
type OAICreatorOutput struct {
	// ID is the ID of the OAI.
	ID *string
}

// OAICreator is an interface for creating OAI.
// OAI is an Origin Access Identity.
type OAICreator interface {
	CreateOAI(context.Context, *OAICreatorInput) (*OAICreatorOutput, error)
}
