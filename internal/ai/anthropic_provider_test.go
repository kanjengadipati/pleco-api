package ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnthropicProviderGenerateSuccess(t *testing.T) {
	provider := NewAnthropicProvider("https://api.anthropic.test", "claude-key", 5*time.Second)
	provider.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/v1/messages", req.URL.Path)
			assert.Equal(t, "claude-key", req.Header.Get("x-api-key"))
			assert.Equal(t, anthropicVersion, req.Header.Get("anthropic-version"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"content":[
						{"type":"text","text":"{\"summary\":\"ok\",\"timeline\":[],\"suspicious_signals\":[],\"recommendations\":[]}"}
					]
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	}

	result, err := provider.Generate(context.Background(), GenerateInput{
		Model:        "claude-sonnet-4-5",
		SystemPrompt: "system",
		UserPrompt:   "user",
		MaxTokens:    300,
	})

	assert.NoError(t, err)
	assert.Contains(t, result.Text, `"summary":"ok"`)
}

func TestAnthropicProviderGenerateReturnsAPIError(t *testing.T) {
	provider := NewAnthropicProvider("https://api.anthropic.test", "claude-key", 5*time.Second)
	provider.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"bad api key"}}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	result, err := provider.Generate(context.Background(), GenerateInput{
		Model:      "claude-sonnet-4-5",
		UserPrompt: "hello",
		MaxTokens:  300,
	})

	assert.Nil(t, result)
	assert.EqualError(t, err, "anthropic error: bad api key")
}
