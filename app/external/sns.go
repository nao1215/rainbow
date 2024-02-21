package external

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

// NewSNSClient returns a new SNSClient.
func NewSNSClient(cfg *model.AWSConfig) *sns.Client {
	return sns.NewFromConfig(*cfg.Config)
}

// SNSPublisher is an implementation for SNSPublisher.
type SNSPublisher struct {
	*sns.Client
}

// SNSPublisherSet is a provider set for SNSPublisher.
//
//nolint:gochecknoglobals
var SNSPublisherSet = wire.NewSet(
	NewSNSPublisher,
	wire.Bind(new(service.SNSPublisher), new(*SNSPublisher)),
)

// NewSNSPublisher creates a new SNSPublisher.
func NewSNSPublisher(c *sns.Client) *SNSPublisher {
	return &SNSPublisher{
		Client: c,
	}
}

var _ service.SNSPublisher = (*SNSPublisher)(nil)

// PublishSNS publishes a message to SNS.
func (p *SNSPublisher) PublishSNS(ctx context.Context, input *service.SNSPublisherInput) (*service.SNSPublisherOutput, error) {
	if _, err := p.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(input.Message),
		TopicArn: aws.String(input.TopicArn),
	}); err != nil {
		return nil, err
	}

	return &service.SNSPublisherOutput{}, nil
}
