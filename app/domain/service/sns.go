package service

import "context"

// SNSPublisherInput is the input of the Publish method.
type SNSPublisherInput struct {
	// TopicArn is the ARN of the topic.
	TopicArn string
	// Message is the message that you want to publish.
	Message string
}

// SNSPublisherOutput is the output of the Publish method.
type SNSPublisherOutput struct{}

// SNSPublisher is the interface that wraps the basic Publish method.
type SNSPublisher interface {
	PublishSNS(ctx context.Context, input *SNSPublisherInput) (*SNSPublisherOutput, error)
}
