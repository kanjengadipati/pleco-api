package monitoring

import "context"

// Monitor is the interface for all monitoring implementations
type Monitor interface {
	// CaptureException logs an error and optionally sends to monitoring service
	CaptureException(err error, ctx context.Context)

	// CaptureMessage logs a message with severity level
	CaptureMessage(msg string, level string, ctx context.Context)

	// RecordMetric records a metric/gauge
	RecordMetric(name string, value float64, tags map[string]string)

	// Close cleanup resources (flushes events, closes connections)
	Close() error
}

// MonitorProvider identifies which monitoring backend to use
type MonitorProvider string

const (
	ProviderNone    MonitorProvider = "none"
	ProviderSentry  MonitorProvider = "sentry"
	ProviderDatadog MonitorProvider = "datadog"
)

// Factory function to create appropriate monitor based on env var
func NewMonitor(provider string) (Monitor, error) {
	switch MonitorProvider(provider) {
	case ProviderSentry:
		return NewSentryMonitor()
	case ProviderDatadog:
		return NewDatadogMonitor()
	case ProviderNone:
		fallthrough
	default:
		return &NoOpMonitor{}, nil
	}
}
