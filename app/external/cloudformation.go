package external

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/domain/service"
)

// NewCloudFormationClient returns a new CloudFormation client.
func NewCloudFormationClient(cfg *model.AWSConfig) *cloudformation.Client {
	return cloudformation.NewFromConfig(
		*cfg.Config,
		func(o *cloudformation.Options) {
			o.RetryMaxAttempts = model.CloudFormationRetryMaxAttempts
			o.RetryMode = aws.RetryModeStandard
		})
}

// NewCloudFormationStackDeleteCompleteWaiter returns a new CloudFormationStackDeleteCompleteWaiter.
func NewCloudFormationStackDeleteCompleteWaiter(client *cloudformation.Client) *cloudformation.StackDeleteCompleteWaiter {
	return cloudformation.NewStackDeleteCompleteWaiter(client)
}

// CFnStackLister implements the CFnStackLister interface.
type CFnStackLister struct {
	client *cloudformation.Client
}

// CFnStackListerSet is a set of CFnStackLister.
//
//nolint:gochecknoglobals
var CFnStackListerSet = wire.NewSet(
	NewCFnStackLister,
	wire.Bind(new(service.CFnStackLister), new(*CFnStackLister)),
)

var _ service.CFnStackLister = (*CFnStackLister)(nil)

// NewCFnStackLister returns a new CloudFormationStackLister.
func NewCFnStackLister(client *cloudformation.Client) *CFnStackLister {
	return &CFnStackLister{client: client}
}

// CFnStackLister returns a list of CloudFormation stacks.
func (l *CFnStackLister) CFnStackLister(ctx context.Context, _ *service.CFnStackListerInput) (*service.CFnStackListerOutput, error) {
	in := &cloudformation.ListStacksInput{}
	stacks := make([]*model.Stack, 0, 100)

	for {
		select {
		case <-ctx.Done():
			return &service.CFnStackListerOutput{
				Stacks: stacks,
			}, ctx.Err()
		default:
		}

		out, err := l.client.ListStacks(ctx, in)
		if err != nil {
			return nil, err
		}

		for _, stack := range out.StackSummaries {
			stacks = append(stacks, &model.Stack{
				CreationTime: stack.CreationTime,
				StackName:    stack.StackName,
				StackStatus:  model.StackStatus(stack.StackStatus),
				DeletionTime: stack.DeletionTime,
				DriftInformation: &model.StackDriftInformationSummary{
					StackDriftStatus:   model.StackDriftStatus(stack.DriftInformation.StackDriftStatus),
					LastCheckTimestamp: stack.DriftInformation.LastCheckTimestamp,
				},
				LastUpdatedTime:     stack.LastUpdatedTime,
				ParentID:            stack.ParentId,
				RootID:              stack.RootId,
				StackID:             stack.StackId,
				StackStatusReason:   stack.StackStatusReason,
				TemplateDescription: stack.TemplateDescription,
			})
		}

		if out.NextToken == nil {
			break
		}
		in.NextToken = out.NextToken
	}

	return &service.CFnStackListerOutput{
		Stacks: stacks,
	}, nil
}

// CFnStackResourceLister implements the CFnStackResourceLister interface.
type CFnStackResourceLister struct {
	client *cloudformation.Client
}

// CFnStackResourceListerSet is a set of CFnStackResourceLister.
//
//nolint:gochecknoglobals
var CFnStackResourceListerSet = wire.NewSet(
	NewCFnStackResourceLister,
	wire.Bind(new(service.CFnStackResourceLister), new(*CFnStackResourceLister)),
)

var _ service.CFnStackResourceLister = (*CFnStackResourceLister)(nil)

// NewCFnStackResourceLister returns a new CloudFormationStackResourceLister.
func NewCFnStackResourceLister(client *cloudformation.Client) *CFnStackResourceLister {
	return &CFnStackResourceLister{client: client}
}

