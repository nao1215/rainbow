package service

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
	CFnStackLister(ctx context.Context, input *CFnStackListerInput) (*CFnStackListerOutput, error)
}

// CFnStackResourceListerInput is the input of the CFnStackResourceLister method.
type CFnStackResourceListerInput struct {
	// StackName is the name of the stack.
	StackName string
	// Region is the region of the stack.
	Region model.Region
}

// CFnStackResourceListerOutput is the output of the CFnStackResourceLister method.
type CFnStackResourceListerOutput struct {
	// Resources is a list of CloudFormation stack resources.
	Resources []*model.StackResource
}

// CFnStackResourceLister is the interface that wraps the basic CFnStackResourceLister method.
type CFnStackResourceLister interface {
	CFnStackResourceLister(ctx context.Context, input *CFnStackResourceListerInput) (*CFnStackResourceListerOutput, error)
}

// CFnStackCreatorInput is the input of the CFnStackCreator method.
type CFnStackCreatorInput struct {
	// StackName is the name of the stack.
	StackName string
	// TemplateBody is the template body.
	TemplateBody string
}

// CFnStackCreatorOutput is the output of the CFnStackCreator method.
type CFnStackCreatorOutput struct {
	// StackID is the ID of the stack.
	StackID string
}

// CFnStackCreator is the interface that wraps the basic CFnStackCreator method.
type CFnStackCreator interface {
	CFnStackCreator(ctx context.Context, input *CFnStackCreatorInput) (*CFnStackCreatorOutput, error)
}

// CFnStackDeleterInput is the input of the CFnStackDeleter method.
type CFnStackDeleterInput struct {
	// StackName is the name of the stack.
	StackName string
}

// CFnStackDeleterOutput is the output of the CFnStackDeleter method.
type CFnStackDeleterOutput struct{}

// CFnStackDeleter is the interface that wraps the basic CFnStackDeleter method.
type CFnStackDeleter interface {
	CFnStackDeleter(ctx context.Context, input *CFnStackDeleterInput) (*CFnStackDeleterOutput, error)
}
