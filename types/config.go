// /resilience/types/config.go
package types

import "time"

type ResilienceConfig struct {
	Bulkhead struct {
		Capacity int           `mapstructure:"capacity"`
		Timeout  time.Duration `mapstructure:"timeout"`
	}

	CircuitBreaker struct {
		FailureRateThreshold float32       `mapstructure:"failure_rate_threshold"`
		FailureThreshold     int           `mapstructure:"failure_threshold"`
		SuccessThreshold     int           `mapstructure:"success_threshold"`
		ConsecutiveFailures  int           `mapstructure:"consecutive_failures"`
		MaxRequest           uint32        `mapstructure:"max_requests"`
		Interval             time.Duration `mapstructure:"interval"`
		Timeout              time.Duration `mapstructure:"timeout"`
	} `mapstructure:"circuit_breaker"`

	RateLimiter struct {
		RPS   int `mapstructure:"rps"`
		Burst int `mapstructure:"burst"`
	} `mapstructure:"rate_limiter"`

	Retry struct {
		MaxRetries int `mapstructure:"max_retries"`
	} `mapstructure:"retry"`
}
