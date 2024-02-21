package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/external"
	"github.com/nao1215/rainbow/app/interactor"
	"github.com/nao1215/rainbow/app/usecase"
)

// CostNotifier is a notifier for cost.
type CostNotifier struct {
	costGetter   usecase.CostGetter
	snsPublisher usecase.SNSPublisher
	publishArn   string
}

// NewCostNotifier returns a new CostNotifier.
func NewCostNotifier(ctx context.Context) (*CostNotifier, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	costGetter := external.NewCostGetter(external.NewCostExplorerClient(&model.AWSConfig{
		Config: &cfg,
	}))

	snsPublisher := external.NewSNSPublisher(external.NewSNSClient(&model.AWSConfig{
		Config: &cfg,
	}))

	return &CostNotifier{
		costGetter:   interactor.NewCostGetter(costGetter),
		snsPublisher: interactor.NewSNSPublisher(snsPublisher),
		publishArn:   os.Getenv("SNS_TOPIC_ARN"),
	}, nil
}

// getDailyCost gets the daily cost.
func (c *CostNotifier) getDailyCost(ctx context.Context) (string, error) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	output, err := c.costGetter.GetCost(ctx, &usecase.CostGetterInput{
		Start: yesterday,
		End:   today,
	})
	if err != nil {
		return "", err
	}
	return output.Cost, nil
}

// getMonthlyCost gets the monthly cost.
func (c *CostNotifier) getMonthlyCost(ctx context.Context) (string, error) {
	today := time.Now()
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	output, err := c.costGetter.GetCost(ctx, &usecase.CostGetterInput{
		Start: startOfMonth,
		End:   today,
	})
	if err != nil {
		return "", err
	}
	return output.Cost, nil
}

// handler is the Lambda function handler
func handler(ctx context.Context) error {
	notifier, err := NewCostNotifier(ctx)
	if err != nil {
		return err
	}

	dailyCost, err := notifier.getDailyCost(ctx)
	if err != nil {
		return err
	}

	monthlyCost, err := notifier.getMonthlyCost(ctx)
	if err != nil {
		return err
	}

	if _, err := notifier.snsPublisher.PublishSNS(ctx, &usecase.SNSPublisherInput{
		Message:  fmt.Sprintf("Daily cost %s USD (Monthly cost %s USD)\n", dailyCost, monthlyCost),
		TopicArn: notifier.publishArn,
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
