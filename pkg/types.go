// /resilience/pkg/types.go
package resilience

import "github.com/sony/gobreaker/v2"

// CircuitBreakerCallbacks defines optional callbacks for the circuit breaker events.
type CircuitBreakerCallbacks struct {
	OnOpen        func()
	OnClose       func()
	OnStateChange func(name string, from, to gobreaker.State)
	OnSuccess     func()
	OnFailure     func()
}
