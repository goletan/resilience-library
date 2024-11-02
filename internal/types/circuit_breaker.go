// /resilience/types/circuit_breaker.go
package types

import (
	"context"

	"github.com/sony/gobreaker/v2"
)

// CircuitBreakerInterface defines the methods needed from the CircuitBreaker.
type CircuitBreakerInterface interface {
	Execute(ctx context.Context, operation func() error, fallback func() error) error
	Shutdown(ctx context.Context) error
}

// CircuitBreakerCallbacks defines optional callbacks for the circuit breaker events.
type CircuitBreakerCallbacks struct {
	// OnStateChange is called when the circuit breaker changes state (e.g., from closed to open).
	OnStateChange func(name string, from, to gobreaker.State)

	// OnSuccess is called when an operation completes successfully.
	OnSuccess func()

	// OnFailure is called when an operation fails.
	OnFailure func()
}
