package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/732124645/promptops/server/internal/db"
	"github.com/732124645/promptops/server/internal/ws"
	"github.com/gin-gonic/gin"
)

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	database, err := db.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	hub := ws.NewHub()
	go hub.Run()
	h := New(database, hub)

	r := gin.New()
	r.POST("/api/login", h.Login)
	api := r.Group("/api")
	api.Use(Auth())
	{
		api.GET("/prompts", h.ListPrompts)
		api.POST("/prompts", h.CreatePrompt)
		api.POST("/prompts/publish", h.PublishPrompt)
		api.POST("/prompts/rollback", h.RollbackPrompt)
		api.GET("/prompts/:id", h.GetPrompt)
		api.PUT("/prompts/:id", h.UpdatePrompt)
		api.DELETE("/prompts/:id", h.DeletePrompt)
		api.GET("/prompts/:id/versions", h.ListVersions)
		api.GET("/sdk/prompts/:key", h.SDKGetPrompt)
		api.POST("/playground/run", h.RunPlayground)
		api.GET("/playground/providers", h.ListProviders)
		api.GET("/agents", h.ListAgents)
		api.POST("/agents", h.CreateAgent)
		api.GET("/agents/:id", h.GetAgent)
		api.PUT("/agents/:id", h.UpdateAgent)
		api.DELETE("/agents/:id", h.DeleteAgent)
		api.POST("/agents/:id/run", h.RunAgent)
		api.GET("/workflows", h.ListWorkflows)
		api.POST("/workflows", h.CreateWorkflow)
		api.GET("/workflows/:id", h.GetWorkflow)
		api.PUT("/workflows/:id", h.UpdateWorkflow)
		api.DELETE("/workflows/:id", h.DeleteWorkflow)
		api.POST("/workflows/:id/run", h.RunWorkflow)
		api.GET("/audit", h.ListAudit)
		api.GET("/runs", h.ListRuns)
		api.GET("/runs/stats", h.RunStats)
	}
	return r
}

func do(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Authorization", "Bearer promptops-dev-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthRejectsBadToken(t *testing.T) {
	r := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/prompts", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

func TestPromptLifecycle(t *testing.T) {
	r := newTestRouter(t)

	// Create.
	w := do(t, r, http.MethodPost, "/api/prompts", gin.H{
		"key": "sql.generator", "name": "SQL Generator",
		"content": "v1 content", "env": "prod", "model": "gpt-4o",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create: want 200, got %d: %s", w.Code, w.Body)
	}
	var created struct {
		Data struct{ ID, Version string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id := created.Data.ID
	if id == "" || created.Data.Version != "v1" {
		t.Fatalf("unexpected create result: %s", w.Body)
	}

	// Publish snapshots current content.
	if w := do(t, r, http.MethodPost, "/api/prompts/publish", gin.H{"id": id}); w.Code != http.StatusOK {
		t.Fatalf("publish: want 200, got %d: %s", w.Code, w.Body)
	}

	// Update content to v2.
	if w := do(t, r, http.MethodPut, "/api/prompts/"+id, gin.H{
		"name": "SQL Generator", "content": "v2 content", "version": "v2", "env": "prod",
	}); w.Code != http.StatusOK {
		t.Fatalf("update: want 200, got %d: %s", w.Code, w.Body)
	}

	// SDK fetch returns latest content.
	w = do(t, r, http.MethodGet, "/api/sdk/prompts/sql.generator?env=prod", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("sdk get: want 200, got %d: %s", w.Code, w.Body)
	}
	var sdk struct{ Content, Version string }
	_ = json.Unmarshal(w.Body.Bytes(), &sdk)
	if sdk.Content != "v2 content" {
		t.Fatalf("sdk content = %q, want v2 content", sdk.Content)
	}

	// Rollback to v1.
	if w := do(t, r, http.MethodPost, "/api/prompts/rollback", gin.H{"id": id, "version": "v1"}); w.Code != http.StatusOK {
		t.Fatalf("rollback: want 200, got %d: %s", w.Code, w.Body)
	}
	w = do(t, r, http.MethodGet, "/api/sdk/prompts/sql.generator?env=prod", nil)
	_ = json.Unmarshal(w.Body.Bytes(), &sdk)
	if sdk.Content != "v1 content" {
		t.Fatalf("after rollback content = %q, want v1 content", sdk.Content)
	}

	// Search finds it.
	w = do(t, r, http.MethodGet, "/api/prompts?q=generator", nil)
	var list struct {
		Data []struct{ Key string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list.Data) != 1 || list.Data[0].Key != "sql.generator" {
		t.Fatalf("search returned %s", w.Body)
	}

	// Delete.
	if w := do(t, r, http.MethodDelete, "/api/prompts/"+id, nil); w.Code != http.StatusOK {
		t.Fatalf("delete: want 200, got %d", w.Code)
	}
	if w := do(t, r, http.MethodGet, "/api/prompts/"+id, nil); w.Code != http.StatusNotFound {
		t.Fatalf("get after delete: want 404, got %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	r := newTestRouter(t)
	if w := do(t, r, http.MethodPost, "/api/login", gin.H{"token": "promptops-dev-token"}); w.Code != http.StatusOK {
		t.Fatalf("login ok: want 200, got %d", w.Code)
	}
	if w := do(t, r, http.MethodPost, "/api/login", gin.H{"token": "nope"}); w.Code != http.StatusUnauthorized {
		t.Fatalf("login bad: want 401, got %d", w.Code)
	}
}
