package resilience

import (
	"context"

	"github.com/goletan/observability/pkg"
	"github.com/goletan/resilience/internal/bulkhead"
	"github.com/goletan/resilience/internal/circuit_breaker"
	"github.com/goletan/resilience/internal/config"
	"github.com/goletan/resilience/internal/metrics"
	"github.com/goletan/resilience/internal/rate_limiter"
	"github.com/goletan/resilience/internal/retry"
	"github.com/goletan/resilience/shared/types"
	"go.uber.org/zap"
)

type DefaultResilienceService struct {
	MaxRetries     int
	ShouldRetry    func(error) bool
	Observability  *observability.Observability
	CircuitBreaker *circuit_breaker.CircuitBreaker
	RetryPolicy    *retry.Policy
}

func NewResilienceService(serviceName string, obs *observability.Observability, shouldRetry func(error) bool, callbacks *types.CircuitBreakerCallbacks) *DefaultResilienceService {
	cfg, err := config.LoadResilienceConfig(obs.Logger)
	if err != nil {
		obs.Logger.Fatal("Failed to load resilience-library configuration", zap.Error(err))
	}

	bulkhead.Init(cfg, serviceName)
	metrics.InitMetrics(obs)

	// Convert public callbacks to internal ones
	internalCallbacks := &types.CircuitBreakerCallbacks{
		OnOpen:        callbacks.OnOpen,
		OnClose:       callbacks.OnClose,
		OnStateChange: callbacks.OnStateChange,
		OnSuccess:     callbacks.OnSuccess,
		OnFailure:     callbacks.OnFailure,
	}

	cb := circuit_breaker.NewCircuitBreaker(cfg, internalCallbacks, obs)
	rate_limiter.NewRateLimiter(cfg, serviceName, obs)

	retryPolicy := retry.NewRetryPolicy(cfg, obs)

	return &DefaultResilienceService{
		MaxRetries:     cfg.Retry.MaxRetries,
		ShouldRetry:    shouldRetry,
		Observability:  obs,
		CircuitBreaker: cb,
		RetryPolicy:    retryPolicy,
	}
}

// ExecuteWithRetry retries an operation with the configured retry policy.
func (r *DefaultResilienceService) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	return r.RetryPolicy.ExecuteWithRetry(ctx, operation)
}

// Shutdown gracefully shuts down the resilience-library service, releasing any resources held by components.
func (r *DefaultResilienceService) Shutdown(ctx *context.Context) error {
	r.Observability.Logger.Info("Shutting down resilience-library service")

	// Add shutdown logic here for each resilience-library component if necessary.
	if err := r.CircuitBreaker.Shutdown(ctx); err != nil {
		r.Observability.Logger.Error("Failed to shutdown circuit breaker", zap.Error(err))
		return err
	}

	r.Observability.Logger.Info("Resilience service shut down successfully")
	return nil
}
