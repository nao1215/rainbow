package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/domain/service"
)

// CFnStackLister is a mock of the CFnStackLister interface.
type CFnStackLister func(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error)

// ListCFnStack calls the CFnStackListerFunc.
func (m CFnStackLister) ListCFnStack(ctx context.Context, input *service.CFnStackListerInput) (*service.CFnStackListerOutput, error) {
	return m(ctx, input)
}
