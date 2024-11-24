// /resilience/retry/metrics.go
package retry

import (
	"time"

	observability "github.com/goletan/observability/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type RetryMetrics struct{}

// Retry Metrics: Track retry attempts and latency.
var (
	RetryAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "retry_attempts_total",
			Help:      "Counts the number of retry attempts for operations.",
		},
		[]string{"operation", "status"},
	)

	RetryLatency = prometheus.NewHistogramVec(
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
	observer.Metrics.Register(&RetryMetrics{})
}

func (rm *RetryMetrics) Register() error {
	if err := prometheus.Register(RetryAttempts); err != nil {
		return err
	}

	if err := prometheus.Register(RetryLatency); err != nil {
		return err
	}

	return nil
}

// CountRetryAttempt logs retry attempts for operations.
func CountRetryAttempt(operation, status string) {
	RetryAttempts.WithLabelValues(operation, status).Inc()
}

// TrackRetryLatency records the latency of a retry attempt.
func TrackRetryLatency(operation string, latency time.Duration) {
	RetryLatency.WithLabelValues(operation).Observe(latency.Seconds())
}
