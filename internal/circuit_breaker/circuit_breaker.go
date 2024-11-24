// /resilience/circuit_breaker/circuit_breaker.go
package circuit_breaker

import (
	"context"
	"time"

	observability "github.com/goletan/observability/pkg"
	"github.com/goletan/resilience/internal/types"
	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
)

type CircuitBreaker struct {
	cb          *gobreaker.CircuitBreaker[types.CircuitBreakerInterface]
	logger      *zap.Logger
	metrics     *CircuitBreakerMetrics
	lastState   gobreaker.State
	stateChange time.Time
}

// NewCircuitBreaker constructs a new CircuitBreaker instance.
func NewCircuitBreaker(cfg *types.ResilienceConfig, callbacks *types.CircuitBreakerCallbacks, obs *observability.Observability) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "GoletanCircuitBreaker",
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.ConsecutiveFailures > uint32(cfg.CircuitBreaker.ConsecutiveFailures) || failureRate > cfg.CircuitBreaker.FailureRateThreshold
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			cb.logger.Info("Circuit breaker state changed", zap.String("name", name), zap.String("from", from.String()), zap.String("to", to.String()))
			RecordCircuitBreakerStateChange(name, from.String(), to.String())

			// Track duration in the previous state using the `stateChange` field from `CircuitBreaker` struct
			if from != to {
				duration := time.Since(cb.stateChange)
				RecordStateDuration(name, from.String(), duration)
			}

			// Update the `stateChange` time for the new state
			cb.stateChange = time.Now()

			if callbacks != nil && callbacks.OnStateChange != nil {
				callbacks.OnStateChange(name, from, to)
			}
		},
	}

	cb := gobreaker.NewCircuitBreaker[types.CircuitBreakerInterface](settings)

	// Initialize the CircuitBreaker with the current time
	return &CircuitBreaker{
		cb:          cb,
		logger:      obs.Logger,
		metrics:     &CircuitBreakerMetrics{},
		lastState:   gobreaker.StateClosed, // Assuming we start with a closed state
		stateChange: time.Now(),
	}
}

// Execute runs the provided operation and handles fallback if the circuit breaker is open.
func (c *CircuitBreaker) Execute(ctx context.Context, operation func() error, fallback func() error) error {
	resultCh := make(chan error, 1)

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

// Shutdown gracefully shuts down the circuit breaker and releases any resources.
func (c *CircuitBreaker) Shutdown(ctx context.Context) error {
	c.logger.Info("Shutting down circuit breaker")
	return nil
}
