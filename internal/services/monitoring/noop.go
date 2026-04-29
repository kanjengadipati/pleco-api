package monitoring

import (
	"context"
	"log/slog"
)

// NoOpMonitor implements Monitor but does nothing except optionally logging
// Used by default to avoid overhead when monitoring is disabled
type NoOpMonitor struct{}

func (m *NoOpMonitor) CaptureException(err error, ctx context.Context) {
	slog.Error("Captured Exception", "error", err)
}

func (m *NoOpMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	slog.Info("Captured Message", "msg", msg, "level", level)
}

func (m *NoOpMonitor) RecordMetric(name string, value float64, tags map[string]string) {
	slog.Debug("Recorded Metric", "name", name, "value", value, "tags", tags)
}

func (m *NoOpMonitor) Close() error {
	return nil
}
