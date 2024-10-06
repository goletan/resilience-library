# Goletan Resilience Library

The Goletan Resilience Library provides reusable resilience mechanisms for distributed systems, such as retries, circuit breakers, and rate limiters. It is designed to be used across different Nemetons within the Goletan ecosystem, ensuring consistent resilience patterns are applied throughout.

## Features

- Retry Mechanism: Implements exponential backoff with jitter to retry failed operations while avoiding the thundering herd problem.
- Rate Limiter: Provides rate limiting capabilities to control request rates and prevent overload.
- Circuit Breaker: Prevents cascading failures by managing dependencies when issues arise.

## Installation

To use this library, add it to your Go module by running: `go install github.com/goletan/resilience`

## Getting Started

The resilience library can be used to implement retry logic, circuit breaking, and rate limiting to protect services from failures.

### Retry Mechanism

The retry mechanism is useful for transient errors that are likely to succeed if retried after a delay.

### Configuration

The resilience library uses domain-specific configuration for each Nemeton. You can use the shared configuration loader from the config library to load resilience-specific configurations.

### Loading Configuration

To load the configuration, use the shared `LoadConfig()` function from the [config](https://github.com/goletan/config) library

## License

This library is released under the Apache 2 License.
