package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/service"
)

// CFnStackLister is a mock of the CFnStackLister interface.
type CFnStackLister func(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error)

// CFnStackLister calls the CFnStackListerFunc.
func (m CFnStackLister) CFnStackLister(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error) {
	return m(ctx, input)
}
