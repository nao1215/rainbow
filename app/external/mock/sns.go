package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/service"
)

// SNSPublisher is a mock of the SNSPublisher interface.
type SNSPublisher func(ctx context.Context, input *service.SNSPublisherInput) (*service.SNSPublisherOutput, error)

// PublishSNS is a mock of the PublishSNS method.
func (m SNSPublisher) PublishSNS(ctx context.Context, input *service.SNSPublisherInput) (*service.SNSPublisherOutput, error) {
	return m(ctx, input)
}
