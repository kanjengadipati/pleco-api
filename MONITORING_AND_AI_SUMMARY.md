# Pleco API — Monitoring & AI Integration Summary

**Status:** Planning for v1.0.0 release  
**Priority:** Medium (infrastructure setup now, activation later)  
**Estimated Effort:** 45 minutes implementation  
**Timeline:** Add to v1.0.0, enable in v1.1.0+

---

## Executive Summary

Implement configurable monitoring infrastructure with AI-powered error analysis capabilities. Architecture allows:

1. **Zero dependencies today** — NoOp implementation (no performance cost)
2. **Flip-a-switch enablement** — Activate Sentry/Datadog with env var
3. **AI-powered error analysis** — Unique differentiator vs Auth0/Clerk
4. **Phased rollout** — Infrastructure now, AI analysis when needed

---

## Architecture Overview

```
┌─────────────────────────────────────────┐
│         Application Code                │
│      (error handling, middleware)       │
└────────────────┬────────────────────────┘
                 │
                 ▼
        ┌────────────────┐
        │    Monitor     │  Interface (configurable)
        │   Interface    │
        └────────┬───────┘
                 │
        ┌────────┴─────────────────────┐
        │                              │
        ▼                              ▼
    ┌──────────┐                  ┌──────────────┐
    │ No-Op    │                  │  Real        │
    │ Monitor  │                  │  Monitor     │
    │          │                  │  (Sentry/    │
    │ (v1.0.0) │                  │   Datadog)   │
    └──────────┘                  └───────┬──────┘
                                          │
                                          ▼
                                  ┌──────────────┐
                                  │  AI Monitor  │
                                  │  Wrapper     │
                                  │  (Optional)  │
                                  └───────┬──────┘
                                          │
                                          ▼
                                  ┌──────────────┐
                                  │ AI Service   │
                                  │ (Ollama/     │
                                  │ OpenAI/      │
                                  │ Gemini)      │
                                  └──────────────┘
```

---

## 1. Core Monitoring Interface

**File:** `internal/services/monitoring/monitor.go`

```go
package monitoring

import "context"

// Monitor is the interface for all monitoring implementations
type Monitor interface {
	// CaptureException logs an error and optionally sends to monitoring service
	CaptureException(err error, ctx context.Context)

	// CaptureMessage logs a message with severity level
	CaptureMessage(msg string, level string, ctx context.Context)

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
```

---

## 2. No-Op Implementation (Default)

**File:** `internal/services/monitoring/noop.go`

```go
package monitoring

import "context"

// NoOpMonitor implements Monitor but does nothing
// Used by default to avoid overhead when monitoring is disabled
type NoOpMonitor struct{}

func (m *NoOpMonitor) CaptureException(err error, ctx context.Context) {
	// No-op
}

func (m *NoOpMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	// No-op
}

func (m *NoOpMonitor) Close() error {
	return nil
}
```

---

## 3. Sentry Implementation (Stub)

**File:** `internal/services/monitoring/sentry.go`

```go
package monitoring

import (
	"context"
	"errors"
	"fmt"
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
	m.client.CaptureException(err, &sentry.EventHint{Context: ctx}, nil)
}

func (m *SentryMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	severity := sentry.LevelInfo
	if level == "error" {
		severity = sentry.LevelError
	} else if level == "warning" {
		severity = sentry.LevelWarning
	}

	m.client.CaptureMessage(msg, severity, &sentry.EventHint{Context: ctx}, nil)
}

func (m *SentryMonitor) Close() error {
	m.client.Flush(2 * time.Second)
	return nil
}
```

---

## 4. Datadog Implementation (Stub)

**File:** `internal/services/monitoring/datadog.go`

```go
package monitoring

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type DatadogMonitor struct {
	client statsd.Client
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
	m.client.Event(&statsd.Event{
		Title: "Exception",
		Text:  err.Error(),
		Tags:  []string{"service:pleco-api"},
	})
}

func (m *DatadogMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	m.client.Event(&statsd.Event{
		Title: level,
		Text:  msg,
		Tags:  []string{"service:pleco-api"},
	})
}

func (m *DatadogMonitor) Close() error {
	return m.client.Close()
}
```

---

## 5. AI-Powered Monitor Wrapper

**File:** `internal/services/monitoring/ai_monitor.go`

