package main

import (
	"log"
	"os"

	"github.com/732124645/promptops/server/internal/db"
	"github.com/732124645/promptops/server/internal/handlers"
	"github.com/732124645/promptops/server/internal/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := getenv("PROMPTOPS_DB", "data/promptops.db")
	addr := getenv("PROMPTOPS_ADDR", ":8080")

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	h := handlers.New(database, hub)
	if err := h.SeedDefaultAdmin(); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	r := gin.Default()
	r.Use(handlers.CORS())

	r.POST("/api/login", h.Login)
	r.GET("/ws", hub.ServeWS)
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Read and run routes — any authenticated user (viewer and above).
	api := r.Group("/api")
	api.Use(h.Authenticate())
	{
		api.GET("/me", h.Me)
		api.POST("/logout", h.Logout)
		api.GET("/prompts", h.ListPrompts)
		api.GET("/prompts/:id", h.GetPrompt)
		api.GET("/prompts/:id/versions", h.ListVersions)
		api.GET("/prompts/:id/rollout", h.GetRollout)
		api.GET("/sdk/prompts/:key", h.SDKGetPrompt)
		api.GET("/playground/providers", h.ListProviders)
		api.POST("/playground/run", h.RunPlayground)
		api.GET("/agents", h.ListAgents)
		api.GET("/agents/:id", h.GetAgent)
		api.POST("/agents/:id/run", h.RunAgent)
		api.GET("/workflows", h.ListWorkflows)
		api.GET("/workflows/:id", h.GetWorkflow)
		api.POST("/workflows/:id/run", h.RunWorkflow)
		api.GET("/audit", h.ListAudit)
		api.GET("/runs", h.ListRuns)
		api.GET("/runs/stats", h.RunStats)
	}

	// Mutating routes — editor and above.
	write := r.Group("/api")
	write.Use(h.Authenticate(), h.RequireRole("editor"))
	{
		write.POST("/prompts", h.CreatePrompt)
		write.PUT("/prompts/:id", h.UpdatePrompt)
		write.DELETE("/prompts/:id", h.DeletePrompt)
		write.POST("/prompts/publish", h.PublishPrompt)
		write.POST("/prompts/rollback", h.RollbackPrompt)
		write.PUT("/prompts/:id/rollout", h.SetRollout)
		write.DELETE("/prompts/:id/rollout", h.DeleteRollout)
		write.POST("/agents", h.CreateAgent)
		write.PUT("/agents/:id", h.UpdateAgent)
		write.DELETE("/agents/:id", h.DeleteAgent)
		write.POST("/workflows", h.CreateWorkflow)
		write.PUT("/workflows/:id", h.UpdateWorkflow)
		write.DELETE("/workflows/:id", h.DeleteWorkflow)
	}

	// User management — admin only.
	admin := r.Group("/api")
	admin.Use(h.Authenticate(), h.RequireRole("admin"))
	{
		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.CreateUser)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
	}

	// Serve the built Web UI if present (production single-binary mode).
	if _, err := os.Stat("web/dist/index.html"); err == nil {
		r.Static("/assets", "web/dist/assets")
		r.StaticFile("/", "web/dist/index.html")
		r.NoRoute(func(c *gin.Context) { c.File("web/dist/index.html") })
	}

	log.Printf("PromptOps server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
