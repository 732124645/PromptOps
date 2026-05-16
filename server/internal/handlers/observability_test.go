package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestObservabilityRecordsAuditAndRuns(t *testing.T) {
	r := newTestRouter(t)

	// Creating a prompt should produce one audit entry.
	if w := do(t, r, http.MethodPost, "/api/prompts", gin.H{
		"key": "obs.test", "content": "hi {{x}}", "env": "dev",
	}); w.Code != http.StatusOK {
		t.Fatalf("create prompt: %d %s", w.Code, w.Body)
	}

	// A playground run should produce one run-log entry.
	if w := do(t, r, http.MethodPost, "/api/playground/run", gin.H{
		"content": "hello {{name}}", "variables": gin.H{"name": "Ada"}, "provider": "mock",
	}); w.Code != http.StatusOK {
		t.Fatalf("playground run: %d %s", w.Code, w.Body)
	}

	w := do(t, r, http.MethodGet, "/api/audit", nil)
	var audit struct {
		Data []struct{ Action, Resource, Key string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &audit)
	if len(audit.Data) != 1 || audit.Data[0].Action != "create" ||
		audit.Data[0].Resource != "prompt" || audit.Data[0].Key != "obs.test" {
		t.Fatalf("audit = %s", w.Body)
	}

	w = do(t, r, http.MethodGet, "/api/runs", nil)
	var runs struct {
		Data []struct {
			Source       string `json:"source"`
			Status       string `json:"status"`
			OutputTokens int    `json:"output_tokens"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &runs)
	if len(runs.Data) != 1 || runs.Data[0].Source != "playground" || runs.Data[0].Status != "ok" {
		t.Fatalf("runs = %s", w.Body)
	}
	if runs.Data[0].OutputTokens <= 0 {
		t.Fatalf("expected non-zero output tokens: %s", w.Body)
	}

	w = do(t, r, http.MethodGet, "/api/runs/stats", nil)
	var stats struct {
		Total int `json:"total"`
		OK    int `json:"ok"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &stats)
	if stats.Total != 1 || stats.OK != 1 {
		t.Fatalf("stats = %s", w.Body)
	}
}
