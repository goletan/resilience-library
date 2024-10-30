// /resilience/rate_limiter/metrics.go
package rate_limiter

import (
	"github.com/goletan/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type RateLimiterMetrics struct{}

// RateLimiter Metrics: Track rate limit attempts.
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
)

func InitMetrics() {
	metrics.NewManager().Register(&RateLimiterMetrics{})
}

func (rlm *RateLimiterMetrics) Register() error {
	if err := prometheus.Register(RateLimitReached); err != nil {
		return err
	}
	return nil
}

// CountRateLimit logs when a rate limit is reached for an operation.
func CountRateLimit(operation string) {
	RateLimitReached.WithLabelValues(operation).Inc()
}
