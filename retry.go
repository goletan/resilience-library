// /resilience/retry/retry.go
package resilience

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"
)

// InitRetry initializes the logger for retry operations
func InitRetry(log *zap.Logger) {
	logger = log
}

// ExecuteWithRetry retries a function with exponential backoff, jitter, and error categorization.
func ExecuteWithRetry(ctx context.Context, operation func() error, maxRetries int, shouldRetry func(error) bool) error {
	var attempt int
	baseBackoff := time.Millisecond * 100
	const maxBackoff = time.Second * 10

	for attempt < maxRetries {
		err := operation()
		if err == nil {
			return nil
		}

		// Check if the error is retryable based on custom logic
		if !shouldRetry(err) {
			logger.Warn("Non-retryable error occurred", zap.Error(err))
			return err
		}

		// Calculate jitter using crypto/rand for added randomness
		jitterValue, _ := rand.Int(rand.Reader, big.NewInt(int64(baseBackoff)))
		jitter := time.Duration(jitterValue.Int64())
		waitTime := min(baseBackoff+jitter, maxBackoff)

		logger.Warn("Operation failed, retrying...", zap.Error(err), zap.Int("attempt", attempt+1), zap.Duration("wait_time", waitTime))

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
		baseBackoff = min(baseBackoff*2, maxBackoff)
		attempt++
	}

	return fmt.Errorf("operation failed after %d retries, last error: %w", attempt, ctx.Err())
}

// min returns the smaller of two time durations
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
