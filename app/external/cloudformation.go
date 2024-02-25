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

// ListCFnStack returns a list of CloudFormation stacks.
func (l *CFnStackLister) ListCFnStack(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error) {
	in := &cloudformation.ListStacksInput{}
	opt := func(o *cloudformation.Options) {
		o.Region = input.Region.String()
	}

	stacks := make([]*model.Stack, 0, 100)
	for {
		select {
		case <-ctx.Done():
			return &service.CFnStackListerOutput{
				Stacks: stacks,
			}, ctx.Err()
		default:
		}

		out, err := l.client.ListStacks(ctx, in, opt)
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

// ListCFnStackResource returns a list of CloudFormation stack resources.
func (l *CFnStackResourceLister) ListCFnStackResource(ctx context.Context, input *service.CFnStackResourceListerInput) (*service.CFnStackResourceListerOutput, error) {
	in := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(input.StackName),
	}
	opt := func(o *cloudformation.Options) {
		o.Region = input.Region.String()
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

		out, err := l.client.ListStackResources(ctx, in, opt)
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

// CreatCFnStack creates a CloudFormation stack.
func (c *CFnStackCreator) CreatCFnStack(ctx context.Context, input *service.CFnStackCreatorInput) (*service.CFnStackCreatorOutput, error) {
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

// DeleteCFnStack deletes a CloudFormation stack.
func (d *CFnStackDeleter) DeleteCFnStack(ctx context.Context, input *service.CFnStackDeleterInput) (*service.CFnStackDeleterOutput, error) {
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

// CFnStackEventsDescriber implements the CFnStackEventsDescriber interface.
type CFnStackEventsDescriber struct {
	client *cloudformation.Client
}

// CFnStackEventsDescriberSet is a set of CFnStackEventsDescriber.
//
//nolint:gochecknoglobals
var CFnStackEventsDescriberSet = wire.NewSet(
	NewCFnStackEventsDescriber,
	wire.Bind(new(service.CFnStackEventsDescriber), new(*CFnStackEventsDescriber)),
)

var _ service.CFnStackEventsDescriber = (*CFnStackEventsDescriber)(nil)

// NewCFnStackEventsDescriber returns a new CloudFormationStackEventsDescriber.
func NewCFnStackEventsDescriber(client *cloudformation.Client) *CFnStackEventsDescriber {
	return &CFnStackEventsDescriber{client: client}
}

// DescribeCFnStackEvents returns a list of CloudFormation stack events.
func (d *CFnStackEventsDescriber) DescribeCFnStackEvents(ctx context.Context, input *service.CFnStackEventsDescriberInput) (*service.CFnStackEventsDescriberOutput, error) {
	in := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(input.StackName),
	}
	opt := func(o *cloudformation.Options) {
		o.Region = input.Region.String()
	}

	out, err := d.client.DescribeStackEvents(ctx, in, opt)
	if err != nil {
		return nil, err
	}

	events := make([]*model.StackEvent, 0, 100)
	for _, event := range out.StackEvents {
		events = append(events, &model.StackEvent{
			EventID:              event.EventId,
			StackID:              event.StackId,
			StackName:            event.StackName,
			Timestamp:            event.Timestamp,
			ClientRequestToken:   event.ClientRequestToken,
			HookFailureMode:      model.HookFailureMode(event.HookFailureMode),
			HookInvocationPoint:  model.HookInvocationPoint(event.HookInvocationPoint),
			HookStatus:           model.HookStatus(event.HookStatus),
			HookStatusReason:     event.HookStatusReason,
			HookType:             event.HookType,
			LogicalResourceID:    event.LogicalResourceId,
			PhysicalResourceID:   event.PhysicalResourceId,
			ResourceProperties:   event.ResourceProperties,
			ResourceStatus:       model.ResourceStatus(event.ResourceStatus),
			ResourceStatusReason: event.ResourceStatusReason,
			ResourceType:         event.ResourceType,
		})
	}
	return &service.CFnStackEventsDescriberOutput{
		Events: events,
	}, nil
}
