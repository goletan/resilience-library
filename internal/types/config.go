package types

import "time"

type ResilienceConfig struct {
	Bulkhead struct {
		Capacity int           `mapstructure:"capacity"`
		Timeout  time.Duration `mapstructure:"timeout"`
	} `mapstructure:"bulkhead"`

	CircuitBreaker struct {
		FailureRateThreshold float64       `mapstructure:"failure_rate_threshold"` // Updated to float64 for more precision
		FailureThreshold     int           `mapstructure:"failure_threshold"`
		SuccessThreshold     int           `mapstructure:"success_threshold"`
		ConsecutiveFailures  int           `mapstructure:"consecutive_failures"`
		MaxRequests          uint32        `mapstructure:"max_requests"`
		Interval             time.Duration `mapstructure:"interval"`
		Timeout              time.Duration `mapstructure:"timeout"`
		StateDuration        time.Duration `mapstructure:"state_duration"` // Duration to track in a specific state
	} `mapstructure:"circuit_breaker"`

	RateLimiter struct {
		RPS   int `mapstructure:"rps"`
		Burst int `mapstructure:"burst"`
	} `mapstructure:"rate_limiter"`

	Retry struct {
		MaxRetries     int           `mapstructure:"max_retries"`
		InitialBackoff time.Duration `mapstructure:"initial_backoff"` // Initial delay for retry backoff
		MaxBackoff     time.Duration `mapstructure:"max_backoff"`     // Maximum delay for retry backoff
		BackoffFactor  float64       `mapstructure:"backoff_factor"`  // Exponential factor for increasing delay
	} `mapstructure:"retry"`

	DLQ struct {
		Topic   string   `mapstructure:"topic"`
		Brokers []string `mapstructure:"brokers"`
		GroupID string   `mapstructure:"group_id"`
	} `mapstructure:"dlq"` // Added DLQ config for Dead Letter Queue
}