// CFnStackResourceLister returns a list of CloudFormation stack resources.
func (l *CFnStackResourceLister) CFnStackResourceLister(ctx context.Context, input *service.CFnStackResourceListerInput) (*service.CFnStackResourceListerOutput, error) {
	in := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(input.StackName),
	}
	resources := make([]*model.StackResource, 0, 100)

	for {
		select {
		case <-ctx.Done():
			return &service.CFnStackResourceListerOutput{
				Resources: resources,
			}, ctx.Err()
		default:
		}

		out, err := l.client.ListStackResources(ctx, in)
		if err != nil {
			return nil, err
		}

		for _, resource := range out.StackResourceSummaries {
			resources = append(resources, &model.StackResource{
				LastUpdatedTimestamp: resource.LastUpdatedTimestamp,
				LogicalResourceID:    resource.LogicalResourceId,
				ResourceStatus:       model.ResourceStatus(resource.ResourceStatus),
				ResourceType:         resource.ResourceType,
				DriftInformation: &model.StackResourceDriftInformationSummary{
					StackResourceDriftStatus: model.StackResourceDriftStatus(resource.DriftInformation.StackResourceDriftStatus),
					LastCheckTimestamp:       resource.DriftInformation.LastCheckTimestamp,
				},
				PhysicalResourceID:   resource.PhysicalResourceId,
				ResourceStatusReason: resource.ResourceStatusReason,
			})
		}
		if out.NextToken == nil {
			break
		}
		in.NextToken = out.NextToken
	}

	return &service.CFnStackResourceListerOutput{
		Resources: resources,
	}, nil
}

// CFnStackCreator implements the CFnStackCreator interface.
type CFnStackCreator struct {
	client *cloudformation.Client
}

// CFnStackCreatorSet is a set of CFnStackCreator.
//
//nolint:gochecknoglobals
var CFnStackCreatorSet = wire.NewSet(
	NewCFnStackCreator,
	wire.Bind(new(service.CFnStackCreator), new(*CFnStackCreator)),
)

var _ service.CFnStackCreator = (*CFnStackCreator)(nil)

// NewCFnStackCreator returns a new CloudFormationStackCreator.
func NewCFnStackCreator(client *cloudformation.Client) *CFnStackCreator {
	return &CFnStackCreator{client: client}
}

// CFnStackCreator creates a CloudFormation stack.
func (c *CFnStackCreator) CFnStackCreator(ctx context.Context, input *service.CFnStackCreatorInput) (*service.CFnStackCreatorOutput, error) {
	in := &cloudformation.CreateStackInput{
		StackName:    aws.String(input.StackName),
		TemplateBody: aws.String(input.TemplateBody),
	}
	out, err := c.client.CreateStack(ctx, in)
	if err != nil {
		return nil, err
	}

	return &service.CFnStackCreatorOutput{
		StackID: *out.StackId,
	}, nil
}

// CFnStackDeleter implements the CFnStackDeleter interface.
type CFnStackDeleter struct {
	client *cloudformation.Client
	waiter *cloudformation.StackDeleteCompleteWaiter
}

// CFnStackDeleterSet is a set of CFnStackDeleter.
//
//nolint:gochecknoglobals
var CFnStackDeleterSet = wire.NewSet(
	NewCFnStackDeleter,
	wire.Bind(new(service.CFnStackDeleter), new(*CFnStackDeleter)),
)

var _ service.CFnStackDeleter = (*CFnStackDeleter)(nil)

// NewCFnStackDeleter returns a new CloudFormationStackDeleter.
func NewCFnStackDeleter(client *cloudformation.Client, waiter *cloudformation.StackDeleteCompleteWaiter) *CFnStackDeleter {
	return &CFnStackDeleter{
		client: client,
		waiter: waiter,
	}
}

// CFnStackDeleter deletes a CloudFormation stack.
func (d *CFnStackDeleter) CFnStackDeleter(ctx context.Context, input *service.CFnStackDeleterInput) (*service.CFnStackDeleterOutput, error) {
	in := &cloudformation.DeleteStackInput{
		StackName: aws.String(input.StackName),
	}
	_, err := d.client.DeleteStack(ctx, in)
	if err != nil {
		return nil, err
	}

	if err = d.waiter.Wait(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(input.StackName),
	}, model.CloudFormationWaitNanoSecTime); err != nil && !strings.Contains(err.Error(), "waiter state transitioned to Failure") {
		return nil, err
	}

	return &service.CFnStackDeleterOutput{}, nil
}
