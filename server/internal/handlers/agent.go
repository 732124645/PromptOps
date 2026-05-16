package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/732124645/promptops/server/internal/providers"
	"github.com/732124645/promptops/server/internal/render"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListAgents returns all agents, optionally filtered by a free-text query.
func (h *Handler) ListAgents(c *gin.Context) {
	tx := h.db.Model(&models.Agent{})
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + q + "%"
		tx = tx.Where("key LIKE ? OR name LIKE ? OR description LIKE ?", like, like, like)
	}
	var agents []models.Agent
	if err := tx.Order("updated_at desc").Find(&agents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agents})
}

// GetAgent returns a single agent by id.
func (h *Handler) GetAgent(c *gin.Context) {
	var a models.Agent
	if err := h.db.First(&a, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

// CreateAgent creates a new agent.
func (h *Handler) CreateAgent(c *gin.Context) {
	var body models.Agent
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if strings.TrimSpace(body.Key) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	body.ID = uuid.NewString()
	if body.Provider == "" {
		body.Provider = "mock"
	}
	now := time.Now()
	body.CreatedAt = now
	body.UpdatedAt = now
	if err := h.db.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": body})
}

// UpdateAgent updates an existing agent.
func (h *Handler) UpdateAgent(c *gin.Context) {
	var a models.Agent
	if err := h.db.First(&a, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body models.Agent
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	a.Name = body.Name
	a.Description = body.Description
	a.Prompt = body.Prompt
	a.Model = body.Model
	if body.Provider != "" {
		a.Provider = body.Provider
	}
	a.UpdatedAt = time.Now()
	if err := h.db.Save(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

// DeleteAgent removes an agent.
func (h *Handler) DeleteAgent(c *gin.Context) {
	if err := h.db.Delete(&models.Agent{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// RunAgent renders the agent's prompt with the supplied variables and runs it
// through the agent's configured model provider.
func (h *Handler) RunAgent(c *gin.Context) {
	var a models.Agent
	if err := h.db.First(&a, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body struct {
		Variables map[string]string `json:"variables"`
		APIKey    string            `json:"api_key"`
		BaseURL   string            `json:"base_url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	rendered := render.Template(a.Prompt, body.Variables)
	prov, err := providers.Get(a.Provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := prov.Run(providers.Request{
		Model:   a.Model,
		APIKey:  body.APIKey,
		BaseURL: body.BaseURL,
		Prompt:  rendered,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "rendered": rendered})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rendered": rendered, "result": result})
}
