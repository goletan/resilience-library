// /resilience/types/circuit_breaker.go
package types

import "github.com/sony/gobreaker/v2"

// CircuitBreakerCallbacks defines optional callbacks for the circuit breaker events.
type CircuitBreakerCallbacks struct {
	// OnStateChange is called when the circuit breaker changes state (e.g., from closed to open).
	OnStateChange func(name string, from, to gobreaker.State)

	// OnSuccess is called when an operation completes successfully.
	OnSuccess func()

	// OnFailure is called when an operation fails.
	OnFailure func()
}
