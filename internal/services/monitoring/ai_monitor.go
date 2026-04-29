package monitoring

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"pleco-api/internal/ai"
)

// ErrorAnalysis represents an AI analysis of an error, stored in the database
type ErrorAnalysis struct {
	gorm.Model
	ErrorMessage       string
	ErrorType          string
	RootCause          string
	AffectedComponents string
	RecommendedAction  string
	Severity           string
}

// AIMonitor wraps a base monitor and adds AI-powered error analysis
type AIMonitor struct {
	baseMonitor Monitor
	aiService   *ai.Service
	db          *gorm.DB
	enabled     bool
	errorCount  int
	threshold   int // Only analyze after N errors
}

func NewAIMonitor(
	baseMonitor Monitor,
	aiService *ai.Service,
	db *gorm.DB,
	enabled bool,
) *AIMonitor {
	threshold := 5 // Default: analyze after 5 errors
	if t := os.Getenv("AI_MONITORING_ERROR_THRESHOLD"); t != "" {
		fmt.Sscanf(t, "%d", &threshold)
	}

	return &AIMonitor{
		baseMonitor: baseMonitor,
		aiService:   aiService,
		db:          db,
		enabled:     enabled,
		threshold:   threshold,
	}
}

func (m *AIMonitor) CaptureException(err error, ctx context.Context) {
	// Always send to base monitor
	m.baseMonitor.CaptureException(err, ctx)

	// Optionally analyze with AI
	if m.enabled && m.aiService != nil && m.db != nil {
		m.errorCount++
		// Sampling: only analyze every Nth error
		if m.errorCount >= m.threshold {
			go m.analyzeError(err, ctx) // Non-blocking
			m.errorCount = 0            // Reset counter
		}
	}
}

func (m *AIMonitor) CaptureMessage(msg string, level string, ctx context.Context) {
	m.baseMonitor.CaptureMessage(msg, level, ctx)
}

func (m *AIMonitor) RecordMetric(name string, value float64, tags map[string]string) {
	m.baseMonitor.RecordMetric(name, value, tags)
}

func (m *AIMonitor) Close() error {
	return m.baseMonitor.Close()
}

// analyzeError uses AI to provide root cause analysis and stores it in DB
func (m *AIMonitor) analyzeError(err error, ctx context.Context) {
	prompt := fmt.Sprintf(`Analyze this error and provide:
1. Root cause hypothesis (1-2 sentences)
2. Affected system components
3. Recommended immediate action
4. Severity level (low/medium/high/critical)

Format your response EXACTLY as follows, with no extra text:
ROOT_CAUSE: <your answer>
COMPONENTS: <your answer>
ACTION: <your answer>
SEVERITY: <your answer>

Error Message: %s
Error Type: %T
`, err.Error(), err)

	// We use a detached context because the original request context might be canceled
	// when the HTTP request finishes, but analysis takes longer.
	analysisCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := ai.GenerateInput{
		SystemPrompt: "You are an expert debugging assistant for an authentication API built with Go.",
		UserPrompt:   prompt,
	}

	response, aiErr := m.aiService.Generate(analysisCtx, input)
	if aiErr != nil {
		slog.Error("AI error analysis failed", "error", aiErr)
		return
	}

	analysis := parseAIResponse(err, response.Text)

	if dbErr := m.db.Create(&analysis).Error; dbErr != nil {
		slog.Error("Failed to save AI error analysis to DB", "error", dbErr)
		return
	}

	slog.Info("AI Error Analysis completed and saved", "analysis_id", analysis.ID)
}

func parseAIResponse(err error, response string) ErrorAnalysis {
	analysis := ErrorAnalysis{
		ErrorMessage: err.Error(),
		ErrorType:    fmt.Sprintf("%T", err),
	}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ROOT_CAUSE:") {
			analysis.RootCause = strings.TrimSpace(strings.TrimPrefix(line, "ROOT_CAUSE:"))
		} else if strings.HasPrefix(line, "COMPONENTS:") {
			analysis.AffectedComponents = strings.TrimSpace(strings.TrimPrefix(line, "COMPONENTS:"))
		} else if strings.HasPrefix(line, "ACTION:") {
			analysis.RecommendedAction = strings.TrimSpace(strings.TrimPrefix(line, "ACTION:"))
		} else if strings.HasPrefix(line, "SEVERITY:") {
			analysis.Severity = strings.TrimSpace(strings.TrimPrefix(line, "SEVERITY:"))
		}
	}
	return analysis
}
