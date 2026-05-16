package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/732124645/promptops/server/internal/auth"
	"github.com/732124645/promptops/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var validRoles = map[string]bool{"admin": true, "editor": true, "viewer": true}

// SeedDefaultAdmin creates an "admin"/"admin" account when no users exist, so
// the Web UI always has a way in alongside the static token.
func (h *Handler) SeedDefaultAdmin() error {
	var count int64
	if err := h.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := auth.HashPassword("admin")
	if err != nil {
		return err
	}
	return h.db.Create(&models.User{
		ID:           uuid.NewString(),
		Username:     "admin",
		PasswordHash: hash,
		Role:         "admin",
		CreatedAt:    time.Now(),
	}).Error
}

// Login authenticates with either a username/password pair or the static token,
// and returns a bearer token plus the resolved role.
func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if body.Token != "" {
		if body.Token == authToken() {
			c.JSON(http.StatusOK, gin.H{
				"ok": true, "token": body.Token, "role": "admin", "username": "admin",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var u models.User
	if err := h.db.First(&u, "username = ?", body.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !auth.VerifyPassword(u.PasswordHash, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token := auth.NewToken()
	if err := h.db.Create(&models.Session{
		Token: token, UserID: u.ID, CreatedAt: time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok": true, "token": token, "role": u.Role, "username": u.Username,
	})
}

// Logout deletes the current session token (no-op for the static token).
func (h *Handler) Logout(c *gin.Context) {
	tok := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tok != "" && tok != authToken() {
		h.db.Delete(&models.Session{}, "token = ?", tok)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Me returns the current user's username and role.
func (h *Handler) Me(c *gin.Context) {
	role, _ := c.Get("role")
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"username": username, "role": role})
}

// ListUsers returns all accounts (admin only).
func (h *Handler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Order("created_at asc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser creates an account (admin only).
func (h *Handler) CreateUser(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}
	if !validRoles[body.Role] {
		body.Role = "viewer"
	}
	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u := models.User{
		ID:           uuid.NewString(),
		Username:     body.Username,
		PasswordHash: hash,
		Role:         body.Role,
		CreatedAt:    time.Now(),
	}
	if err := h.db.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

// UpdateUser changes a user's role and/or password (admin only).
func (h *Handler) UpdateUser(c *gin.Context) {
	var u models.User
	if err := h.db.First(&u, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body struct {
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if validRoles[body.Role] {
		u.Role = body.Role
	}
	if body.Password != "" {
		hash, err := auth.HashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		u.PasswordHash = hash
	}
	if err := h.db.Save(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

// DeleteUser removes a user and their sessions (admin only).
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.db.Delete(&models.Session{}, "user_id = ?", id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
