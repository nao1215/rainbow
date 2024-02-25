package interactor

import (
	"context"

	"github.com/google/wire"
	"github.com/nao1215/rainbow/app/domain/service"
	"github.com/nao1215/rainbow/app/usecase"
)

// CFnStackListerSet is a set of CFnStackLister.
//
//nolint:gochecknoglobals
var CFnStackListerSet = wire.NewSet(
	NewCFnStackLister,
	wire.Bind(new(usecase.CFnStackLister), new(*CFnStackLister)),
)

var _ usecase.CFnStackLister = (*CFnStackLister)(nil)

// CFnStackLister is an implementation for CFnStackLister.
type CFnStackLister struct {
	service.CFnStackLister
}

// NewCFnStackLister returns a new CFnStackLister struct.
func NewCFnStackLister(lister service.CFnStackLister) *CFnStackLister {
	return &CFnStackLister{
		CFnStackLister: lister,
	}
}

// ListCFnStack returns a list of CloudFormation stacks.
func (l *CFnStackLister) ListCFnStack(ctx context.Context, input *usecase.CFnStackListerInput) (*usecase.CFnStackListerOutput, error) {
	output, err := l.CFnStackLister.ListCFnStack(ctx, &service.CFnStackListerInput{
		Region: input.Region,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.CFnStackListerOutput{
		Stacks: output.Stacks,
	}, nil
}

// CFnStackEventsDescriberSet is a set of CFnStackEventsDescriber.
//
//nolint:gochecknoglobals
var CFnStackEventsDescriberSet = wire.NewSet(
	NewCFnStackEventsDescriber,
	wire.Bind(new(usecase.CFnStackEventsDescriber), new(*CFnStackEventsDescriber)),
)

var _ usecase.CFnStackEventsDescriber = (*CFnStackEventsDescriber)(nil)

// CFnStackEventsDescriber is an implementation for CFnStackEventsDescriber.
type CFnStackEventsDescriber struct {
	service.CFnStackEventsDescriber
}

// NewCFnStackEventsDescriber returns a new CFnStackEventsDescriber struct.
func NewCFnStackEventsDescriber(describer service.CFnStackEventsDescriber) *CFnStackEventsDescriber {
	return &CFnStackEventsDescriber{
		CFnStackEventsDescriber: describer,
	}
}

// DescribeCFnStackEvents returns a list of CloudFormation stack events.
func (d *CFnStackEventsDescriber) DescribeCFnStackEvents(ctx context.Context, input *usecase.CFnStackEventsDescriberInput) (*usecase.CFnStackEventsDescriberOutput, error) {
	output, err := d.CFnStackEventsDescriber.DescribeCFnStackEvents(ctx, &service.CFnStackEventsDescriberInput{
		StackName: input.StackName,
		Region:    input.Region,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.CFnStackEventsDescriberOutput{
		Events: output.Events,
	}, nil
}
