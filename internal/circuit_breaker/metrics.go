package circuit_breaker

import (
	"time"

	observability "github.com/goletan/observability-library/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type CircuitBreakerMetrics struct{}

var (
	CircuitBreakerStateChange = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience_library",
			Name:      "circuit_breaker_state_changes_total",
			Help:      "Tracks state changes in circuit breakers.",
		},
		[]string{"circuit", "from", "to"},
	)

	CircuitBreakerRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience_library",
			Name:      "circuit_breaker_requests_total",
			Help:      "Tracks the number of requests through the circuit breaker.",
		},
		[]string{"circuit", "state"},
	)

	CircuitBreakerStateDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "resilience_library",
			Name:      "circuit_breaker_state_duration_seconds",
			Help:      "Tracks the duration of time spent in each circuit breaker state.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"circuit", "state"},
	)
)

func InitMetrics(observer *observability.Observability) {
	observer.Metrics.Register(&CircuitBreakerMetrics{})
}

func (cbm *CircuitBreakerMetrics) Register() error {
	if err := prometheus.Register(CircuitBreakerStateChange); err != nil {
		return err
	}
	if err := prometheus.Register(CircuitBreakerRequestCount); err != nil {
		return err
	}
	if err := prometheus.Register(CircuitBreakerStateDuration); err != nil {
		return err
	}
	return nil
}

// RecordCircuitBreakerStateChange logs state changes in the circuit breaker.
func RecordCircuitBreakerStateChange(circuit, from, to string) {
	CircuitBreakerStateChange.WithLabelValues(circuit, from, to).Inc()
	CircuitBreakerRequestCount.WithLabelValues(circuit, to).Inc()
}

// RecordStateDuration logs the duration spent in each state.
func RecordStateDuration(circuit, state string, duration time.Duration) {
	CircuitBreakerStateDuration.WithLabelValues(circuit, state).Observe(duration.Seconds())
}
