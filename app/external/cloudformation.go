package external

import (
	"context"

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

// CloudFormationStackLister implements the CloudFormationStackLister interface.
type CloudFormationStackLister struct {
	client *cloudformation.Client
}

// CloudFormationStackListerSet is a set of CloudFormationStackLister.
//
//nolint:gochecknoglobals
var CloudFormationStackListerSet = wire.NewSet(
	NewCloudFormationStackLister,
	wire.Bind(new(service.CloudFormationStackLister), new(*CloudFormationStackLister)),
)

var _ service.CloudFormationStackLister = (*CloudFormationStackLister)(nil)

// NewCloudFormationStackLister returns a new CloudFormationStackLister.
func NewCloudFormationStackLister(client *cloudformation.Client) *CloudFormationStackLister {
	return &CloudFormationStackLister{client: client}
}

// CloudFormationStackLister returns a list of CloudFormation stacks.
func (l *CloudFormationStackLister) CloudFormationStackLister(ctx context.Context, input *service.CloudFormationStackListerInput) (*service.CloudFormationStackListerOutput, error) {
	in := &cloudformation.ListStacksInput{}
	stacks := make([]*model.Stack, 0, 100)

	for {
		select {
		case <-ctx.Done():
			return &service.CloudFormationStackListerOutput{
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

	return &service.CloudFormationStackListerOutput{
		Stacks: stacks,
	}, nil
}
