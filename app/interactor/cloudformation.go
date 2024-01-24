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
	output, err := l.CFnStackLister.CFnStackLister(ctx, &service.CFnStackListerInput{
		Region: input.Region,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.CFnStackListerOutput{
		Stacks: output.Stacks,
	}, nil
}
