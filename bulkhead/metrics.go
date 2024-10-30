// /resilience/bulkhead/metrics.go
package bulkhead

import (
	"github.com/goletan/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type BulkheadMetrics struct{}

var (
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
	metrics.NewManager().Register(&BulkheadMetrics{})
}

func (bm *BulkheadMetrics) Register() error {
	if err := prometheus.Register(BulkheadLimitReached); err != nil {
		return err
	}
	return nil
}

// CountBulkheadLimitReached logs when a bulkhead limit is reached for a specific service.
func CountBulkheadLimitReached(service string) {
	BulkheadLimitReached.WithLabelValues(service).Inc()
}
