package external

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/nao1215/rainbow/app/domain/model"
)

var _ aws.RetryerV2 = (*Retryer)(nil)

// Retryer implements the aws.RetryerV2 interface.
type Retryer struct {
	// isErrorRetryableFunc is a function that determines whether the error is retryable.
	isErrorRetryableFunc func(error) bool
	// delayTimeSec is the delay time in seconds.
	delayTimeSec int
}

// NewRetryer creates a new Retryer.
func NewRetryer(isErrorRetryableFunc func(error) bool, delayTimeSec int) *Retryer {
	return &Retryer{
		isErrorRetryableFunc: isErrorRetryableFunc,
		delayTimeSec:         delayTimeSec,
	}
}

// IsErrorRetryable returns true if the error is retryable.
func (r *Retryer) IsErrorRetryable(err error) bool {
	return r.isErrorRetryableFunc(err)
}

// MaxAttempts returns the maximum number of attempts.
func (r *Retryer) MaxAttempts() int {
	return model.MaxS3DeleteObjectsRetryCount
}

// RetryDelay returns the delay time.
func (r *Retryer) RetryDelay(int, error) (time.Duration, error) {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(r.delayTimeSec)))
	if err != nil {
		return 0, err
	}
	waitTime := 1 + int(randomInt.Int64())
	return time.Duration(waitTime) * time.Second, nil
}

// GetRetryToken returns the retry token. This is not used.
func (r *Retryer) GetRetryToken(context.Context, error) (func(error) error, error) {
	return func(error) error { return nil }, nil
}

// GetInitialToken returns the initial token. This is not used.
func (r *Retryer) GetInitialToken() func(error) error {
	return func(error) error { return nil }
}

// GetAttemptToken returns the attempt token. This is not used.
func (r *Retryer) GetAttemptToken(context.Context) (func(error) error, error) {
	return func(error) error { return nil }, nil
}
