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

	r := gin.Default()
	r.Use(handlers.CORS())

	r.POST("/api/login", h.Login)
	r.GET("/ws", hub.ServeWS)
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	api := r.Group("/api")
	api.Use(handlers.Auth())
	{
		api.GET("/prompts", h.ListPrompts)
		api.POST("/prompts", h.CreatePrompt)
		api.POST("/prompts/publish", h.PublishPrompt)
		api.POST("/prompts/rollback", h.RollbackPrompt)
		api.GET("/prompts/:id", h.GetPrompt)
		api.PUT("/prompts/:id", h.UpdatePrompt)
		api.DELETE("/prompts/:id", h.DeletePrompt)
		api.GET("/prompts/:id/versions", h.ListVersions)
		api.GET("/sdk/prompts/:key", h.SDKGetPrompt)
		api.POST("/playground/run", h.RunPlayground)
		api.GET("/playground/providers", h.ListProviders)

		api.GET("/agents", h.ListAgents)
		api.POST("/agents", h.CreateAgent)
		api.GET("/agents/:id", h.GetAgent)
		api.PUT("/agents/:id", h.UpdateAgent)
		api.DELETE("/agents/:id", h.DeleteAgent)
		api.POST("/agents/:id/run", h.RunAgent)

		api.GET("/workflows", h.ListWorkflows)
		api.POST("/workflows", h.CreateWorkflow)
		api.GET("/workflows/:id", h.GetWorkflow)
		api.PUT("/workflows/:id", h.UpdateWorkflow)
		api.DELETE("/workflows/:id", h.DeleteWorkflow)
		api.POST("/workflows/:id/run", h.RunWorkflow)
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
