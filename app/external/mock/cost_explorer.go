// Package mock is mock for external package.
package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/service"
)

// CostGetter is a mock of the CostGetter interface.
type CostGetter func(ctx context.Context, input *service.CostGetterInput) (*service.CostGetterOutput, error)

// GetCost calls the GetCostFunc.
func (m CostGetter) GetCost(ctx context.Context, input *service.CostGetterInput) (*service.CostGetterOutput, error) {
	return m(ctx, input)
}
