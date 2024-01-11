package service

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// CloudFormationStackListerInput is the input of the CloudFormationStackLister method.
type CloudFormationStackListerInput struct{}

// CloudFormationStackListerOutput is the output of the CloudFormationStackLister method.
type CloudFormationStackListerOutput struct {
	// Stacks is a list of CloudFormation stacks.
	Stacks []*model.Stack
}

// CloudFormationStackLister is the interface that wraps the basic CloudFormationStackLister method.
type CloudFormationStackLister interface {
	CloudFormationStackLister(ctx context.Context, input *CloudFormationStackListerInput) (*CloudFormationStackListerOutput, error)
}
