package resilience

import (
	"github.com/goletan/resilience/bulkhead"
	"github.com/goletan/resilience/circuit_breaker"
	"github.com/goletan/resilience/rate_limiter"
	"github.com/goletan/resilience/types"
	"go.uber.org/zap"
)

// NewResilienceService initializes a new DefaultResilienceService with bulkhead, circuit breaker, and rate limiter.
func NewResilienceService(serviceName string, logger *zap.Logger, shouldRetry func(error) bool, callbacks *types.CircuitBreakerCallbacks) *types.DefaultResilienceService {
	cfg, err := LoadResilienceConfig(logger)
	if err != nil {
		logger.Fatal("Failed to load resilience configuration", zap.Error(err))
	}

	// Initialize Bulkhead
	bulkhead.Init(cfg, serviceName)
	bulkhead.InitMetrics()

	// Initialize Circuit Breaker
	cb := circuit_breaker.NewCircuitBreaker(cfg, callbacks)
	circuit_breaker.InitMetrics()

	// Initialize Rate Limiter
	rate_limiter.NewRateLimiter(serviceName, cfg, logger)
	rate_limiter.InitMetrics()

	return &types.DefaultResilienceService{
		MaxRetries:     cfg.Retry.MaxRetries,
		ShouldRetry:    shouldRetry,
		Logger:         logger,
		CircuitBreaker: cb,
	}
}
