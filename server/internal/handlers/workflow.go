package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/732124645/promptops/server/internal/workflow"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type workflowBody struct {
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	WorkspaceID string          `json:"workspace_id"`
	Steps       []workflow.Step `json:"steps"`
}

// workflowResponse converts a stored workflow into an API payload, decoding
// the JSON-encoded steps into a structured array.
func workflowResponse(w models.Workflow) gin.H {
	steps := []workflow.Step{}
	if w.Steps != "" {
		_ = json.Unmarshal([]byte(w.Steps), &steps)
	}
	return gin.H{
		"id":           w.ID,
		"workspace_id": w.WorkspaceID,
		"key":          w.Key,
		"name":         w.Name,
		"description":  w.Description,
		"steps":        steps,
		"created_at":   w.CreatedAt,
		"updated_at":   w.UpdatedAt,
	}
}

// ListWorkflows returns all workflows, optionally filtered by a free-text query.
func (h *Handler) ListWorkflows(c *gin.Context) {
	tx := h.db.Model(&models.Workflow{})
	if ws := c.Query("workspace"); ws != "" {
		tx = tx.Where("workspace_id = ?", ws)
	}
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + q + "%"
		tx = tx.Where("key LIKE ? OR name LIKE ? OR description LIKE ?", like, like, like)
	}
	var workflows []models.Workflow
	if err := tx.Order("updated_at desc").Find(&workflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]gin.H, 0, len(workflows))
	for _, w := range workflows {
		out = append(out, workflowResponse(w))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// GetWorkflow returns a single workflow by id.
func (h *Handler) GetWorkflow(c *gin.Context) {
	var w models.Workflow
	if err := h.db.First(&w, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": workflowResponse(w)})
}

// CreateWorkflow creates a new workflow.
func (h *Handler) CreateWorkflow(c *gin.Context) {
	var body workflowBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if strings.TrimSpace(body.Key) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	steps, _ := json.Marshal(body.Steps)
	now := time.Now()
	w := models.Workflow{
		ID:          uuid.NewString(),
		WorkspaceID: body.WorkspaceID,
		Key:         body.Key,
		Name:        body.Name,
		Description: body.Description,
		Steps:       string(steps),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.db.Create(&w).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("create", "workflow", w.ID, w.Key, w.Name)
	c.JSON(http.StatusOK, gin.H{"data": workflowResponse(w)})
}

// UpdateWorkflow updates an existing workflow.
func (h *Handler) UpdateWorkflow(c *gin.Context) {
	var w models.Workflow
	if err := h.db.First(&w, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body workflowBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	steps, _ := json.Marshal(body.Steps)
	w.Name = body.Name
	w.Description = body.Description
	w.Steps = string(steps)
	w.UpdatedAt = time.Now()
	if err := h.db.Save(&w).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("update", "workflow", w.ID, w.Key, w.Name)
	c.JSON(http.StatusOK, gin.H{"data": workflowResponse(w)})
}

// DeleteWorkflow removes a workflow.
func (h *Handler) DeleteWorkflow(c *gin.Context) {
	id := c.Param("id")
	var existing models.Workflow
	h.db.First(&existing, "id = ?", id)
	if err := h.db.Delete(&models.Workflow{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("delete", "workflow", id, existing.Key, existing.Name)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// RunWorkflow executes a workflow's steps and returns the per-step trace.
func (h *Handler) RunWorkflow(c *gin.Context) {
	var w models.Workflow
	if err := h.db.First(&w, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	steps := []workflow.Step{}
	if w.Steps != "" {
		_ = json.Unmarshal([]byte(w.Steps), &steps)
	}
	var body struct {
		Variables map[string]string `json:"variables"`
		APIKey    string            `json:"api_key"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	start := time.Now()
	result, err := workflow.Run(steps, body.Variables, body.APIKey)
	if err != nil {
		h.recordRun("workflow", w.Key, "workflow", w.Key, "", "", time.Since(start), err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	h.recordRun("workflow", w.Key, "workflow", w.Key, "", result.Output, time.Since(start), nil)
	c.JSON(http.StatusOK, gin.H{"output": result.Output, "steps": result.Steps})
}
