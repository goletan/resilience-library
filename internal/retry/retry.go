package retry

import (
	"context"
	observability "github.com/goletan/observability/pkg"
	"math/rand"
	"time"

	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

// Policy holds settings for retry behavior.
type Policy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
	ShouldRetry    func(error) bool
	obs            *observability.Observability
}

var _ types.RetryPolicyInterface = (*Policy)(nil)

const (
	operationName       = "operation_name"
	successStatus       = "success"
	failureStatus       = "failure"
	exceededRetriesMsg  = "Retry attempts exceeded"
	operationFailedMsg  = "Operation failed, retrying..."
	nonRetryableErrMsg  = "Non-retryable error occurred"
	retryWithBackoffMsg = "Retry attempt with backoff"
)

// NewRetryPolicy initializes a new Policy based on the configuration.
func NewRetryPolicy(cfg *types.ResilienceConfig, obs *observability.Observability) *Policy {
	return &Policy{
		MaxRetries:     cfg.Retry.MaxRetries,
		InitialBackoff: cfg.Retry.InitialBackoff,
		MaxBackoff:     cfg.Retry.MaxBackoff,
		BackoffFactor:  cfg.Retry.BackoffFactor,
		ShouldRetry:    func(err error) bool { return true }, // Default retry policy
		obs:            obs,
	}
}

func (rp *Policy) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	currentBackoff := rp.InitialBackoff

	for attempt := 0; attempt < rp.MaxRetries; attempt++ {
		if err := rp.tryOperation(ctx, operation, &currentBackoff, attempt); err == nil {
			return nil
		}
	}

	CountAttempt(operationName, exceededRetriesMsg)
	return ctx.Err()
}

func (rp *Policy) tryOperation(ctx context.Context, operation func() error, currentBackoff *time.Duration, attempt int) error {
	start := time.Now()
	err := operation()
	if err == nil {
		CountAttempt(operationName, successStatus)
		TrackLatency(operationName, time.Since(start))
		return nil
	}

	rp.logRetryAttempt(err, attempt)

	if !rp.ShouldRetry(err) {
		rp.logFailure(err)
		return err
	}

	waitTime := rp.calculateBackoffWithJitter(*currentBackoff)
	rp.logBackoff(attempt, waitTime)

	if err := rp.handleRetry(ctx, waitTime); err != nil {
		return err
	}

	*currentBackoff = duration(*currentBackoff*time.Duration(rp.BackoffFactor), rp.MaxBackoff)
	return nil
}

func (rp *Policy) logRetryAttempt(err error, attempt int) {
	rp.obs.Logger.Warn(operationFailedMsg, zap.Error(err), zap.Int("attempt", attempt+1))
}

func (rp *Policy) logFailure(err error) {
	rp.obs.Logger.Warn(nonRetryableErrMsg, zap.Error(err))
	CountAttempt(operationName, failureStatus)
}

func (rp *Policy) logBackoff(attempt int, waitTime time.Duration) {
	rp.obs.Logger.Warn(retryWithBackoffMsg, zap.Int("attempt", attempt+1), zap.Duration("wait_time", waitTime))
}

func (rp *Policy) handleRetry(ctx context.Context, waitTime time.Duration) error {
	retryCtx, cancel := context.WithTimeout(ctx, waitTime)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-retryCtx.Done():
		return nil
	}
}

// calculateBackoffWithJitter returns a randomized backoff duration using jitter.
func (rp *Policy) calculateBackoffWithJitter(baseBackoff time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(baseBackoff)))
	return baseBackoff + jitter
}

// duration returns the smaller of two time durations.
func duration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
