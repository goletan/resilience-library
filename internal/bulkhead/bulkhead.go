package bulkhead

import (
	"context"
	"sync"
	"time"

	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

// Bulkhead controls the number of concurrent operations allowed in a specific section of code to limit resource consumption.
type Bulkhead struct {
	capacity    int
	semaphore   chan struct{}
	timeout     time.Duration
	serviceName string
}

var (
	once             sync.Once
	bulkheadInstance *Bulkhead
)

// NewBulkhead creates a new Bulkhead with a given capacity and timeout duration.
func NewBulkhead(cfg *types.ResilienceConfig, serviceName string) *Bulkhead {
	return &Bulkhead{
		capacity:    cfg.Bulkhead.Capacity,
		semaphore:   make(chan struct{}, cfg.Bulkhead.Capacity),
		timeout:     cfg.Bulkhead.Timeout,
		serviceName: serviceName,
	}
}

// Init once a new Bulkhead instance within its config.
func Init(cfg *types.ResilienceConfig, serviceName string) {
	once.Do(func() {
		bulkheadInstance = NewBulkhead(cfg, serviceName)
	})
}

// GetInstance returns current insteance of the Bulkhead.
func GetInstance() *Bulkhead {
	return bulkheadInstance
}

// Execute attempts to acquire a permit and run the given function within the bulkhead's capacity and timeout.
func (b *Bulkhead) Execute(ctx context.Context, fn func() error, fallback func() error, logger *zap.Logger) error {
	select {
	case b.semaphore <- struct{}{}:
		defer func() { <-b.semaphore }()
		logger.Info("Bulkhead permit acquired", zap.String("service", b.serviceName), zap.Int("occupied_slots", b.Usage()), zap.Int("total_capacity", b.capacity))

		errCh := make(chan error, 1)
		go func() {
			errCh <- fn()
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}

	case <-time.After(b.timeout):
		logger.Warn("Bulkhead limit reached", zap.String("service", b.serviceName), zap.Int("occupied_slots", b.Usage()), zap.Int("total_capacity", b.capacity))
		CountBulkheadLimitReached(b.serviceName)
		if fallback != nil {
			return fallback()
		}
		return context.DeadlineExceeded

	case <-ctx.Done():
		return ctx.Err()
	}
}

// Capacity returns the current capacity of the bulkhead.
func (b *Bulkhead) Capacity() int {
	return b.capacity
}

// SetCapacity allows dynamically updating the bulkhead's capacity.
func (b *Bulkhead) SetCapacity(newCapacity int) {
	if newCapacity > b.capacity {
		for i := 0; i < newCapacity-b.capacity; i++ {
			b.semaphore <- struct{}{}
		}
	} else {
		for i := 0; i < b.capacity-newCapacity; i++ {
			<-b.semaphore
		}
	}
	b.capacity = newCapacity
}

// Usage returns the current number of occupied slots.
func (b *Bulkhead) Usage() int {
	return len(b.semaphore)
}
