// /resilience/retry/retry.go
package retry

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

// RetryPolicy holds settings for retry behavior.
type RetryPolicy struct {
	MaxRetries  int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
	ShouldRetry func(error) bool
	Logger      *zap.Logger
}

var _ types.RetryPolicyInterface = (*RetryPolicy)(nil)

// NewRetryPolicy initializes a new RetryPolicy with default values.
func NewRetryPolicy(cfg *types.ResilienceConfig, logger *zap.Logger) *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:  cfg.Retry.MaxRetries,
		BaseBackoff: time.Millisecond * 100,               // Default base backoff
		MaxBackoff:  time.Second * 10,                     // Default max backoff
		ShouldRetry: func(err error) bool { return true }, // Default retry policy
		Logger:      logger,
	}
}

// ExecuteWithRetry retries a function with exponential backoff, jitter, and error categorization.
func (rp *RetryPolicy) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	var attempt int
	baseBackoff := rp.BaseBackoff

	for attempt < rp.MaxRetries {
		err := operation()
		if err == nil {
			// Operation succeeded
			CountRetryAttempt("operation_name", "success") // Update success metric
			return nil
		}

		// Check if the error is retryable based on custom logic
		if !rp.ShouldRetry(err) {
			rp.Logger.Warn("Non-retryable error occurred", zap.Error(err))
			CountRetryAttempt("operation_name", "failure") // Update failure metric
			return err
		}

		// Calculate jitter using crypto/rand for added randomness
		jitterValue, _ := rand.Int(rand.Reader, big.NewInt(int64(baseBackoff)))
		jitter := time.Duration(jitterValue.Int64())
		waitTime := min(baseBackoff+jitter, rp.MaxBackoff)

		rp.Logger.Warn("Operation failed, retrying...", zap.Error(err), zap.Int("attempt", attempt+1), zap.Duration("wait_time", waitTime))

		// Dynamic retry delay using context
		retryCtx, cancel := context.WithTimeout(ctx, waitTime)
		defer cancel()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-retryCtx.Done():
			// Continue to next retry
		}

		// Ensure the backoff does not exceed maxBackoff
		baseBackoff = min(baseBackoff*2, rp.MaxBackoff)
		attempt++
	}

	CountRetryAttempt("operation_name", "exceeded") // Update metric for retries exhausted
	return fmt.Errorf("operation failed after %d retries, last error: %w", attempt, ctx.Err())
}

// min returns the smaller of two time durations.
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
