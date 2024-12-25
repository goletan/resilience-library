package metrics

import (
	observability "github.com/goletan/observability-library/pkg"
	"github.com/goletan/resilience-library/internal/bulkhead"
	"github.com/goletan/resilience-library/internal/circuit_breaker"
	"github.com/goletan/resilience-library/internal/rate_limiter"
	"github.com/goletan/resilience-library/internal/retry"
)

// InitMetrics initializes all metrics for resilience-library components
func InitMetrics(observer *observability.Observability) {
	bulkhead.InitMetrics(observer)
	circuit_breaker.InitMetrics(observer)
	rate_limiter.InitMetrics(observer)
	retry.InitMetrics(observer)
}
