// /resilience/types/resilience.go
package types

import (
	"context"
)

// ResilienceService defines methods for executing operations with resilience mechanisms.
type ResilienceService interface {
	ExecuteWithResilience(ctx context.Context, operation func() error) error
}

type ResilienceMetrics struct{}
