package handlers

import (
	"net/http"
	"time"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/732124645/promptops/server/internal/providers"
	"github.com/732124645/promptops/server/internal/render"
	"github.com/gin-gonic/gin"
)

// RunPlayground renders a prompt with the supplied variables and runs it
// through the selected model provider.
func (h *Handler) RunPlayground(c *gin.Context) {
	var body struct {
		ID        string            `json:"id"`
		Content   string            `json:"content"`
		Variables map[string]string `json:"variables"`
		Provider  string            `json:"provider"`
		Model     string            `json:"model"`
		APIKey    string            `json:"api_key"`
		BaseURL   string            `json:"base_url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	content := body.Content
	if body.ID != "" {
		var p models.Prompt
		if err := h.db.First(&p, "id = ?", body.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found"})
			return
		}
		content = p.Content
		if body.Model == "" {
			body.Model = p.Model
		}
	}
	rendered := render.Template(content, body.Variables)

	prov, err := providers.Get(body.Provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start := time.Now()
	result, err := prov.Run(providers.Request{
		Model:   body.Model,
		APIKey:  body.APIKey,
		BaseURL: body.BaseURL,
		Prompt:  rendered,
	})
	if err != nil {
		h.recordRun("playground", "", body.Provider, body.Model, rendered, "", time.Since(start), err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "rendered": rendered})
		return
	}
	h.recordRun("playground", "", result.Provider, result.Model, rendered, result.Output, time.Since(start), nil)
	c.JSON(http.StatusOK, gin.H{"rendered": rendered, "result": result})
}

// ListProviders returns the available model providers.
func (h *Handler) ListProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": providers.Names()})
}
