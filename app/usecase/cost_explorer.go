package usecase

import (
	"context"
	"time"
)

// CostGetterInput is the input of the GetCost method.
type CostGetterInput struct {
	// Start is the start date of the period.
	Start time.Time
	// End is the end date of the period.
	End time.Time
}

// CostGetterOutput is the output of the GetCost method.
type CostGetterOutput struct {
	// Cost is the cost of the period. The unit is USD.
	Cost string
}

// CostGetter is the interface that wraps the basic GetCost method.
type CostGetter interface {
	GetCost(ctx context.Context, input *CostGetterInput) (*CostGetterOutput, error)
}
