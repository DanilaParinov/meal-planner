package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"meal-planner/internal/api"
	"meal-planner/internal/config"
	"meal-planner/internal/db"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("🚀 Starting Meal Planner API (env: %s, port: %s)", cfg.ServerEnv, cfg.ServerPort)

	// Connect to the database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create the repository
	repo := db.NewRepository(database)

	// Configure Gin
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS for the frontend
	router.Use(corsMiddleware())

	// Register routes
	handler := api.NewHandler(repo)
	handler.RegisterRoutes(router)

	// Serve static frontend files
	router.Static("/ui", "./frontend")
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/admin.html", "./frontend/admin.html")
	router.StaticFile("/style.css", "./frontend/style.css")
	router.StaticFile("/app.js", "./frontend/app.js")

	// Start the server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("✓ Server listening on http://localhost%s", addr)
	log.Printf("✓ UI available at http://localhost%s/ui", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
