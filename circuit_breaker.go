// /resilience/circuitbreaker/circuit_breaker.go
package resilience

import (
	"context"
	"sync"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

var (
	circuitBreakerInstance *gobreaker.CircuitBreaker
	once                   sync.Once
	logger                 *zap.Logger
)

// InitCircuitBreaker initializes the circuit breaker with settings and registers state changes.
func InitCircuitBreaker(log *zap.Logger) {
	logger = log

	once.Do(func() {
		settings := createCircuitBreakerSettings()
		circuitBreakerInstance = gobreaker.NewCircuitBreaker(settings)
	})
}

// createCircuitBreakerSettings returns the settings used to configure the circuit breaker.
func createCircuitBreakerSettings() gobreaker.Settings {
	return gobreaker.Settings{
		Name:        "GoletanCircuitBreaker",
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.ConsecutiveFailures > 3 || failureRate > 0.5
		},
		OnStateChange: handleStateChange,
	}
}

// handleStateChange logs circuit breaker state changes.
func handleStateChange(name string, from gobreaker.State, to gobreaker.State) {
	logger.Info("Circuit breaker state changed",
		zap.String("name", name),
		zap.String("from", from.String()),
		zap.String("to", to.String()))
	RecordCircuitBreakerStateChange(name, from.String(), to.String())
}

// GetCircuitBreakerInstance ensures the circuit breaker is initialized and returns the instance.
func GetCircuitBreakerInstance(log *zap.Logger) *gobreaker.CircuitBreaker {
	InitCircuitBreaker(log)

	return circuitBreakerInstance
}

// ExecuteWithCircuitBreaker executes a function with circuit breaker protection.
func ExecuteWithCircuitBreaker(ctx context.Context, log *zap.Logger, action func() (interface{}, error), errorHandler func(error) (interface{}, error)) (interface{}, error) {
	cb := GetCircuitBreakerInstance(log)
	return cb.Execute(func() (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			result, err := action()
			if err != nil {
				return errorHandler(err)
			}
			return result, nil
		}
	})
}
