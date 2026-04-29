package monitoring

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryMonitor struct {
	client *sentry.Client
}

func NewSentryMonitor() (*SentryMonitor, error) {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		return nil, errors.New("SENTRY_DSN not set when provider=sentry")
	}

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      os.Getenv("ENV_NAME"),
		Release:          "1.0.0",
		TracesSampleRate: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Sentry client: %w", err)
	}

	return &SentryMonitor{client: client}, nil
}

func (m *SentryMonitor) CaptureException(err error, ctx context.Context) {
	slog.Error("Sentry: Exception captured", "error", err)
	m.client.CaptureException(err, &sentry.EventHint{Context: ctx}, nil)
}

func (m *SentryMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	slog.Info("Sentry: Message captured", "msg", msg, "level", level)
	severity := sentry.LevelInfo
	switch level {
	case "error":
		severity = sentry.LevelError
	case "warning":
		severity = sentry.LevelWarning
	}

	event := sentry.NewEvent()
	event.Message = msg
	event.Level = severity
	m.client.CaptureEvent(event, &sentry.EventHint{Context: ctx}, nil)
}

func (m *SentryMonitor) RecordMetric(name string, value float64, tags map[string]string) {
	slog.Debug("Sentry: Metric recorded", "name", name, "value", value, "tags", tags)
	// Sentry supports metrics but it's often more complex or requires specific config.
	// Placeholder for Sentry metrics integration if needed later.
}

func (m *SentryMonitor) Close() error {
	m.client.Flush(2 * time.Second)
	return nil
}
