package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-api-starterkit/internal/ai"
)

type InvestigatorService struct {
	AI   *ai.Service
	Repo Repository
}

func NewInvestigatorService(repo Repository, aiService *ai.Service) *InvestigatorService {
	return &InvestigatorService{
		AI:   aiService,
		Repo: repo,
	}
}

func (s *InvestigatorService) Investigate(ctx context.Context, filter Filter) (*InvestigationResult, []AuditLog, error) {
	if s == nil || s.AI == nil || !s.AI.Enabled() {
		return nil, nil, errors.New("ai investigator is not enabled")
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	filter.Page = 1

	logs, _, err := s.Repo.FindAllWithFilter(filter)
	if err != nil {
		return nil, nil, err
	}
	if len(logs) == 0 {
		return nil, nil, errors.New("no audit logs found for investigation")
	}

	input := ai.BuildJSONPrompt(
		"Summarize these audit logs into JSON with keys: summary, timeline, suspicious_signals, recommendations.",
		buildInvestigationContext(logs),
	)
	result, err := s.AI.Generate(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	var parsed InvestigationResult
	if err := json.Unmarshal([]byte(strings.TrimSpace(result.Text)), &parsed); err != nil {
		return nil, nil, fmt.Errorf("failed to parse ai investigation response: %w", err)
	}

	return &parsed, logs, nil
}

func buildInvestigationContext(logs []AuditLog) string {
	lines := make([]string, 0, len(logs)+1)
	lines = append(lines, fmt.Sprintf("Audit log count: %d", len(logs)))
	for _, logEntry := range logs {
		lines = append(lines, fmt.Sprintf(
			"- time=%s action=%s resource=%s status=%s actor_user_id=%s resource_id=%s ip=%s description=%q",
			logEntry.CreatedAt.UTC().Format(time.RFC3339),
			logEntry.Action,
			logEntry.Resource,
			logEntry.Status,
			uintPointerString(logEntry.ActorUserID),
			uintPointerString(logEntry.ResourceID),
			logEntry.IPAddress,
			logEntry.Description,
		))
	}
	return strings.Join(lines, "\n")
}
