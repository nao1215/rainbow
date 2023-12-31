package interactor

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/service"

	"github.com/nao1215/rainbow/app/usecase"
)

// CloudFrontCreatorSet is a set of CloudFrontCreator.
//
//nolint:gochecknoglobals
var CloudFrontCreatorSet = wire.NewSet(
	NewCloudFrontCreator,
	wire.Struct(new(CloudFrontCreatorOptions), "*"),
	wire.Bind(new(usecase.CloudFrontCreator), new(*CloudFrontCreator)),
)

var _ usecase.CloudFrontCreator = (*CloudFrontCreator)(nil)

// CloudFrontCreator is an implementation for CloudFrontCreator.
type CloudFrontCreator struct {
	opts *CloudFrontCreatorOptions
}

// CloudFrontCreatorOptions is an option struct for CloudFrontCreator.
type CloudFrontCreatorOptions struct {
	service.CloudFrontCreator
	service.OAICreator
}

// NewCloudFrontCreator returns a new CloudFrontCreator struct.
func NewCloudFrontCreator(opts *CloudFrontCreatorOptions) *CloudFrontCreator {
	return &CloudFrontCreator{
		opts: opts,
	}
}

// CreateCloudFront creates a CDN.
func (c *CloudFrontCreator) CreateCloudFront(ctx context.Context, input *usecase.CreateCloudFrontInput) (*usecase.CreateCloudFrontOutput, error) {
	oaiOutput, err := c.opts.OAICreator.CreateOAI(ctx, &service.OAICreatorInput{})
	if err != nil {
		return nil, err
	}

	createCDNOutput, err := c.opts.CloudFrontCreator.CreateCloudFront(ctx, &service.CloudFrontCreatorInput{
		Bucket: input.Bucket,
		OAIID:  oaiOutput.ID,
	})
	if err != nil {
		return nil, err
	}

	return &usecase.CreateCloudFrontOutput{
		Domain: createCDNOutput.Domain,
	}, nil
}
