package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// remoteProvider calls a real LLM HTTP API. kind selects the wire format.
type remoteProvider struct {
	kind string
}

func (p remoteProvider) Name() string { return p.kind }

func (p remoteProvider) Run(req Request) (Result, error) {
	switch p.kind {
	case "openai":
		return p.runOpenAI(req)
	case "claude":
		return p.runClaude(req)
	case "ollama":
		return p.runOllama(req)
	case "gemini":
		return p.runGemini(req)
	default:
		return Result{}, fmt.Errorf("unknown provider %q", p.kind)
	}
}

func (p remoteProvider) runOpenAI(req Request) (Result, error) {
	if req.APIKey == "" {
		return Result{}, errors.New("openai: api_key is required")
	}
	model := def(req.Model, "gpt-4o-mini")
	body, err := postJSON(
		def(req.BaseURL, "https://api.openai.com")+"/v1/chat/completions",
		map[string]string{"Authorization": "Bearer " + req.APIKey},
		map[string]any{
			"model":    model,
			"messages": []map[string]string{{"role": "user", "content": req.Prompt}},
		},
	)
	if err != nil {
		return Result{}, err
	}
	return Result{Provider: "openai", Model: model,
		Output: dig(body, "choices", 0, "message", "content")}, nil
}

func (p remoteProvider) runClaude(req Request) (Result, error) {
	if req.APIKey == "" {
		return Result{}, errors.New("claude: api_key is required")
	}
	model := def(req.Model, "claude-sonnet-4-6")
	body, err := postJSON(
		def(req.BaseURL, "https://api.anthropic.com")+"/v1/messages",
		map[string]string{"x-api-key": req.APIKey, "anthropic-version": "2023-06-01"},
		map[string]any{
			"model":      model,
			"max_tokens": 1024,
			"messages":   []map[string]string{{"role": "user", "content": req.Prompt}},
		},
	)
	if err != nil {
		return Result{}, err
	}
	return Result{Provider: "claude", Model: model,
		Output: dig(body, "content", 0, "text")}, nil
}

func (p remoteProvider) runOllama(req Request) (Result, error) {
	model := def(req.Model, "llama3")
	body, err := postJSON(
		def(req.BaseURL, "http://localhost:11434")+"/api/generate",
		nil,
		map[string]any{"model": model, "prompt": req.Prompt, "stream": false},
	)
	if err != nil {
		return Result{}, err
	}
	return Result{Provider: "ollama", Model: model,
		Output: dig(body, "response")}, nil
}

func (p remoteProvider) runGemini(req Request) (Result, error) {
	if req.APIKey == "" {
		return Result{}, errors.New("gemini: api_key is required")
	}
	model := def(req.Model, "gemini-1.5-flash")
	base := def(req.BaseURL, "https://generativelanguage.googleapis.com")
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", base, model, req.APIKey)
	body, err := postJSON(url, nil, map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": req.Prompt}}},
		},
	})
	if err != nil {
		return Result{}, err
	}
	return Result{Provider: "gemini", Model: model,
		Output: dig(body, "candidates", 0, "content", "parts", 0, "text")}, nil
}

// postJSON sends a JSON POST and decodes the JSON response body.
func postJSON(url string, headers map[string]string, payload any) (map[string]any, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}
	return out, nil
}

// dig walks nested decoded JSON (maps keyed by string, arrays indexed by int)
// and returns the string at the path, or "" if the path does not resolve.
func dig(v any, path ...any) string {
	for _, step := range path {
		switch key := step.(type) {
		case string:
			m, ok := v.(map[string]any)
			if !ok {
				return ""
			}
			v = m[key]
		case int:
			a, ok := v.([]any)
			if !ok || key < 0 || key >= len(a) {
				return ""
			}
			v = a[key]
		}
	}
	s, _ := v.(string)
	return s
}
