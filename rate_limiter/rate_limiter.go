// /resilience/rate_limiter/rate_limiter.go
package rate_limiter

import (
	"context"
	"fmt"
	"sync"

	"github.com/goletan/resilience/types"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiter defines a rate limiter with specific settings.
type RateLimiter struct {
	limiter *rate.Limiter
	logger  *zap.Logger
}

// rateLimiterRegistry stores multiple rate limiters for different operations.
var (
	rateLimiterRegistry = make(map[string]*RateLimiter)
	mu                  sync.Mutex
)

// NewRateLimiter initializes a new rate limiter for a given operation and configuration.
func NewRateLimiter(operation string, cfg *types.ResilienceConfig, logger *zap.Logger) {
	mu.Lock()
	defer mu.Unlock()

	// Initialize a new rate limiter for the given operation
	rateLimiterRegistry[operation] = &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(cfg.RateLimiter.RPS), cfg.RateLimiter.Burst),
		logger:  logger,
	}
	logger.Info("Rate limiter initialized", zap.String("operation", operation), zap.Int("rps", cfg.RateLimiter.RPS), zap.Int("burst", cfg.RateLimiter.Burst))
}

// GetRateLimiter retrieves a rate limiter for the given operation.
func GetRateLimiter(operation string) (*RateLimiter, bool) {
	mu.Lock()
	defer mu.Unlock()

	rl, exists := rateLimiterRegistry[operation]
	return rl, exists
}

// ExecuteWithRateLimiting executes the provided operation if the rate limit allows it.
func ExecuteWithRateLimiting(ctx context.Context, operation string, fn func() error) error {
	// Get the rate limiter for the operation
	rl, exists := GetRateLimiter(operation)
	if !exists {
		return fmt.Errorf("%s: Rate limiter not initialized", operation)
	}

	// Wait for permission to proceed
	if err := rl.limiter.Wait(ctx); err != nil {
		rl.logger.Warn("Rate limit exceeded", zap.String("operation", operation), zap.Error(err))
		CountRateLimit(operation) // Update metric for rate limit reached
		return err
	}

	// Execute the operation
	return fn()
}
