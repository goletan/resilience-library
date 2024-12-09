package circuit_breaker

import (
	"context"
	"errors"
	"github.com/goletan/observability/shared/logger"
	"time"

	"github.com/goletan/resilience/internal/types"
	sharedTypes "github.com/goletan/resilience/shared/types"
	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
)

type CircuitBreaker struct {
	cb          *gobreaker.CircuitBreaker[sharedTypes.CircuitBreakerInterface]
	logger      *logger.ZapLogger
	metrics     *CircuitBreakerMetrics
	lastState   gobreaker.State
	stateChange time.Time
	callbacks   *sharedTypes.CircuitBreakerCallbacks
}

// NewCircuitBreaker constructs a new CircuitBreaker instance.
func NewCircuitBreaker(cfg *types.ResilienceConfig, callbacks *sharedTypes.CircuitBreakerCallbacks, log *logger.ZapLogger) *CircuitBreaker {
	cb := &CircuitBreaker{
		logger:      log,
		metrics:     &CircuitBreakerMetrics{},
		lastState:   gobreaker.StateClosed, // Assuming we start with a closed state
		stateChange: time.Now(),
		callbacks:   callbacks,
	}

	settings := gobreaker.Settings{
		Name:        "GoletanCircuitBreaker",
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.ConsecutiveFailures > uint32(cfg.CircuitBreaker.ConsecutiveFailures) || failureRate > cfg.CircuitBreaker.FailureRateThreshold
		},
		OnStateChange: cb.handleStateChange,
	}

	cb.cb = gobreaker.NewCircuitBreaker[sharedTypes.CircuitBreakerInterface](settings)

	return cb
}

// Execute runs the provided operation and handles fallback if the circuit breaker is open.
func (c *CircuitBreaker) Execute(ctx context.Context, operation func() error, fallback func() error) error {
	resultCh := make(chan error, 1)

	go func() {
		resultCh <- func() error {
			_, err := c.cb.Execute(func() (sharedTypes.CircuitBreakerInterface, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					return nil, operation()
				}
			})

			// If the circuit breaker is open and a fallback is provided, execute the fallback.
			if errors.Is(err, gobreaker.ErrOpenState) && fallback != nil {
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
func (c *CircuitBreaker) Shutdown(ctx *context.Context) error {
	c.logger.Info("Shutting down circuit breaker")
	return nil
}

// handleStateChange is a private function to handle state changes for the circuit breaker.
func (c *CircuitBreaker) handleStateChange(name string, from, to gobreaker.State) {
	c.logger.Info("Circuit breaker state changed", zap.String("name", name), zap.String("from", from.String()), zap.String("to", to.String()))
	RecordCircuitBreakerStateChange(name, from.String(), to.String())

	// Track duration in the previous state using the stateChange field from the struct
	if from != to {
		duration := time.Since(c.stateChange)
		RecordStateDuration(name, from.String(), duration)
	}

	// Update the stateChange time for the new state
	c.stateChange = time.Now()

	if c.callbacks != nil && c.callbacks.OnStateChange != nil {
		c.callbacks.OnStateChange(name, from, to)
	}
}
