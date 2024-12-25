package metrics

import (
	observability "github.com/goletan/observability-library/pkg"
	"github.com/goletan/resilience/internal/bulkhead"
	"github.com/goletan/resilience/internal/circuit_breaker"
	"github.com/goletan/resilience/internal/rate_limiter"
	"github.com/goletan/resilience/internal/retry"
)

// InitMetrics initializes all metrics for resilience-library components
func InitMetrics(observer *observability.Observability) {
	bulkhead.InitMetrics(observer)
	circuit_breaker.InitMetrics(observer)
	rate_limiter.InitMetrics(observer)
	retry.InitMetrics(observer)
}
