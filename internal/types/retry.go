package types

import "context"

type RetryPolicyInterface interface {
	ExecuteWithRetry(ctx context.Context, operation func() error) error
}
