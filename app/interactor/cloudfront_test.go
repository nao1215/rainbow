package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/external/mock"
	"github.com/nao1215/rainbow/app/usecase"
)

func TestCloudFrontCreator_CreateCloudFront(t *testing.T) {
	t.Parallel()

	t.Run("success to create CloudFront", func(t *testing.T) {
		t.Parallel()

		oaiCreator := mock.OAICreator(func(ctx context.Context, input *service.OAICreatorInput) (*service.OAICreatorOutput, error) {
			return &service.OAICreatorOutput{
				ID: aws.String("test-oai-id"),
			}, nil
		})

		cloudFrontCreator := mock.CloudFrontCreator(func(ctx context.Context, input *service.CloudFrontCreatorInput) (*service.CloudFrontCreatorOutput, error) {
			want := &service.CloudFrontCreatorInput{
				Bucket: model.Bucket("test-bucket"),
				OAIID:  aws.String("test-oai-id"),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.CloudFrontCreatorOutput{
				Domain: model.Domain("test.cloudfront.net"),
			}, nil
		})

		creator := NewCloudFrontCreator(&CloudFrontCreatorOptions{
			OAICreator:        oaiCreator,
			CloudFrontCreator: cloudFrontCreator,
		})

		got, err := creator.CreateCloudFront(context.Background(), &usecase.CreateCloudFrontInput{
			Bucket: model.Bucket("test-bucket"),
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := &usecase.CreateCloudFrontOutput{
			Domain: model.Domain("test.cloudfront.net"),
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("fail to create CloudFront", func(t *testing.T) {
		t.Parallel()

		oaiCreator := mock.OAICreator(func(ctx context.Context, input *service.OAICreatorInput) (*service.OAICreatorOutput, error) {
			return &service.OAICreatorOutput{
				ID: aws.String("test-oai-id"),
			}, nil
		})

		cloudFrontCreator := mock.CloudFrontCreator(func(ctx context.Context, input *service.CloudFrontCreatorInput) (*service.CloudFrontCreatorOutput, error) {
			want := &service.CloudFrontCreatorInput{
				Bucket: model.Bucket("test-bucket"),
				OAIID:  aws.String("test-oai-id"),
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return nil, errors.New("some error")
		})

		creator := NewCloudFrontCreator(&CloudFrontCreatorOptions{
			OAICreator:        oaiCreator,
			CloudFrontCreator: cloudFrontCreator,
		})

		_, err := creator.CreateCloudFront(context.Background(), &usecase.CreateCloudFrontInput{
			Bucket: model.Bucket("test-bucket"),
		})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("fail to create OAI", func(t *testing.T) {
		t.Parallel()

		oaiCreator := mock.OAICreator(func(ctx context.Context, input *service.OAICreatorInput) (*service.OAICreatorOutput, error) {
			return nil, errors.New("some error")
		})

		creator := NewCloudFrontCreator(&CloudFrontCreatorOptions{
			OAICreator:        oaiCreator,
			CloudFrontCreator: nil,
		})

		if _, err := creator.CreateCloudFront(context.Background(), &usecase.CreateCloudFrontInput{
			Bucket: model.Bucket("test-bucket"),
		}); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("fail to validate s3 bucket name", func(t *testing.T) {
		t.Parallel()

		creator := NewCloudFrontCreator(&CloudFrontCreatorOptions{
			OAICreator:        nil,
			CloudFrontCreator: nil,
		})

		if _, err := creator.CreateCloudFront(context.Background(), &usecase.CreateCloudFrontInput{
			Bucket: model.Bucket(""),
		}); err == nil {
			t.Error("expected error")
		}
	})
}
