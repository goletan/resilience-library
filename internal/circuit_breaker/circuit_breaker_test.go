// /circuit_breaker/circuit_breaker_test.go
package circuit_breaker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goletan/resilience/internal/types"
)

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cfg := &types.ResilienceConfig{
		CircuitBreaker: struct {
			FailureRateThreshold float32       `mapstructure:"failure_rate_threshold"`
			FailureThreshold     int           `mapstructure:"failure_threshold"`
			SuccessThreshold     int           `mapstructure:"success_threshold"`
			ConsecutiveFailures  int           `mapstructure:"consecutive_failures"`
			MaxRequest           uint32        `mapstructure:"max_requests"`
			Interval             time.Duration `mapstructure:"interval"`
			Timeout              time.Duration `mapstructure:"timeout"`
		}{
			FailureRateThreshold: 0.5,
			ConsecutiveFailures:  3,
			MaxRequest:           5,
			Interval:             60 * time.Second,
			Timeout:              10 * time.Second,
		},
	}

	cb := NewCircuitBreaker(cfg)
	err := cb.Execute(context.Background(), func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCircuitBreaker_Execute_FailureAndOpen(t *testing.T) {
	cfg := &types.ResilienceConfig{
		CircuitBreaker: struct {
			FailureRateThreshold float32       `mapstructure:"failure_rate_threshold"`
			FailureThreshold     int           `mapstructure:"failure_threshold"`
			SuccessThreshold     int           `mapstructure:"success_threshold"`
			ConsecutiveFailures  int           `mapstructure:"consecutive_failures"`
			MaxRequest           uint32        `mapstructure:"max_requests"`
			Interval             time.Duration `mapstructure:"interval"`
			Timeout              time.Duration `mapstructure:"timeout"`
		}{
			FailureRateThreshold: 0.5,
			ConsecutiveFailures:  3,
			MaxRequest:           5,
			Interval:             60 * time.Second,
			Timeout:              10 * time.Second,
		},
	}

	cb := NewCircuitBreaker(cfg)

	// Trigger consecutive failures
	for i := 0; i < 3; i++ {
		err := cb.Execute(context.Background(), func() error {
			return errors.New("operation failed")
		})
		if err == nil {
			t.Errorf("Expected error, got none")
		}
	}

	// Next attempt should be blocked by the open circuit breaker
	err := cb.Execute(context.Background(), func() error {
		return nil
	})
	if err == nil {
		t.Errorf("Expected error due to open circuit breaker, got none")
	}
}
