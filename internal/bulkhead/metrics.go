package bulkhead

import (
	observability "github.com/goletan/observability/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents the metrics-related operations, including registration with a monitoring system.
type Metrics struct{}

var (
	// LimitReached is a CounterVec metric that tracks the number of times bulkhead limits are reached, labeled by service.
	LimitReached = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "bulkhead_limit_reached_total",
			Help:      "Counts the number of times bulkhead limits have been reached.",
		},
		[]string{"service"},
	)
)

// InitMetrics registers and initializes metrics for observability components.
func InitMetrics(observer *observability.Observability) {
	observer.Metrics.Register(&Metrics{})
}

// Register attempts to register the LimitReached metric with the Prometheus registry.
// Returns an error if registration fails.
func (bm *Metrics) Register() error {
	if err := prometheus.Register(LimitReached); err != nil {
		return err
	}
	return nil
}

// CountLimitReached logs when a bulkhead limit is reached for a specific service.
func CountLimitReached(service string) {
	LimitReached.WithLabelValues(service).Inc()
}
