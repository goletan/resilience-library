// /resilience/circuit_breaker/circuit_breaker.go
package circuit_breaker

import (
	"context"
	"sync"

	"github.com/goletan/resilience/internal/types"
	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

type CircuitBreaker struct {
	cb     *gobreaker.CircuitBreaker[types.CircuitBreakerInterface]
	logger *zap.Logger
}

// Modify NewCircuitBreaker to take a logger argument
func NewCircuitBreaker(cfg *types.ResilienceConfig, callbacks *types.CircuitBreakerCallbacks, logger *zap.Logger) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "GoletanCircuitBreaker",
		MaxRequests: uint32(cfg.CircuitBreaker.MaxRequest),
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.ConsecutiveFailures > uint32(cfg.CircuitBreaker.ConsecutiveFailures) || failureRate > float64(cfg.CircuitBreaker.FailureRateThreshold)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Info("Circuit breaker state changed", zap.String("name", name), zap.String("from", from.String()), zap.String("to", to.String()))
			RecordCircuitBreakerStateChange(name, from.String(), to.String())

			if callbacks != nil && callbacks.OnStateChange != nil {
				callbacks.OnStateChange(name, from, to)
			}
		},
	}

	cb := gobreaker.NewCircuitBreaker[types.CircuitBreakerInterface](settings)
	return &CircuitBreaker{cb: cb, logger: logger}
}

func (c *CircuitBreaker) Shutdown(ctx context.Context) error {
	// Here, you could reset states, stop monitoring, or clean up resources.
	c.logger.Info("Shutting down circuit breaker")
	return nil
}

// Execute runs the provided operation and handles fallback if the circuit breaker is open.
func (c *CircuitBreaker) Execute(ctx context.Context, operation func() error, fallback func() error) error {
	resultCh := make(chan error, 1)

	// Run the operation in a separate goroutine.
	go func() {
		resultCh <- func() error {
			_, err := c.cb.Execute(func() (types.CircuitBreakerInterface, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					return nil, operation()
				}
			})

			// If the circuit breaker is open and a fallback is provided, execute the fallback.
			if err == gobreaker.ErrOpenState && fallback != nil {
				c.logger.Warn("Circuit breaker open, executing fallback")
				return fallback()
			}

			return err
		}()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}
