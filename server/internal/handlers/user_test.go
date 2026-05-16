package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func doAs(t *testing.T, r *gin.Engine, token, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func loginToken(t *testing.T, r *gin.Engine, body gin.H) string {
	t.Helper()
	w := doAs(t, r, "", http.MethodPost, "/api/login", body)
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: %d %s", w.Code, w.Body)
	}
	var resp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Token
}

func TestSeededAdminLogin(t *testing.T) {
	r := newTestRouter(t)
	w := doAs(t, r, "", http.MethodPost, "/api/login",
		gin.H{"username": "admin", "password": "admin"})
	if w.Code != http.StatusOK {
		t.Fatalf("admin login: want 200, got %d: %s", w.Code, w.Body)
	}
	var resp struct {
		Role  string `json:"role"`
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Role != "admin" || resp.Token == "" {
		t.Fatalf("unexpected login response: %s", w.Body)
	}
}

func TestRBACEnforcement(t *testing.T) {
	r := newTestRouter(t)
	adminTok := loginToken(t, r, gin.H{"username": "admin", "password": "admin"})

	if w := doAs(t, r, adminTok, http.MethodPost, "/api/users", gin.H{
		"username": "vince", "password": "pw", "role": "viewer",
	}); w.Code != http.StatusOK {
		t.Fatalf("create user: want 200, got %d: %s", w.Code, w.Body)
	}

	viewerTok := loginToken(t, r, gin.H{"username": "vince", "password": "pw"})

	if w := doAs(t, r, viewerTok, http.MethodGet, "/api/prompts", nil); w.Code != http.StatusOK {
		t.Fatalf("viewer read: want 200, got %d", w.Code)
	}
	if w := doAs(t, r, viewerTok, http.MethodPost, "/api/prompts", gin.H{"key": "x"}); w.Code != http.StatusForbidden {
		t.Fatalf("viewer write: want 403, got %d: %s", w.Code, w.Body)
	}
	if w := doAs(t, r, viewerTok, http.MethodGet, "/api/users", nil); w.Code != http.StatusForbidden {
		t.Fatalf("viewer users: want 403, got %d", w.Code)
	}

	w := doAs(t, r, viewerTok, http.MethodGet, "/api/me", nil)
	var me struct {
		Role string `json:"role"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &me)
	if me.Role != "viewer" {
		t.Fatalf("me role = %q", me.Role)
	}

	if w := doAs(t, r, "", http.MethodPost, "/api/login",
		gin.H{"username": "vince", "password": "bad"}); w.Code != http.StatusUnauthorized {
		t.Fatalf("bad password: want 401, got %d", w.Code)
	}
}
