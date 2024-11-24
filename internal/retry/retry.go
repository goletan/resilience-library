// /resilience/retry/retry.go
package retry

import (
	"context"
	"math/rand"
	"time"

	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

// RetryPolicy holds settings for retry behavior.
type RetryPolicy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
	ShouldRetry    func(error) bool
	Logger         *zap.Logger
}

var _ types.RetryPolicyInterface = (*RetryPolicy)(nil)

// NewRetryPolicy initializes a new RetryPolicy based on the configuration.
func NewRetryPolicy(cfg *types.ResilienceConfig, logger *zap.Logger) *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:     cfg.Retry.MaxRetries,
		InitialBackoff: cfg.Retry.InitialBackoff,
		MaxBackoff:     cfg.Retry.MaxBackoff,
		BackoffFactor:  cfg.Retry.BackoffFactor,
		ShouldRetry:    func(err error) bool { return true }, // Default retry policy
		Logger:         logger,
	}
}

// ExecuteWithRetry retries a function with exponential backoff and jitter.
func (rp *RetryPolicy) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	backoff := rp.InitialBackoff

	for attempt := 0; attempt < rp.MaxRetries; attempt++ {
		start := time.Now()
		err := operation()
		if err == nil {
			// Operation succeeded
			CountRetryAttempt("operation_name", "success") // Update success metric
			TrackRetryLatency("operation_name", time.Since(start))
			return nil
		}

		// Log retry attempt
		rp.Logger.Warn("Operation failed, retrying...", zap.Error(err), zap.Int("attempt", attempt+1))

		// Check if the error is retryable based on custom logic
		if !rp.ShouldRetry(err) {
			rp.Logger.Warn("Non-retryable error occurred", zap.Error(err))
			CountRetryAttempt("operation_name", "failure") // Update failure metric
			return err
		}

		// Calculate exponential backoff with jitter
		waitTime := rp.calculateBackoffWithJitter(backoff)
		rp.Logger.Warn("Retry attempt with backoff", zap.Int("attempt", attempt+1), zap.Duration("wait_time", waitTime))

		// Dynamic retry delay using context
		retryCtx, cancel := context.WithTimeout(ctx, waitTime)
		defer cancel()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-retryCtx.Done():
			// Continue to next retry after backoff
		}

		// Ensure the backoff does not exceed maxBackoff
		backoff = min(backoff*time.Duration(rp.BackoffFactor), rp.MaxBackoff)
	}

	CountRetryAttempt("operation_name", "exceeded") // Update metric for retries exhausted
	return ctx.Err()
}

// calculateBackoffWithJitter returns a randomized backoff duration using jitter.
func (rp *RetryPolicy) calculateBackoffWithJitter(baseBackoff time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(baseBackoff)))
	return baseBackoff + jitter
}

// min returns the smaller of two time durations.
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
