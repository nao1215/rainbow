package interactor

import (
	"context"
	"errors"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/usecase"
)

// CostGetterSet is a set of CostGetter.
//
//nolint:gochecknoglobals
var CostGetterSet = wire.NewSet(
	NewCostGetter,
	wire.Bind(new(usecase.CostGetter), new(*CostGetter)),
)

var _ usecase.CostGetter = (*CostGetter)(nil)

// CostGetter is an implementation for CostGetter.
type CostGetter struct {
	service.CostGetter
}

// NewCostGetter returns a new CostGetter struct.
func NewCostGetter(c service.CostGetter) *CostGetter {
	return &CostGetter{CostGetter: c}
}

// GetCost gets the cost.
func (c *CostGetter) GetCost(ctx context.Context, input *usecase.CostGetterInput) (*usecase.CostGetterOutput, error) {
	if input.End.Before(input.Start) {
		return nil, errors.New("end date is before the start date")
	}
	output, err := c.CostGetter.GetCost(ctx, &service.CostGetterInput{
		Start: input.Start,
		End:   input.End,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.CostGetterOutput{Cost: output.Cost}, nil
}
