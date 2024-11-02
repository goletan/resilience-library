package resilience

import (
	"context"

	observability "github.com/goletan/observability/pkg"
	"github.com/goletan/resilience/internal/bulkhead"
	"github.com/goletan/resilience/internal/circuit_breaker"
	"github.com/goletan/resilience/internal/config"
	"github.com/goletan/resilience/internal/metrics"
	"github.com/goletan/resilience/internal/rate_limiter"
	"github.com/goletan/resilience/internal/retry"
	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

type DefaultResilienceService struct {
	MaxRetries     int
	ShouldRetry    func(error) bool
	Logger         *zap.Logger
	CircuitBreaker types.CircuitBreakerInterface
	RetryPolicy    types.RetryPolicyInterface
}

// NewResilienceService initializes a new DefaultResilienceService with bulkhead, circuit breaker, and rate limiter.
func NewResilienceService(serviceName string, logger *zap.Logger, shouldRetry func(error) bool, callbacks *types.CircuitBreakerCallbacks) *DefaultResilienceService {
	cfg, err := config.LoadResilienceConfig(logger)
	if err != nil {
		logger.Fatal("Failed to load resilience configuration", zap.Error(err))
	}

	observer, err := observability.NewObserver()
	if err != nil {
		logger.Error("Failed to initialize observability", zap.Error(err))
		return nil
	}

	bulkhead.Init(cfg, serviceName)
	metrics.InitMetrics(observer)

	cb := circuit_breaker.NewCircuitBreaker(cfg, callbacks, logger)
	rate_limiter.NewRateLimiter(cfg, serviceName, observer)

	retryPolicy := retry.NewRetryPolicy(cfg, logger)

	return &DefaultResilienceService{
		MaxRetries:     cfg.Retry.MaxRetries,
		ShouldRetry:    shouldRetry,
		Logger:         logger,
		CircuitBreaker: cb,
		RetryPolicy:    retryPolicy,
	}
}

// ExecuteWithRetry retries an operation with the configured retry policy.
func (r *DefaultResilienceService) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	return r.RetryPolicy.ExecuteWithRetry(ctx, operation)
}

// Shutdown gracefully shuts down the resilience service, releasing any resources held by components.
func (r *DefaultResilienceService) Shutdown(ctx context.Context) error {
	r.Logger.Info("Shutting down resilience service")

	// Add shutdown logic here for each resilience component if necessary.
	if err := r.CircuitBreaker.Shutdown(ctx); err != nil {
		r.Logger.Error("Failed to shutdown circuit breaker", zap.Error(err))
		return err
	}

	r.Logger.Info("Resilience service shut down successfully")
	return nil
}
