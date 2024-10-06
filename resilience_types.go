package resilience

import (
	"context"

	"go.uber.org/zap"
)

// ResilienceService defines methods for executing operations with resilience mechanisms.
type ResilienceService interface {
	ExecuteWithResilience(ctx context.Context, operation func() error) error
}

// ResilienceConfig holds all resilience-related configurations.
type ResilienceConfig struct {
	Retry struct {
		MaxRetries int `mapstructure:"max_retries"`
	} `mapstructure:"retry"`
}

type DefaultResilienceService struct {
	MaxRetries  int
	ShouldRetry func(error) bool
	Logger      *zap.Logger
}
