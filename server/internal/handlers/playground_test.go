package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPlaygroundMockRun(t *testing.T) {
	r := newTestRouter(t)

	w := do(t, r, http.MethodPost, "/api/playground/run", gin.H{
		"content":   "Hello {{name}}, you are a {{role}}.",
		"variables": gin.H{"name": "Ada", "role": "engineer"},
		"provider":  "mock",
		"model":     "mock-1",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("playground run: want 200, got %d: %s", w.Code, w.Body)
	}

	var resp struct {
		Rendered string `json:"rendered"`
		Result   struct {
			Provider string `json:"provider"`
			Output   string `json:"output"`
		} `json:"result"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Rendered != "Hello Ada, you are a engineer." {
		t.Fatalf("rendered = %q", resp.Rendered)
	}
	if resp.Result.Provider != "mock" {
		t.Fatalf("provider = %q, want mock", resp.Result.Provider)
	}
	if !strings.Contains(resp.Result.Output, "mock completion") {
		t.Fatalf("unexpected output: %q", resp.Result.Output)
	}
}

func TestPlaygroundUnknownProvider(t *testing.T) {
	r := newTestRouter(t)
	w := do(t, r, http.MethodPost, "/api/playground/run", gin.H{
		"content":  "hi",
		"provider": "does-not-exist",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unknown provider: want 400, got %d", w.Code)
	}
}

func TestPlaygroundRunsStoredPrompt(t *testing.T) {
	r := newTestRouter(t)

	w := do(t, r, http.MethodPost, "/api/prompts", gin.H{
		"key": "pg.test", "content": "Translate: {{text}}", "env": "dev",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create: %d %s", w.Code, w.Body)
	}
	var created struct {
		Data struct{ ID string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	w = do(t, r, http.MethodPost, "/api/playground/run", gin.H{
		"id":        created.Data.ID,
		"variables": gin.H{"text": "hello"},
		"provider":  "mock",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("run stored: want 200, got %d: %s", w.Code, w.Body)
	}
	var resp struct {
		Rendered string `json:"rendered"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Rendered != "Translate: hello" {
		t.Fatalf("rendered = %q", resp.Rendered)
	}
}

func TestListProviders(t *testing.T) {
	r := newTestRouter(t)
	w := do(t, r, http.MethodGet, "/api/playground/providers", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("providers: want 200, got %d", w.Code)
	}
	var resp struct {
		Data []string `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) == 0 || resp.Data[0] != "mock" {
		t.Fatalf("providers = %v", resp.Data)
	}
}
