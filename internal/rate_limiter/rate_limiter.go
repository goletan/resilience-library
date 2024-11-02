// /resilience/rate_limiter/rate_limiter.go
package rate_limiter

import (
	"context"
	"fmt"
	"sync"

	observability "github.com/goletan/observability/pkg"
	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiter defines a rate limiter with specific settings.
type RateLimiter struct {
	limiter *rate.Limiter
	logger  *zap.Logger
}

// rateLimiterRegistry stores multiple rate limiters for different serviceNames.
var (
	rateLimiterRegistry = make(map[string]*RateLimiter)
	mu                  sync.Mutex
)

// NewRateLimiter initializes a new rate limiter for a given serviceName and configuration.
func NewRateLimiter(cfg *types.ResilienceConfig, serviceName string, observer *observability.Observability) {
	mu.Lock()
	defer mu.Unlock()

	// Initialize a new rate limiter for the given serviceName
	rateLimiterRegistry[serviceName] = &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(cfg.RateLimiter.RPS), cfg.RateLimiter.Burst),
		logger:  observer.Logger,
	}

	observer.Logger.Info("Rate limiter initialized", zap.String("serviceName", serviceName), zap.Int("rps", cfg.RateLimiter.RPS), zap.Int("burst", cfg.RateLimiter.Burst))
}

// GetRateLimiter retrieves a rate limiter for the given serviceName.
func GetRateLimiter(serviceName string) (*RateLimiter, bool) {
	mu.Lock()
	defer mu.Unlock()

	rl, exists := rateLimiterRegistry[serviceName]
	return rl, exists
}

// ExecuteWithRateLimiting executes the provided serviceName if the rate limit allows it.
func ExecuteWithRateLimiting(ctx context.Context, serviceName string, fn func() error) error {
	// Get the rate limiter for the serviceName
	rl, exists := GetRateLimiter(serviceName)
	if !exists {
		return fmt.Errorf("%s: Rate limiter not initialized", serviceName)
	}

	// Wait for permission to proceed
	if err := rl.limiter.Wait(ctx); err != nil {
		rl.logger.Warn("Rate limit exceeded", zap.String("serviceName", serviceName), zap.Error(err))
		CountRateLimit(serviceName) // Update metric for rate limit reached
		return err
	}

	// Execute the serviceName
	return fn()
}
