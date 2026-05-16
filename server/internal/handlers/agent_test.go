package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAgentLifecycleAndRun(t *testing.T) {
	r := newTestRouter(t)

	w := do(t, r, http.MethodPost, "/api/agents", gin.H{
		"key": "support.agent", "name": "Support", "provider": "mock",
		"model": "mock-1", "prompt": "You are {{persona}}. Answer: {{question}}",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create agent: want 200, got %d: %s", w.Code, w.Body)
	}
	var created struct {
		Data struct{ ID, Provider string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.Data.ID == "" || created.Data.Provider != "mock" {
		t.Fatalf("unexpected create result: %s", w.Body)
	}
	id := created.Data.ID

	w = do(t, r, http.MethodPost, "/api/agents/"+id+"/run", gin.H{
		"variables": gin.H{"persona": "a helpful bot", "question": "hi"},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("run agent: want 200, got %d: %s", w.Code, w.Body)
	}
	var run struct {
		Rendered string                    `json:"rendered"`
		Result   struct{ Provider string } `json:"result"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &run)
	if run.Rendered != "You are a helpful bot. Answer: hi" {
		t.Fatalf("rendered = %q", run.Rendered)
	}
	if run.Result.Provider != "mock" {
		t.Fatalf("provider = %q", run.Result.Provider)
	}

	w = do(t, r, http.MethodGet, "/api/agents", nil)
	var list struct {
		Data []struct{ Key string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list.Data) != 1 || list.Data[0].Key != "support.agent" {
		t.Fatalf("list = %s", w.Body)
	}

	if w := do(t, r, http.MethodDelete, "/api/agents/"+id, nil); w.Code != http.StatusOK {
		t.Fatalf("delete agent: got %d", w.Code)
	}
	if w := do(t, r, http.MethodGet, "/api/agents/"+id, nil); w.Code != http.StatusNotFound {
		t.Fatalf("get after delete: want 404, got %d", w.Code)
	}
}

func TestWorkflowLifecycleAndRun(t *testing.T) {
	r := newTestRouter(t)

	w := do(t, r, http.MethodPost, "/api/workflows", gin.H{
		"key":  "summarize.flow",
		"name": "Summarize",
		"steps": []gin.H{
			{"name": "build", "type": "render", "template": "Summarize: {{text}}"},
			{"name": "call", "type": "model", "provider": "mock", "model": "mock-1"},
			{"name": "shout", "type": "transform", "op": "upper"},
		},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create workflow: want 200, got %d: %s", w.Code, w.Body)
	}
	var created struct {
		Data struct {
			ID    string                  `json:"id"`
			Steps []struct{ Type string } `json:"steps"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.Data.ID == "" || len(created.Data.Steps) != 3 {
		t.Fatalf("unexpected create result: %s", w.Body)
	}
	id := created.Data.ID

	w = do(t, r, http.MethodPost, "/api/workflows/"+id+"/run", gin.H{
		"variables": gin.H{"text": "hello world"},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("run workflow: want 200, got %d: %s", w.Code, w.Body)
	}
	var run struct {
		Output string `json:"output"`
		Steps  []struct {
			Name   string `json:"name"`
			Output string `json:"output"`
		} `json:"steps"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &run)
	if len(run.Steps) != 3 {
		t.Fatalf("want 3 step results, got %d: %s", len(run.Steps), w.Body)
	}
	// transform "upper" runs last, so the final output must be all upper-case.
	if run.Output != strings.ToUpper(run.Output) || run.Output == "" {
		t.Fatalf("final output not upper-cased: %q", run.Output)
	}
	// The render step must have substituted the variable.
	if !strings.Contains(run.Steps[0].Output, "hello world") {
		t.Fatalf("render step output = %q", run.Steps[0].Output)
	}

	if w := do(t, r, http.MethodDelete, "/api/workflows/"+id, nil); w.Code != http.StatusOK {
		t.Fatalf("delete workflow: got %d", w.Code)
	}
}
