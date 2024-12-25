# Goletan Resilience Library

The Goletan Resilience Library provides reusable resilience mechanisms for distributed systems, such as retries, circuit breakers, bulkheads, and rate limiters. It is designed to be used across different Nemetons within the Goletan ecosystem, ensuring consistent resilience patterns are applied throughout.

## Features

- **Retry Mechanism**: Implements exponential backoff with jitter to retry failed operations while avoiding the thundering herd problem.
- **Rate Limiter**: Provides rate limiting capabilities to control request rates and prevent overload.
- **Circuit Breaker**: Prevents cascading failures by managing dependencies when issues arise, including state transitions and fallback mechanisms.
- **Bulkhead**: Limits the number of concurrent operations to prevent resource exhaustion, ensuring that critical services remain available.

## Installation

To use this library, add it to your Go module by running: `go get github.com/goletan/resilience`

## Getting Started

The resilience library can be used to implement retry logic, circuit breaking, bulkheading, and rate limiting to protect services from failures and ensure reliability.

### Usage Examples

#### Retry Mechanism
The retry mechanism is useful for transient errors that are likely to succeed if retried after a delay. Configure and use it to add resilience to operations prone to intermittent failures.

```go
import (
    "context"
    "github.com/goletan/resilience-library/retry"
)

func main() {
    ctx := context.Background()
    err := retry.ExecuteWithRetry(ctx, myOperation, 3, shouldRetry)
    if err != nil {
        // Handle failure after retries
    }
}
```

#### Circuit Breaker
The circuit breaker helps manage downstream dependencies by preventing repeated failures from overwhelming the system.

```go
import (
    "context"
    "github.com/goletan/resilience-library/circuit_breaker"
    "github.com/goletan/resilience-library/types"
)

func main() {
    cfg := types.LoadDefaultCircuitBreakerConfig()
    cb := circuit_breaker.NewCircuitBreaker(cfg, nil)

    err := cb.Execute(context.Background(), myOperation, fallbackOperation)
    if err != nil {
        // Handle circuit breaker open state
    }
}
```

#### Rate Limiter
The rate limiter enforces limits on how many operations can be performed within a given timeframe.

```go
import (
    "context"
    "github.com/goletan/resilience-library/rate_limiter"
)

func main() {
    rateLimiter := rate_limiter.NewRateLimiter("operationName", cfg, logger)
    err := rateLimiter.ExecuteWithRateLimiting(context.Background(), myOperation)
    if err != nil {
        // Handle rate limit exceeded
    }
}
```

#### Bulkhead
Bulkheading limits concurrent operations to avoid resource exhaustion.

```go
import (
    "context"
    "github.com/goletan/resilience-library/bulkhead"
)

func main() {
    bulkhead := bulkhead.NewBulkhead(cfg, "serviceName")
    err := bulkhead.Execute(context.Background(), myOperation, fallbackOperation, logger)
    if err != nil {
        // Handle bulkhead limit reached
    }
}
```

### Configuration

The resilience library uses domain-specific configuration for each Nemeton. You can use the shared configuration loader from the config library to load resilience-specific configurations. An example configuration file (`config.yaml`) may include:

```yaml
resilience:
  max_retries: 3
  base_backoff: 100ms
  bulkhead:
    capacity: 5
    timeout: 10s
  circuit_breaker:
    failure_rate_threshold: 0.5
    failure_threshold: 5
    success_threshold: 2
    consecutive_failures: 3
    max_requests: 5
    interval: 60s
    timeout: 30s
  rate_limiter:
    rps: 100
    burst: 10
```

## License

This library is released under the Apache 2 License.
