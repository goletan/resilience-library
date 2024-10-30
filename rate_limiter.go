package resilience

import (
	"context"

	"github.com/goletan/config"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var (
	rateLimiterInstance *rate.Limiter
)

// Init initializes the rate limiter module by utilizing the configuration library.
// This function loads the relevant configuration parameters and sets up the rate limiter accordingly.
func InitRateLimiter(logger *zap.Logger) {
	var cfg ResilienceConfig

	// Load the configuration
	err := config.LoadConfig("Resilience", &cfg, logger)
	if err != nil {
		logger.Fatal("Failed to load rate limiter configuration", zap.Error(err))
	}

	once.Do(func() {
		rateLimiterInstance = rate.NewLimiter(rate.Limit(cfg.RateLimiter.RPS), cfg.RateLimiter.Burst)
		logger.Info("Rate limiter initialized", zap.Int("rps", cfg.RateLimiter.RPS), zap.Int("burst", cfg.RateLimiter.Burst))
	})
}

// GetRateLimiterInstance guarantees that the rate limiter is properly initialized and returns its instance.
func GetRateLimiterInstance() *rate.Limiter {
	return rateLimiterInstance
}

// ExecuteWithRateLimiting enforces rate limiting when executing the provided operation.
// If the rate limit is exceeded, it logs a warning and returns an error; otherwise, it proceeds to execute the operation.
func ExecuteWithRateLimiting(ctx context.Context, operation func() error) error {
	limiter := GetRateLimiterInstance()
	if err := limiter.Wait(ctx); err != nil {
		logger.Warn("Rate limit exceeded", zap.Error(err))
		return err
	}
	return operation()
}
