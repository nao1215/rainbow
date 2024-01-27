package mock

import (
	"context"

	"github.com/nao1215/rainbow/app/usecase"
)

// S3ObjectsLister is a mock of the S3ObjectLister interface.
type S3ObjectsLister func(ctx context.Context, input *usecase.S3ObjectsListerInput) (*usecase.S3ObjectsListerOutput, error)

// ListS3Objects calls the ListS3ObjectsFunc.
func (m S3ObjectsLister) ListS3Objects(ctx context.Context, input *usecase.S3ObjectsListerInput) (*usecase.S3ObjectsListerOutput, error) {
	return m(ctx, input)
}
