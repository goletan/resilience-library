package resilience

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
	OnOpen        func(name string, from, to gobreaker.State)
	OnClose       func(name string, from, to gobreaker.State)
	OnStateChange func(name string, from, to gobreaker.State)
	OnSuccess     func(name string)
	OnFailure     func(name string, err error)
}
