package usecase

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/model"
)

// CFnStackListerInput is the input of the CFnStackLister method.
type CFnStackListerInput struct {
	// Region is the region of the stack.
	Region model.Region
}

// CFnStackListerOutput is the output of the CFnStackLister method.
type CFnStackListerOutput struct {
	// Stacks is a list of CloudFormation stacks.
	Stacks []*model.Stack
}

// CFnStackLister is the interface that wraps the basic CFnStackLister method.
type CFnStackLister interface {
	ListCFnStack(ctx context.Context, input *CFnStackListerInput) (*CFnStackListerOutput, error)
}

// CFnStackEventsDescriberInput is the input of the CFnStackEventsDescriber method.
type CFnStackEventsDescriberInput struct {
	// StackName is the name of the stack.
	StackName string
	// Region is the region of the stack.
	Region model.Region
}

// CFnStackEventsDescriberOutput is the output of the CFnStackEventsDescriber method.
type CFnStackEventsDescriberOutput struct {
	// Events is a list of CloudFormation stack events.
	Events []*model.StackEvent
}

// CFnStackEventsDescriber is the interface that wraps the basic CFnStackEventsDescriber method.
type CFnStackEventsDescriber interface {
	DescribeCFnStackEvents(ctx context.Context, input *CFnStackEventsDescriberInput) (*CFnStackEventsDescriberOutput, error)
}
