// /resilience/types/resilience.go
package types

import (
	"context"

	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
)

// ResilienceService defines methods for executing operations with resilience mechanisms.
type ResilienceService interface {
	ExecuteWithResilience(ctx context.Context, operation func() error) error
}

type ResilienceMetrics struct{}

type DefaultResilienceService struct {
	MaxRetries     int
	ShouldRetry    func(error) bool
	Logger         *zap.Logger
	CircuitBreaker *gobreaker.CircuitBreaker[any]
}
