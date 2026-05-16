package handlers

import (
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// estimateTokens approximates a token count (~4 characters per token).
func estimateTokens(s string) int {
	if s == "" {
		return 0
	}
	return (utf8.RuneCountInString(s) + 3) / 4
}

// recordAudit best-effort records a mutation; failures must not break the
// triggering request.
func (h *Handler) recordAudit(action, resource, id, key, summary string) {
	h.db.Create(&models.AuditLog{
		ID:         uuid.NewString(),
		Action:     action,
		Resource:   resource,
		ResourceID: id,
		Key:        key,
		Summary:    summary,
		CreatedAt:  time.Now(),
	})
}

// recordRun best-effort records a model invocation for observability.
func (h *Handler) recordRun(source, refKey, provider, model, promptText, outputText string, latency time.Duration, runErr error) {
	entry := models.RunLog{
		ID:           uuid.NewString(),
		Source:       source,
		RefKey:       refKey,
		Provider:     provider,
		Model:        model,
		PromptTokens: estimateTokens(promptText),
		OutputTokens: estimateTokens(outputText),
		LatencyMs:    latency.Milliseconds(),
		Status:       "ok",
		CreatedAt:    time.Now(),
	}
	if runErr != nil {
		entry.Status = "error"
		entry.Error = runErr.Error()
	}
	h.db.Create(&entry)
}

// ListAudit returns the most recent audit-log entries.
func (h *Handler) ListAudit(c *gin.Context) {
	var logs []models.AuditLog
	if err := h.db.Order("created_at desc").Limit(200).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// ListRuns returns the most recent run-log entries.
func (h *Handler) ListRuns(c *gin.Context) {
	var runs []models.RunLog
	if err := h.db.Order("created_at desc").Limit(200).Find(&runs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": runs})
}

// RunStats aggregates run logs into summary metrics.
func (h *Handler) RunStats(c *gin.Context) {
	var runs []models.RunLog
	if err := h.db.Find(&runs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ok, fail, promptTokens, outputTokens int
	var totalLatency int64
	byProvider := map[string]int{}
	for _, r := range runs {
		if r.Status == "ok" {
			ok++
		} else {
			fail++
		}
		promptTokens += r.PromptTokens
		outputTokens += r.OutputTokens
		totalLatency += r.LatencyMs
		byProvider[r.Provider]++
	}

	avgLatency := 0
	if len(runs) > 0 {
		avgLatency = int(totalLatency / int64(len(runs)))
	}
	providers := make([]gin.H, 0, len(byProvider))
	for name, count := range byProvider {
		providers = append(providers, gin.H{"provider": name, "count": count})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":          len(runs),
		"ok":             ok,
		"error":          fail,
		"prompt_tokens":  promptTokens,
		"output_tokens":  outputTokens,
		"avg_latency_ms": avgLatency,
		"by_provider":    providers,
	})
}