```go
package monitoring

import (
	"context"
	"fmt"
	"log"
	"os"
	"pleco-api/internal/services"
)

// AIMonitor wraps a base monitor and adds AI-powered error analysis
type AIMonitor struct {
	baseMonitor Monitor
	aiService   *services.AIService
	enabled     bool
	errorCount  int
	threshold   int // Only analyze after N errors
}

func NewAIMonitor(
	baseMonitor Monitor,
	aiService *services.AIService,
	enabled bool,
) *AIMonitor {
	threshold := 5 // Default: analyze after 5 errors
	if t := os.Getenv("AI_MONITORING_ERROR_THRESHOLD"); t != "" {
		fmt.Sscanf(t, "%d", &threshold)
	}

	return &AIMonitor{
		baseMonitor: baseMonitor,
		aiService:   aiService,
		enabled:     enabled,
		threshold:   threshold,
	}
}

func (m *AIMonitor) CaptureException(err error, ctx context.Context) {
	// Always send to base monitor
	m.baseMonitor.CaptureException(err, ctx)

	// Optionally analyze with AI
	if m.enabled && m.aiService != nil {
		m.errorCount++
		if m.errorCount >= m.threshold {
			go m.analyzeError(err, ctx) // Non-blocking
			m.errorCount = 0              // Reset counter
		}
	}
}

func (m *AIMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	m.baseMonitor.CaptureMessage(msg, level, ctx)
}

func (m *AIMonitor) Close() error {
	return m.baseMonitor.Close()
}

// analyzeError uses AI to provide root cause analysis
func (m *AIMonitor) analyzeError(err error, ctx context.Context) {
	prompt := fmt.Sprintf(`
You are an expert debugging assistant for an authentication API built with Go.

Analyze this error and provide:
1. Root cause hypothesis (1-2 sentences)
2. Affected system components
3. Recommended immediate action
4. Severity level (low/medium/high/critical)

Error Message: %s
Error Type: %T

Respond in concise, actionable format.
`, err.Error(), err)

	response, err := m.aiService.Complete(ctx, prompt, "")
	if err != nil {
		log.Printf("AI error analysis failed: %v", err)
		return
	}

	log.Printf("[AI Error Analysis] %s", response)
	// Future: send to Slack, email, dashboard, etc
}
```

---

## 6. Integration into Server Setup

**File:** `internal/appsetup/server.go` (modifications)

```go
func RunAPI(registerDocs func(*gin.Engine)) error {
	config.LoadEnv()
	appConfig := config.LoadAppConfig()
	db := config.ConnectDB(appConfig.DatabaseURL)
	RunStartupTasks(appConfig, db)

	jwtService := services.NewJWTService(appConfig.JWTSecret)

	// Initialize monitoring
	provider := os.Getenv("MONITORING_PROVIDER") // default: "none"
	baseMonitor, err := NewMonitor(provider)
	if err != nil {
		log.Printf("Warning: monitoring initialization failed, falling back to no-op: %v", err)
		baseMonitor = &NoOpMonitor{}
	}
	defer baseMonitor.Close()

	// Wrap with AI monitoring if enabled
	var monitor Monitor = baseMonitor
	if os.Getenv("AI_MONITORING_ENABLED") == "true" {
		aiService := services.NewAIService(appConfig) // Reuse existing AI service
		monitor = NewAIMonitor(baseMonitor, aiService, true)
	}

	// Build router with monitor middleware
	router, err := BuildRouter(db, appConfig, jwtService)
	if err != nil {
		return err
	}

	// Add monitoring middleware for error capture
	router.Use(func(c *gin.Context) {
		c.Next()

		// Capture 5xx errors
		if c.Writer.Status() >= 500 {
			err := fmt.Errorf(
				"HTTP %d on %s %s",
				c.Writer.Status(),
				c.Request.Method,
				c.Request.URL.Path,
			)
			monitor.CaptureException(err, c.Request.Context())
		}
	})

	// ... rest of server setup
}
```

---

## 7. Environment Configuration

**File:** `.env.example` (additions)

