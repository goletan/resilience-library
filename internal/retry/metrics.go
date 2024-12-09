package retry

import (
	"time"

	observability "github.com/goletan/observability/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct{}

var (
	// Attempts counts the number of retry attempts for operations, labeled by operation type and status.
	Attempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "retry_attempts_total",
			Help:      "Counts the number of retry attempts for operations.",
		},
		[]string{"operation", "status"},
	)

	// Latency records the latency of retry attempts in seconds, labeled by operation, using default Prometheus buckets.
	Latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "retry_latency_seconds",
			Help:      "Latency of retry attempts in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func InitMetrics(observer *observability.Observability) {
	observer.Metrics.Register(&Metrics{})
}

func (rm *Metrics) Register() error {
	if err := prometheus.Register(Attempts); err != nil {
		return err
	}

	if err := prometheus.Register(Latency); err != nil {
		return err
	}

	return nil
}

// CountAttempt logs retry attempts for operations.
func CountAttempt(operation, status string) {
	Attempts.WithLabelValues(operation, status).Inc()
}

// TrackLatency records the latency of a retry attempt.
func TrackLatency(operation string, latency time.Duration) {
	Latency.WithLabelValues(operation).Observe(latency.Seconds())
}
