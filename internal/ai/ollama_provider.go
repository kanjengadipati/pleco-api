package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OllamaProvider struct {
	baseURL string
	client  *http.Client
}

type ollamaRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	System      string  `json:"system,omitempty"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature,omitempty"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

func NewOllamaProvider(baseURL string) *OllamaProvider {
	return &OllamaProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *OllamaProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error) {
	body, err := json.Marshal(ollamaRequest{
		Model:       input.Model,
		Prompt:      input.UserPrompt,
		System:      input.SystemPrompt,
		Stream:      false,
		Temperature: input.Temperature,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsed ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if parsed.Error != "" {
			return nil, fmt.Errorf("ollama error: %s", parsed.Error)
		}
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}
	if parsed.Error != "" {
		return nil, fmt.Errorf("ollama error: %s", parsed.Error)
	}

	return &GenerateResult{Text: parsed.Response}, nil
}
