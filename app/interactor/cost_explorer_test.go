package interactor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/external/mock"
	"github.com/nao1215/rainbow/app/usecase"
)

func TestCostGetter_GetCost(t *testing.T) {
	t.Run("Success to get cost", func(t *testing.T) {
		t.Parallel()

		costGetter := mock.CostGetter(func(_ context.Context, input *service.CostGetterInput) (*service.CostGetterOutput, error) {
			want := &service.CostGetterInput{
				Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			if diff := cmp.Diff(want, input); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return &service.CostGetterOutput{Cost: "1000"}, nil
		})

		getter := NewCostGetter(costGetter)
		got, err := getter.GetCost(context.Background(), &usecase.CostGetterInput{
			Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 31, 23, 59, 59, 0, time.UTC),
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := &usecase.CostGetterOutput{Cost: "1000"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("Fail to get cost", func(t *testing.T) {
		t.Parallel()

		costGetter := mock.CostGetter(func(_ context.Context, _ *service.CostGetterInput) (*service.CostGetterOutput, error) {
			want := &service.CostGetterInput{
				Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 31, 23, 59, 59, 0, time.UTC),
			}
			if diff := cmp.Diff(want, want); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
			return nil, errors.New("failed to get cost")
		})

		getter := NewCostGetter(costGetter)
		_, err := getter.GetCost(context.Background(), &usecase.CostGetterInput{
			Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 31, 23, 59, 59, 0, time.UTC),
		})
		if err == nil {
			t.Error("expected error, but not occurred")
		}
	})
}
