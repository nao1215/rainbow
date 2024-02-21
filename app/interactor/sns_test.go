package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/external/mock"
	"github.com/nao1215/rainbow/app/usecase"
)

func TestSNSPublisher_PublishSNS(t *testing.T) {
	t.Parallel()

	t.Run("success to publish sns", func(t *testing.T) {
		t.Parallel()

		snsPublisher := mock.SNSPublisher(func(_ context.Context, input *service.SNSPublisherInput) (*service.SNSPublisherOutput, error) {
			want := &service.SNSPublisherInput{
				Message:  "message",
				TopicArn: "topicArn",
			}

			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.SNSPublisherOutput{}, nil
		})

		publisher := NewSNSPublisher(snsPublisher)
		if _, err := publisher.PublishSNS(context.Background(), &usecase.SNSPublisherInput{
			Message:  "message",
			TopicArn: "topicArn",
		}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("fail to publish sns", func(t *testing.T) {
		t.Parallel()

		snsPublisher := mock.SNSPublisher(func(_ context.Context, input *service.SNSPublisherInput) (*service.SNSPublisherOutput, error) {
			want := &service.SNSPublisherInput{
				Message:  "message",
				TopicArn: "topicArn",
			}

			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return nil, errors.New("failed to publish sns")
		})

		publisher := NewSNSPublisher(snsPublisher)
		if _, err := publisher.PublishSNS(context.Background(), &usecase.SNSPublisherInput{
			Message:  "message",
			TopicArn: "topicArn",
		}); err == nil {
			t.Error("expected error, but not occurred")
		}
	})
}
