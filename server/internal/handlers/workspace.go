package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const defaultWorkspaceID = "default"

// SeedDefaultWorkspace ensures a "default" workspace always exists.
func (h *Handler) SeedDefaultWorkspace() error {
	var count int64
	if err := h.db.Model(&models.Workspace{}).
		Where("id = ?", defaultWorkspaceID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return h.db.Create(&models.Workspace{
		ID:        defaultWorkspaceID,
		Name:      "Default",
		Slug:      "default",
		CreatedAt: time.Now(),
	}).Error
}

// ListWorkspaces returns all team workspaces.
func (h *Handler) ListWorkspaces(c *gin.Context) {
	var workspaces []models.Workspace
	if err := h.db.Order("created_at asc").Find(&workspaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": workspaces})
}

// CreateWorkspace creates a new workspace.
func (h *Handler) CreateWorkspace(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	w := models.Workspace{
		ID:        uuid.NewString(),
		Name:      body.Name,
		Slug:      strings.TrimSpace(body.Slug),
		CreatedAt: time.Now(),
	}
	if err := h.db.Create(&w).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("create", "workspace", w.ID, w.Slug, w.Name)
	c.JSON(http.StatusOK, gin.H{"data": w})
}

// DeleteWorkspace removes a workspace; the default workspace cannot be deleted.
func (h *Handler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == defaultWorkspaceID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the default workspace cannot be deleted"})
		return
	}
	if err := h.db.Delete(&models.Workspace{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("delete", "workspace", id, "", "")
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
