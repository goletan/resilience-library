// /resilience/rate_limiter/metrics.go
package rate_limiter

import (
	"time"

	observability "github.com/goletan/observability/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type RateLimiterMetrics struct{}

// RateLimiter Metrics: Track rate limit attempts and latency.
var (
	RateLimitReached = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "rate_limit_reached_total",
			Help:      "Counts the number of times rate limit has been reached.",
		},
		[]string{"operation"},
	)

	RateLimitLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "rate_limit_latency_seconds",
			Help:      "Latency for rate-limited operations in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func InitMetrics(observer *observability.Observability) {
	observer.Metrics.Register(&RateLimiterMetrics{})
}

func (rlm *RateLimiterMetrics) Register() error {
	if err := prometheus.Register(RateLimitReached); err != nil {
		return err
	}

	if err := prometheus.Register(RateLimitLatency); err != nil {
		return err
	}

	return nil
}

// CountRateLimit logs when a rate limit is reached for an operation.
func CountRateLimit(operation string) {
	RateLimitReached.WithLabelValues(operation).Inc()
}

// TrackRateLimitLatency records the latency for rate-limited operations.
func TrackRateLimitLatency(operation string, latency time.Duration) {
	RateLimitLatency.WithLabelValues(operation).Observe(latency.Seconds())
}
