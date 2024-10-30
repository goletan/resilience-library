// /bulkhead/bulkhead_test.go
package bulkhead

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goletan/resilience/types"
	"go.uber.org/zap"
)

func TestBulkhead_Execute_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	b := NewBulkhead(&types.ResilienceConfig{
		Bulkhead: struct {
			Capacity int           `mapstructure:"capacity"`
			Timeout  time.Duration `mapstructure:"timeout"`
		}{
			Capacity: 2,
			Timeout:  1 * time.Second,
		},
	}, "bulkhead_test")

	err := b.Execute(
		context.Background(),
		func() error {
			return nil
		},
		func() error {
			return errors.New("fallback executed")
		},
		logger,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBulkhead_Execute_Timeout(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create a Bulkhead with a capacity of 1 and a timeout of 100 milliseconds.
	b := NewBulkhead(&types.ResilienceConfig{
		Bulkhead: struct {
			Capacity int           `mapstructure:"capacity"`
			Timeout  time.Duration `mapstructure:"timeout"`
		}{
			Capacity: 1,
			Timeout:  100 * time.Millisecond,
		},
	}, "bulkhead_test")

	// First execution acquires the permit, runs for 500 milliseconds.
	go func() {
		err := b.Execute(context.Background(), func() error {
			time.Sleep(500 * time.Millisecond)
			return nil
		}, nil, logger)

		if err != nil {
			t.Errorf("Expected no error for first execution, got %v", err)
		}
	}()

	// Give some time to let the first execution acquire the semaphore.
	time.Sleep(10 * time.Millisecond)

	// Second execution should timeout and fallback should be executed.
	err := b.Execute(
		context.Background(),
		func() error {
			return nil
		},
		func() error {
			return errors.New("fallback executed")
		},
		logger,
	)

	if err == nil || err.Error() != "fallback executed" {
		t.Errorf("Expected fallback executed error, got %v", err)
	}
}

func TestBulkhead_Execute_CancelledContext(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	b := NewBulkhead(&types.ResilienceConfig{
		Bulkhead: struct {
			Capacity int           `mapstructure:"capacity"`
			Timeout  time.Duration `mapstructure:"timeout"`
		}{
			Capacity: 1,
			Timeout:  1 * time.Second,
		},
	}, "bulkhead_test")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	err := b.Execute(
		ctx,
		func() error {
			return nil
		},
		func() error {
			return errors.New("fallback executed")
		},
		logger,
	)

	if err != context.Canceled {
		t.Errorf("Expected context canceled error, got %v", err)
	}
}
