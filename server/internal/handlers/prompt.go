package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/732124645/promptops/server/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	hub *ws.Hub
}

func New(db *gorm.DB, hub *ws.Hub) *Handler {
	return &Handler{db: db, hub: hub}
}

// notify broadcasts a hot-reload event so connected SDKs refresh their cache.
func (h *Handler) notify(p models.Prompt, event string) {
	msg, _ := json.Marshal(gin.H{
		"type":    event,
		"id":      p.ID,
		"key":     p.Key,
		"env":     p.Env,
		"version": p.Version,
	})
	h.hub.Broadcast(msg)
}

// ListPrompts returns prompts filtered by free-text query, env, category and tag.
func (h *Handler) ListPrompts(c *gin.Context) {
	tx := h.db.Model(&models.Prompt{})
	if ws := c.Query("workspace"); ws != "" {
		tx = tx.Where("workspace_id = ?", ws)
	}
	if env := c.Query("env"); env != "" {
		tx = tx.Where("env = ?", env)
	}
	if category := c.Query("category"); category != "" {
		tx = tx.Where("category = ?", category)
	}
	if tag := c.Query("tag"); tag != "" {
		tx = tx.Where("tags LIKE ?", "%"+tag+"%")
	}
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + q + "%"
		tx = tx.Where("key LIKE ? OR name LIKE ? OR content LIKE ?", like, like, like)
	}
	var prompts []models.Prompt
	if err := tx.Order("updated_at desc").Find(&prompts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": prompts})
}

// GetPrompt returns a single prompt by id.
func (h *Handler) GetPrompt(c *gin.Context) {
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// CreatePrompt creates a new prompt.
func (h *Handler) CreatePrompt(c *gin.Context) {
	var body models.Prompt
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if strings.TrimSpace(body.Key) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	body.ID = uuid.NewString()
	if body.Version == "" {
		body.Version = "v1"
	}
	if body.Env == "" {
		body.Env = "dev"
	}
	now := time.Now()
	body.CreatedAt = now
	body.UpdatedAt = now
	if err := h.db.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("create", "prompt", body.ID, body.Key, body.Name)
	c.JSON(http.StatusOK, gin.H{"data": body})
}

// UpdatePrompt updates an existing prompt's editable fields.
func (h *Handler) UpdatePrompt(c *gin.Context) {
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body models.Prompt
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	p.Name = body.Name
	p.Content = body.Content
	p.Category = body.Category
	p.Tags = body.Tags
	p.Model = body.Model
	if body.Version != "" {
		p.Version = body.Version
	}
	if body.Env != "" {
		p.Env = body.Env
	}
	p.UpdatedAt = time.Now()
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.notify(p, "prompt.updated")
	h.recordAudit("update", "prompt", p.ID, p.Key, p.Name)
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// DeletePrompt removes a prompt and its version history.
func (h *Handler) DeletePrompt(c *gin.Context) {
	id := c.Param("id")
	var p models.Prompt
	h.db.First(&p, "id = ?", id)
	if err := h.db.Delete(&models.Prompt{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.db.Delete(&models.PromptVersion{}, "prompt_id = ?", id)
	h.recordAudit("delete", "prompt", id, p.Key, p.Name)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ListVersions returns the version history of a prompt, newest first.
func (h *Handler) ListVersions(c *gin.Context) {
	var versions []models.PromptVersion
	if err := h.db.Where("prompt_id = ?", c.Param("id")).
		Order("created_at desc").Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": versions})
}

// PublishPrompt snapshots the prompt's current content as an immutable version.
func (h *Handler) PublishPrompt(c *gin.Context) {
	var body struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", body.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	v := models.PromptVersion{
		ID:        uuid.NewString(),
		PromptID:  p.ID,
		Key:       p.Key,
		Version:   p.Version,
		Env:       p.Env,
		Content:   p.Content,
		CreatedAt: time.Now(),
	}
	if err := h.db.Create(&v).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.notify(p, "prompt.published")
	h.recordAudit("publish", "prompt", p.ID, p.Key, "version "+p.Version)
	c.JSON(http.StatusOK, gin.H{"data": v})
}

// RollbackPrompt restores a prompt's content from a previously published version.
func (h *Handler) RollbackPrompt(c *gin.Context) {
	var body struct {
		ID      string `json:"id"`
		Version string `json:"version"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", body.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var v models.PromptVersion
	if err := h.db.Where("prompt_id = ? AND version = ?", body.ID, body.Version).
		Order("created_at desc").First(&v).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}
	p.Content = v.Content
	p.Version = v.Version
	p.UpdatedAt = time.Now()
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.notify(p, "prompt.updated")
	h.recordAudit("rollback", "prompt", p.ID, p.Key, "to "+v.Version)
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// SDKGetPrompt is the runtime endpoint SDKs call to fetch a prompt by key + env.
func (h *Handler) SDKGetPrompt(c *gin.Context) {
	key := c.Param("key")
	env := c.DefaultQuery("env", "prod")
	var p models.Prompt
	if err := h.db.Where("key = ? AND env = ?", key, env).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	content := p.Content
	version := p.Version
	variant := ""

	// Apply a gray-release rollout: weighted pick between two published versions.
	var ro models.Rollout
	if err := h.db.Where("key = ? AND env = ? AND enabled = ?", key, env, true).
		First(&ro).Error; err == nil {
		picked, label := pickVariant(ro)
		var pv models.PromptVersion
		if err := h.db.Where("key = ? AND env = ? AND version = ?", key, env, picked).
			Order("created_at desc").First(&pv).Error; err == nil {
			content = pv.Content
			version = pv.Version
			variant = label
		}
	}

	resp := gin.H{
		"key":     p.Key,
		"version": version,
		"env":     p.Env,
		"model":   p.Model,
		"content": content,
	}
	if variant != "" {
		resp["variant"] = variant
	}
	c.JSON(http.StatusOK, resp)
}
