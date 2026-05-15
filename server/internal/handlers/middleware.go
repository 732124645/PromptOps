package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// authToken returns the configured admin token (PROMPTOPS_TOKEN), or a dev default.
func authToken() string {
	if v := os.Getenv("PROMPTOPS_TOKEN"); v != "" {
		return v
	}
	return "promptops-dev-token"
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

// Auth checks the Bearer token on protected routes.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		got := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if got != authToken() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
