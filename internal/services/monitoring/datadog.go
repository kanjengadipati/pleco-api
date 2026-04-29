package monitoring

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type DatadogMonitor struct {
	client statsd.ClientInterface
}

func NewDatadogMonitor() (*DatadogMonitor, error) {
	apiKey := os.Getenv("DATADOG_API_KEY")
	if apiKey == "" {
		return nil, errors.New("DATADOG_API_KEY not set when provider=datadog")
	}

	client, err := statsd.New("localhost:8125")
	if err != nil {
		return nil, fmt.Errorf("failed to create Datadog client: %w", err)
	}

	return &DatadogMonitor{client: client}, nil
}

func (m *DatadogMonitor) CaptureException(err error, ctx context.Context) {
	slog.Error("Datadog: Exception captured", "error", err)
	m.client.Event(&statsd.Event{
		Title: "Exception",
		Text:  err.Error(),
		Tags:  []string{"service:pleco-api"},
	})
}

func (m *DatadogMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	slog.Info("Datadog: Message captured", "msg", msg, "level", level)
	m.client.Event(&statsd.Event{
		Title: level,
		Text:  msg,
		Tags:  []string{"service:pleco-api"},
	})
}

func (m *DatadogMonitor) RecordMetric(name string, value float64, tags map[string]string) {
	slog.Debug("Datadog: Metric recorded", "name", name, "value", value, "tags", tags)

	var statsdTags []string
	for k, v := range tags {
		statsdTags = append(statsdTags, fmt.Sprintf("%s:%s", k, v))
	}

	m.client.Gauge(name, value, statsdTags, 1)
}

func (m *DatadogMonitor) Close() error {
	return m.client.Close()
}