```env
# === MONITORING & OBSERVABILITY ===

# Monitoring Provider: none | sentry | datadog
MONITORING_PROVIDER=none

# Sentry (required if MONITORING_PROVIDER=sentry)
SENTRY_DSN=

# Datadog (required if MONITORING_PROVIDER=datadog)
DATADOG_API_KEY=
DATADOG_ENV=production
DATADOG_SERVICE=pleco-api

# === AI-POWERED MONITORING (Optional) ===

# Enable AI error analysis (requires MONITORING_PROVIDER to be set)
AI_MONITORING_ENABLED=false

# AI provider for monitoring: none | mock | ollama | openai | gemini
AI_MONITORING_PROVIDER=mock
AI_MONITORING_MODEL=mock-model

# Ollama configuration (required if AI_MONITORING_PROVIDER=ollama)
AI_MONITORING_BASE_URL=http://localhost:11434

# OpenAI/Gemini configuration (required for those providers)
AI_MONITORING_API_KEY=

# Error threshold: analyze after N errors
AI_MONITORING_ERROR_THRESHOLD=5
```

---

## 8. README Section

**Add to README:**

```markdown
## Monitoring & Observability

Pleco includes optional monitoring with AI-powered error analysis capabilities.

### Basic Monitoring

Monitor errors with Sentry or Datadog:

```bash
# Using Sentry
export MONITORING_PROVIDER=sentry
export SENTRY_DSN=https://key@sentry.io/project
go run ./cmd/api
```

```bash
# Using Datadog
export MONITORING_PROVIDER=datadog
export DATADOG_API_KEY=your_key
go run ./cmd/api
```

### AI-Powered Error Analysis

When enabled, AI analyzes error patterns and provides root cause analysis:

```bash
export MONITORING_PROVIDER=sentry
export SENTRY_DSN=...
export AI_MONITORING_ENABLED=true
export AI_MONITORING_PROVIDER=ollama
export AI_MONITORING_MODEL=qwen2.5:3b
export AI_MONITORING_BASE_URL=http://localhost:11434
```

AI analysis currently logs to stdout. Future: Slack/email integration.

### Current Status

- v1.0.0: Infrastructure in place, disabled by default (zero overhead)
- v1.1.0: Sentry/Datadog integration
- v1.2.0: AI error analysis (mock/ollama)
- v1.3.0+: Real AI providers + smart alerts
```

---

## 9. Implementation Checklist

### Phase 1: Infrastructure (v1.0.0) — 45 minutes

```
[ ] Create internal/services/monitoring/ package
[ ] Create monitor.go (interface + factory)
[ ] Create noop.go (no-op implementation)
[ ] Create sentry.go (stub implementation)
[ ] Create datadog.go (stub implementation)
[ ] Create ai_monitor.go (AI wrapper)
[ ] Add monitoring initialization to appsetup/server.go
[ ] Add monitoring middleware for error capture
[ ] Update .env.example with new variables
[ ] Update README with monitoring section
[ ] Test no-op behavior (default case)
[ ] Add to go.mod: github.com/getsentry/sentry-go (as optional)
[ ] Add to go.mod: github.com/DataDog/datadog-go/v5 (as optional)
```

**Result:** Monitoring interfaces defined, no external dependencies yet, zero overhead.

### Phase 2: Sentry Integration (v1.1.0) — When you have users

```
[ ] go get github.com/getsentry/sentry-go
[ ] Test Sentry integration locally
[ ] Create Sentry project + get DSN
[ ] Set MONITORING_PROVIDER=sentry in production
[ ] Verify errors flowing to Sentry dashboard
[ ] Setup Sentry alerts
```

### Phase 3: AI Error Analysis (v1.2.0) — With paying customers

```
[ ] Reuse internal/services/ai from audit investigator
[ ] Enable AI_MONITORING_ENABLED=true
[ ] Setup Ollama or use OpenAI provider
[ ] Test AI analysis on common error patterns
[ ] Improve prompt for better analysis quality
```

### Phase 4: Advanced Features (v1.3.0+)

```
[ ] Slack/email integration for alerts
[ ] Anomaly detection (error rate spikes)
[ ] Pattern detection (repeated errors)
[ ] Predictive alerts (based on trends)
[ ] Custom dashboard
```

---

## 10. Key Design Decisions

### Why Interface-based?

- **Testability:** Easy to mock Monitor in tests
- **Flexibility:** Can swap Sentry ↔ Datadog without code changes
- **Extensibility:** Easy to add custom implementations
- **Performance:** No-op implementation has zero overhead

### Why wrap with AI Monitor?

