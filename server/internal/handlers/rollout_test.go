package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRolloutSplitsTraffic(t *testing.T) {
	r := newTestRouter(t)

	w := do(t, r, http.MethodPost, "/api/prompts", gin.H{
		"key": "ab.test", "content": "content v1", "version": "v1", "env": "prod",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create: %d %s", w.Code, w.Body)
	}
	var created struct {
		Data struct{ ID string } `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id := created.Data.ID

	if w := do(t, r, http.MethodPost, "/api/prompts/publish", gin.H{"id": id}); w.Code != http.StatusOK {
		t.Fatalf("publish v1: %d", w.Code)
	}
	if w := do(t, r, http.MethodPut, "/api/prompts/"+id, gin.H{
		"content": "content v2", "version": "v2", "env": "prod",
	}); w.Code != http.StatusOK {
		t.Fatalf("update v2: %d %s", w.Code, w.Body)
	}
	if w := do(t, r, http.MethodPost, "/api/prompts/publish", gin.H{"id": id}); w.Code != http.StatusOK {
		t.Fatalf("publish v2: %d", w.Code)
	}

	sdkContent := func() string {
		w := do(t, r, http.MethodGet, "/api/sdk/prompts/ab.test?env=prod", nil)
		var resp struct {
			Content string `json:"content"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		return resp.Content
	}

	// 100% of traffic to variant A (v1).
	if w := do(t, r, http.MethodPut, "/api/prompts/"+id+"/rollout", gin.H{
		"enabled": true, "variant_a": "v1", "variant_b": "v2", "weight_a": 100,
	}); w.Code != http.StatusOK {
		t.Fatalf("set rollout A: %d %s", w.Code, w.Body)
	}
	if got := sdkContent(); got != "content v1" {
		t.Fatalf("weight 100 -> A: got %q, want content v1", got)
	}

	// 100% of traffic to variant B (v2).
	if w := do(t, r, http.MethodPut, "/api/prompts/"+id+"/rollout", gin.H{
		"enabled": true, "variant_a": "v1", "variant_b": "v2", "weight_a": 0,
	}); w.Code != http.StatusOK {
		t.Fatalf("set rollout B: %d", w.Code)
	}
	if got := sdkContent(); got != "content v2" {
		t.Fatalf("weight 0 -> B: got %q, want content v2", got)
	}

	// Removing the rollout falls back to the live prompt content.
	if w := do(t, r, http.MethodDelete, "/api/prompts/"+id+"/rollout", nil); w.Code != http.StatusOK {
		t.Fatalf("delete rollout: %d", w.Code)
	}
	if got := sdkContent(); got != "content v2" {
		t.Fatalf("no rollout: got %q, want live content v2", got)
	}
}
