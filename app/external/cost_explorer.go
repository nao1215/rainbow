package external

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

// NewCostExplorerClient returns a new CostExplorer client.
func NewCostExplorerClient(cfg *model.AWSConfig) *costexplorer.Client {
	return costexplorer.NewFromConfig(*cfg.Config)
}

// CostGetter is an interface for getting cost.
type CostGetter struct {
	*costexplorer.Client
}

// CostGetterSet is a provider set for CostGetter.
//
//nolint:gochecknoglobals
var CostGetterSet = wire.NewSet(
	NewCostGetter,
	wire.Bind(new(service.CostGetter), new(*CostGetter)),
)

var _ service.CostGetter = (*CostGetter)(nil)

// NewCostGetter creates a new CostGetter.
func NewCostGetter(c *costexplorer.Client) *CostGetter {
	return &CostGetter{Client: c}
}

// GetCost gets the cost.
func (c *CostGetter) GetCost(ctx context.Context, input *service.CostGetterInput) (*service.CostGetterOutput, error) {
	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(input.Start.Format("2006-01-02")),
			End:   aws.String(input.End.Format("2006-01-02")),
		},
		Granularity: types.GranularityDaily,
		Metrics:     []string{"UnblendedCost"},
	}

	resp, err := c.GetCostAndUsage(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(resp.ResultsByTime) == 0 || len(resp.ResultsByTime[0].Total) == 0 {
		return nil, fmt.Errorf("no cost data available for the specified time period")
	}

	return &service.CostGetterOutput{Cost: *resp.ResultsByTime[0].Total["UnblendedCost"].Amount}, nil
}
