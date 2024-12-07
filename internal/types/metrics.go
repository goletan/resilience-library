package types

// Metrics is a common interface to register metrics in Prometheus.
type Metrics interface {
	Register() error
}
