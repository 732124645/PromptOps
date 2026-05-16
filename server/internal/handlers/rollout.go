package handlers

import (
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// pickVariant chooses a rollout variant by weight: WeightA percent of traffic
// goes to variant A, the rest to B.
func pickVariant(r models.Rollout) (version, label string) {
	if rand.IntN(100) < r.WeightA {
		return r.VariantA, "A"
	}
	return r.VariantB, "B"
}

// GetRollout returns the rollout configured for a prompt's key+env, or null.
func (h *Handler) GetRollout(c *gin.Context) {
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found"})
		return
	}
	var ro models.Rollout
	if err := h.db.Where("key = ? AND env = ?", p.Key, p.Env).First(&ro).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ro})
}

// SetRollout creates or updates the rollout for a prompt's key+env.
func (h *Handler) SetRollout(c *gin.Context) {
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found"})
		return
	}
	var body struct {
		Enabled  bool   `json:"enabled"`
		VariantA string `json:"variant_a"`
		VariantB string `json:"variant_b"`
		WeightA  int    `json:"weight_a"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if body.WeightA < 0 {
		body.WeightA = 0
	}
	if body.WeightA > 100 {
		body.WeightA = 100
	}

	now := time.Now()
	var ro models.Rollout
	if err := h.db.Where("key = ? AND env = ?", p.Key, p.Env).First(&ro).Error; err != nil {
		ro = models.Rollout{ID: uuid.NewString(), Key: p.Key, Env: p.Env, CreatedAt: now}
	}
	ro.Enabled = body.Enabled
	ro.VariantA = body.VariantA
	ro.VariantB = body.VariantB
	ro.WeightA = body.WeightA
	ro.UpdatedAt = now
	if err := h.db.Save(&ro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.recordAudit("rollout", "prompt", p.ID, p.Key, "gray release updated")
	c.JSON(http.StatusOK, gin.H{"data": ro})
}

// DeleteRollout removes the rollout for a prompt's key+env.
func (h *Handler) DeleteRollout(c *gin.Context) {
	var p models.Prompt
	if err := h.db.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found"})
		return
	}
	h.db.Delete(&models.Rollout{}, "key = ? AND env = ?", p.Key, p.Env)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
