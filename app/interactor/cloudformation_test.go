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

func TestCFnStackLister_ListCFnStack(t *testing.T) {
	t.Parallel()

	t.Run("success to get cloudformation stack list", func(t *testing.T) {
		t.Parallel()

		stackLister := mock.CFnStackLister(func(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error) {
			want := &service.CFnStackListerInput{
				Region: model.RegionAPEast1,
			}
			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}

			return &service.CFnStackListerOutput{
				Stacks: []*model.Stack{
					{
						StackName: aws.String("stackName1"),
						StackID:   aws.String(model.StackStatusCreateComplete.String()),
					},
					{
						StackName: aws.String("stackName2"),
						StackID:   aws.String(model.StackStatusCreateComplete.String()),
					},
				},
			}, nil
		})

		lister := NewCFnStackLister(stackLister)
		output, err := lister.ListCFnStack(context.Background(), &usecase.CFnStackListerInput{
			Region: model.RegionAPEast1,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := &usecase.CFnStackListerOutput{
			Stacks: []*model.Stack{
				{
					StackName: aws.String("stackName1"),
					StackID:   aws.String(model.StackStatusCreateComplete.String()),
				},
				{
					StackName: aws.String("stackName2"),
					StackID:   aws.String(model.StackStatusCreateComplete.String()),
				},
			},
		}
		if diff := cmp.Diff(want, output); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("fail to get cloudformation stack list", func(t *testing.T) {
		t.Parallel()

		stackLister := mock.CFnStackLister(func(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error) {
			return nil, errors.New("some error")
		})

		lister := NewCFnStackLister(stackLister)
		_, err := lister.ListCFnStack(context.Background(), &usecase.CFnStackListerInput{
			Region: model.RegionAPEast1,
		})
		if err == nil {
			t.Error("expected error, but nil")
		}
	})
}
