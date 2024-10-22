// resilience/resilience_service.go
package resilience

import (
	"context"

	"github.com/goletan/config"
	cb "github.com/goletan/resilience/circuitbreaker"
	retry "github.com/goletan/resilience/retry"
	"go.uber.org/zap"
)

var cfg ResilienceConfig

func LoadResilienceConfig(logger *zap.Logger) (*ResilienceConfig, error) {
	if err := config.LoadConfig("Resilience", &cfg, logger); err != nil {
		logger.Error("Failed to load resilience configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}

func NewDefaultResilienceService(logger *zap.Logger, shouldRetry func(error) bool) *DefaultResilienceService {
	cfg, err := LoadResilienceConfig(logger)
	if err != nil {
		logger.Fatal("Failed to load resilience configuration", zap.Error(err))
	}

	return &DefaultResilienceService{
		MaxRetries:  cfg.Retry.MaxRetries,
		ShouldRetry: shouldRetry,
		Logger:      logger,
	}
}

func (r *DefaultResilienceService) ExecuteWithResilience(ctx context.Context, operation func() error) error {
	errorHandler := func(err error) (interface{}, error) {
		r.Logger.Warn("Circuit breaker fallback triggered", zap.Error(err))
		return nil, err
	}

	_, err := cb.ExecuteWithCircuitBreaker(ctx, func() (interface{}, error) {
		return nil, retry.ExecuteWithRetry(ctx, operation, r.MaxRetries, r.ShouldRetry)
	}, errorHandler)

	if err != nil {
		r.Logger.Error("Failed to execute operation with resilience", zap.Error(err))
	}

	return err
}