- **Separation of concerns:** AI analysis is optional layer
- **Reusability:** Can reuse existing AIService
- **Non-blocking:** AI analysis runs async (doesn't slow down API)
- **Future-proof:** Easy to add more analyzers

### Why error threshold?

- **Cost control:** Only analyze every Nth error to avoid API calls
- **Meaningful insights:** Wait for patterns to emerge before analyzing
- **Configurable:** Can adjust for different environments

---

## 11. Testing Strategy

```go
// Test that no-op monitor doesn't affect performance
func TestNoOpMonitor_NoOverhead(t *testing.T) {
	monitor := &NoOpMonitor{}
	start := time.Now()
	for i := 0; i < 10000; i++ {
		monitor.CaptureException(errors.New("test"), context.Background())
	}
	elapsed := time.Since(start)
	assert.Less(t, elapsed, 100*time.Millisecond) // Should be very fast
}

// Test monitoring initialization
func TestMonitorInitialization(t *testing.T) {
	tests := []struct {
		provider string
		expectErr bool
	}{
		{"none", false},
		{"sentry", true}, // No SENTRY_DSN env var
		{"datadog", true}, // No DATADOG_API_KEY env var
	}

	for _, tt := range tests {
		monitor, err := NewMonitor(tt.provider)
		if tt.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, monitor)
		}
	}
}

// Test AI Monitor doesn't break if AI service unavailable
func TestAIMonitor_Graceful_Degradation(t *testing.T) {
	baseMonitor := &NoOpMonitor{}
	aiMonitor := NewAIMonitor(baseMonitor, nil, true)
	
	// Should not panic even with nil AI service
	aiMonitor.CaptureException(errors.New("test"), context.Background())
}
```

---

## 12. Deployment Notes

### Local Development

```bash
# Default: no monitoring overhead
make docker-up
# or
go run ./cmd/api
```

### Staging (with Sentry)

```bash
export MONITORING_PROVIDER=sentry
export SENTRY_DSN=<staging-sentry-dsn>
go run ./cmd/api
```

### Production (with Sentry + AI)

```bash
export MONITORING_PROVIDER=sentry
export SENTRY_DSN=<prod-sentry-dsn>
export AI_MONITORING_ENABLED=true
export AI_MONITORING_PROVIDER=ollama
export AI_MONITORING_BASE_URL=<ollama-url>
export AI_MONITORING_ERROR_THRESHOLD=5
go run ./cmd/api
```

---

## 13. Unique Selling Points

This architecture makes Pleco **unique in the market:**

```
Traditional Auth Backends:
- Error logging: ✅
- Error monitoring: ✅
- AI-powered analysis: ❌

Pleco:
- Error logging: ✅
- Error monitoring: ✅
- AI-powered analysis: ✅ (unique)
- Auditable AI reasoning: ✅ (future)
```

Marketing angle: **"The only auth backend with AI-assisted debugging"**

---

## 14. Future Enhancements

```go
// Future: Anomaly detection
type AnomalyDetector struct {
	errorRateMean   float64
	errorRateStdDev float64
}

func (d *AnomalyDetector) Detect(currentRate float64) bool {
	zScore := (currentRate - d.errorRateMean) / d.errorRateStdDev
	return zScore > 2.0 // 2 std devs = anomaly
}

// Future: Pattern clustering
func (m *AIMonitor) DetectPatterns(window time.Duration) {
	errors := m.getRecentErrors(window)
	patterns := cluster(errors) // Group similar errors
	
	for _, pattern := range patterns {
		if pattern.Frequency > threshold {
			m.analyzePattern(pattern)
		}
	}
}

// Future: Predictive alerts
func (m *AIMonitor) PredictIssues(ctx context.Context) {
	trends := m.getErrorTrends(time.Hour)
	prompt := fmt.Sprintf(
		"Based on these error trends, predict issues in next 30 min: %v",
		trends,
	)
	prediction, _ := m.aiService.Complete(ctx, prompt, "")
	log.Printf("[Prediction] %s", prediction)
}
```

---

## Summary

**What to build:** Configurable monitoring infrastructure with optional AI integration  
**Effort:** 45 minutes  
**Dependencies:** Minimal (interfaces only, implementations are stubs)  
**Go-live:** v1.0.0 with zero overhead (default: no-op)  
**Activation:** Flip env var when needed (v1.1.0+)  
**Differentiation:** Only auth backend with AI-powered error analysis  

**Status:** Ready for handoff to engineer ✅

---

## Questions for Engineer

Before starting:

1. Should we add structured logging (JSON) to complement monitoring?
2. Should monitoring middleware also capture 4xx errors or just 5xx?
3. Should AI analysis results be stored in database for dashboard view?
4. Do you want metrics/gauge support (latency, throughput) or just errors?
5. Should we add request sampling (only analyze 1 in N requests)?
