package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWorkspaceScopingAndCRUD(t *testing.T) {
	r := newTestRouter(t)

	// The default workspace is seeded.
	w := do(t, r, http.MethodGet, "/api/workspaces", nil)
	var list struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	hasDefault := false
	for _, ws := range list.Data {
		if ws.ID == "default" {
			hasDefault = true
		}
	}
	if !hasDefault {
		t.Fatalf("default workspace missing: %s", w.Body)
	}

	// Create a second workspace.
	w = do(t, r, http.MethodPost, "/api/workspaces", gin.H{"name": "Team B"})
	if w.Code != http.StatusOK {
		t.Fatalf("create workspace: %d %s", w.Code, w.Body)
	}
	var created struct {
		Data struct{ ID string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	teamB := created.Data.ID

	// A prompt in each workspace.
	if w := do(t, r, http.MethodPost, "/api/prompts", gin.H{
		"key": "p.default", "env": "dev", "workspace_id": "default",
	}); w.Code != http.StatusOK {
		t.Fatalf("create default prompt: %d %s", w.Code, w.Body)
	}
	if w := do(t, r, http.MethodPost, "/api/prompts", gin.H{
		"key": "p.teamb", "env": "dev", "workspace_id": teamB,
	}); w.Code != http.StatusOK {
		t.Fatalf("create teamB prompt: %d %s", w.Code, w.Body)
	}

	// Listing scoped to Team B returns only its prompt.
	w = do(t, r, http.MethodGet, "/api/prompts?workspace="+teamB, nil)
	var prompts struct {
		Data []struct{ Key string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &prompts)
	if len(prompts.Data) != 1 || prompts.Data[0].Key != "p.teamb" {
		t.Fatalf("workspace-scoped list = %s", w.Body)
	}

	// The default workspace cannot be deleted.
	if w := do(t, r, http.MethodDelete, "/api/workspaces/default", nil); w.Code != http.StatusBadRequest {
		t.Fatalf("delete default: want 400, got %d", w.Code)
	}
	// A normal workspace can be deleted.
	if w := do(t, r, http.MethodDelete, "/api/workspaces/"+teamB, nil); w.Code != http.StatusOK {
		t.Fatalf("delete Team B: %d", w.Code)
	}
}
