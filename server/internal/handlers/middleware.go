package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/732124645/promptops/server/internal/models"
	"github.com/gin-gonic/gin"
)

// authToken returns the static admin token (PROMPTOPS_TOKEN), or a dev default.
// It always works as an "admin" escape hatch alongside user accounts.
func authToken() string {
	if v := os.Getenv("PROMPTOPS_TOKEN"); v != "" {
		return v
	}
	return "promptops-dev-token"
}

var roleRank = map[string]int{"viewer": 1, "editor": 2, "admin": 3}

func roleAtLeast(role, min string) bool {
	return roleRank[role] >= roleRank[min]
}

// CORS allows the Vite dev server (and SDKs) to call the API from any origin.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Authenticate accepts either the static admin token or a session token and
// stores the resolved role (and user, when available) on the request context.
func (h *Handler) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tok == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if tok == authToken() {
			c.Set("role", "admin")
			c.Set("username", "admin")
			c.Next()
			return
		}
		var s models.Session
		if err := h.db.First(&s, "token = ?", tok).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var u models.User
		if err := h.db.First(&u, "id = ?", s.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("role", u.Role)
		c.Set("username", u.Username)
		c.Set("user_id", u.ID)
		c.Next()
	}
}

// RequireRole aborts the request unless the authenticated role is at least min.
func (h *Handler) RequireRole(min string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		current, _ := role.(string)
		if !roleAtLeast(current, min) {
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "forbidden: requires " + min + " role"})
			return
		}
		c.Next()
	}
}
