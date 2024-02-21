package interactor

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/usecase"
)

// SNSPublisherSet is a set of SNSPublisher.
//
//nolint:gochecknoglobals
var SNSPublisherSet = wire.NewSet(
	NewSNSPublisher,
	wire.Bind(new(usecase.SNSPublisher), new(*SNSPublisher)),
)

var _ usecase.SNSPublisher = (*SNSPublisher)(nil)

// SNSPublisher is an implementation for SNSPublisher.
type SNSPublisher struct {
	service.SNSPublisher
}

// NewSNSPublisher returns a new SNSPublisher struct.
func NewSNSPublisher(s service.SNSPublisher) *SNSPublisher {
	return &SNSPublisher{SNSPublisher: s}
}

// PublishSNS publishes a message to SNS.
func (s *SNSPublisher) PublishSNS(ctx context.Context, input *usecase.SNSPublisherInput) (*usecase.SNSPublisherOutput, error) {
	if _, err := s.SNSPublisher.PublishSNS(ctx, &service.SNSPublisherInput{
		Message:  input.Message,
		TopicArn: input.TopicArn,
	}); err != nil {
		return nil, err
	}
	return &usecase.SNSPublisherOutput{}, nil
}
