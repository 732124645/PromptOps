package handlers

import (
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// parsePaging reads `limit` and `offset` query params with sane defaults and a
// hard cap, so a single request can never pull the whole table.
func parsePaging(c *gin.Context) (limit, offset int) {
	limit, offset = 50, 0
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		limit = v
	}
	if limit > 200 {
		limit = 200
	}
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v > 0 {
		offset = v
	}
	return limit, offset
}

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

// ListAudit returns audit-log entries, newest first, with limit/offset paging.
func (h *Handler) ListAudit(c *gin.Context) {
	limit, offset := parsePaging(c)
	var total int64
	h.db.Model(&models.AuditLog{}).Count(&total)
	var logs []models.AuditLog
	if err := h.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs, "total": total})
}

// ListRuns returns run-log entries, newest first, with limit/offset paging.
func (h *Handler) ListRuns(c *gin.Context) {
	limit, offset := parsePaging(c)
	var total int64
	h.db.Model(&models.RunLog{}).Count(&total)
	var runs []models.RunLog
	if err := h.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&runs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": runs, "total": total})
}

// ListClients returns the live hot-reload connections (SDK instances and
// browser sessions) currently held open by the WebSocket hub.
func (h *Handler) ListClients(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": h.hub.Clients()})
}

// RunStats aggregates run logs into summary metrics. The aggregation runs in
// SQL, so the whole run-log table never has to be loaded into memory.
func (h *Handler) RunStats(c *gin.Context) {
	runs := h.db.Model(&models.RunLog{})

	var agg struct {
		Total        int64
		OK           int64
		PromptTokens int64
		OutputTokens int64
		AvgLatency   float64
	}
	if err := runs.Select(
		"COUNT(*) AS total, " +
			"COALESCE(SUM(CASE WHEN status = 'ok' THEN 1 ELSE 0 END), 0) AS ok, " +
			"COALESCE(SUM(prompt_tokens), 0) AS prompt_tokens, " +
			"COALESCE(SUM(output_tokens), 0) AS output_tokens, " +
			"COALESCE(AVG(latency_ms), 0) AS avg_latency",
	).Scan(&agg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type providerRow struct {
		Provider string
		Count    int
	}
	var providerRows []providerRow
	if err := h.db.Model(&models.RunLog{}).
		Select("provider, COUNT(*) AS count").
		Group("provider").
		Scan(&providerRows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	providers := make([]gin.H, 0, len(providerRows))
	for _, p := range providerRows {
		providers = append(providers, gin.H{"provider": p.Provider, "count": p.Count})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":          agg.Total,
		"ok":             agg.OK,
		"error":          agg.Total - agg.OK,
		"prompt_tokens":  agg.PromptTokens,
		"output_tokens":  agg.OutputTokens,
		"avg_latency_ms": int(agg.AvgLatency),
		"by_provider":    providers,
	})
}
