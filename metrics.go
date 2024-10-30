package resilience

import (
	"github.com/goletan/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type ResilienceMetrics struct{}

// Resilience Metrics: Track resilience patterns like retries and circuit breakers.
var (
	CircuitBreakerStateChange = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "circuit_breaker_state_changes_total",
			Help:      "Tracks state changes in circuit breakers.",
		},
		[]string{"circuit", "from", "to"},
	)
	RetryAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "retry_attempts_total",
			Help:      "Counts the number of retry attempts for operations.",
		},
		[]string{"operation", "status"},
	)
	Timeouts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "timeouts_total",
			Help:      "Counts the number of timeouts in various operations.",
		},
		[]string{"operation", "service"},
	)
	BulkheadLimitReached = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "bulkhead_limit_reached_total",
			Help:      "Counts the number of times bulkhead limits have been reached.",
		},
		[]string{"service"},
	)
)

func InitMetrics() {
	metrics.NewManager().Register(&ResilienceMetrics{})
}

func (em *ResilienceMetrics) Register() error {
	if err := prometheus.Register(CircuitBreakerStateChange); err != nil {
		return err
	}

	if err := prometheus.Register(RetryAttempts); err != nil {
		return err
	}

	if err := prometheus.Register(Timeouts); err != nil {
		return err
	}

	if err := prometheus.Register(BulkheadLimitReached); err != nil {
		return err
	}

	return nil
}

// RecordCircuitBreakerStateChange logs state changes in the circuit breaker.
func RecordCircuitBreakerStateChange(circuit, from, to string) {
	CircuitBreakerStateChange.WithLabelValues(circuit, from, to).Inc()
}

// CountRetryAttempt logs retry attempts for operations.
func CountRetryAttempt(operation, status string) {
	RetryAttempts.WithLabelValues(operation, status).Inc()
}

// CountTimeout logs a timeout event for a specific operation and service.
func CountTimeout(operation, service string) {
	Timeouts.WithLabelValues(operation, service).Inc()
}

// CountBulkheadLimitReached logs when a bulkhead limit is reached for a specific service.
func CountBulkheadLimitReached(service string) {
	BulkheadLimitReached.WithLabelValues(service).Inc()
}
