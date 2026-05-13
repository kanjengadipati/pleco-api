package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const defaultAnthropicBaseURL = "https://api.anthropic.com"
const anthropicVersion = "2023-06-01"

type AnthropicProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type anthropicMessagesRequest struct {
	Model       string                    `json:"model"`
	System      string                    `json:"system,omitempty"`
	Messages    []anthropicMessage        `json:"messages"`
	Temperature *float64                  `json:"temperature,omitempty"`
	MaxTokens   int                       `json:"max_tokens"`
	Metadata    *anthropicRequestMetadata `json:"metadata,omitempty"`
}

type anthropicRequestMetadata struct {
	UserID string `json:"user_id,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicMessagesResponse struct {
	Type    string                  `json:"type"`
	Error   *anthropicErrorItem     `json:"error"`
	Content []anthropicContentBlock `json:"content"`
}

type anthropicErrorItem struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewAnthropicProvider(baseURL, apiKey string, timeout time.Duration) *AnthropicProvider {
	return &AnthropicProvider{
		baseURL: normalizeBaseURL(baseURL, defaultAnthropicBaseURL),
		apiKey:  apiKey,
		client:  &http.Client{Timeout: timeout},
	}
}

func (p *AnthropicProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error) {
	reqBody := anthropicMessagesRequest{
		Model: input.Model,
		System: strings.TrimSpace(strings.Join([]string{
			strings.TrimSpace(input.SystemPrompt),
			"Return only valid JSON matching the requested investigation structure.",
		}, "\n\n")),
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: input.UserPrompt,
			},
		},
		MaxTokens: input.MaxTokens,
	}
	if input.Temperature > 0 {
		reqBody.Temperature = &input.Temperature
	}

	body, err := marshalBody(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := p.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		return nil, fmt.Errorf("anthropic is unavailable: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed anthropicMessagesResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to decode anthropic response: %s", strings.TrimSpace(string(bodyBytes)))
	}

	if resp.StatusCode >= 400 {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return nil, fmt.Errorf("anthropic error: %s", parsed.Error.Message)
		}
		return nil, fmt.Errorf("anthropic returned status %d", resp.StatusCode)
	}

	text := extractAnthropicText(parsed.Content)
	if text == "" {
		return nil, ErrInvalidStructuredOutput
	}

	return &GenerateResult{Text: text}, nil
}

func extractAnthropicText(items []anthropicContentBlock) string {
	var parts []string
	for _, item := range items {
		if item.Type == "text" && strings.TrimSpace(item.Text) != "" {
			parts = append(parts, item.Text)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}
