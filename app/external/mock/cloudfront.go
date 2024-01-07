package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/service"
)

// CloudFrontCreator is a mock of the CloudFrontCreator interface.
type CloudFrontCreator func(ctx context.Context, input *service.CloudFrontCreatorInput) (*service.CloudFrontCreatorOutput, error)

// CreateCloudFront calls the CreateCloudFrontFunc.
func (m CloudFrontCreator) CreateCloudFront(ctx context.Context, input *service.CloudFrontCreatorInput) (*service.CloudFrontCreatorOutput, error) {
	return m(ctx, input)
}

// OAICreator is a mock of the OAICreator interface.
type OAICreator func(ctx context.Context, input *service.OAICreatorInput) (*service.OAICreatorOutput, error)

// CreateOAI calls the CreateOAIFunc.
func (m OAICreator) CreateOAI(ctx context.Context, input *service.OAICreatorInput) (*service.OAICreatorOutput, error) {
	return m(ctx, input)
}
