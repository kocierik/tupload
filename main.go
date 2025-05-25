package main

import (
	"fmt"
	"log"

	"github.com/erik/tupload/config"
	"github.com/erik/tupload/handlers"
	"github.com/erik/tupload/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	fileStorage, err := storage.NewFileStorage(cfg.Storage.Path)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize handlers
	h := handlers.NewHandler(fileStorage, cfg.Domain)

	// Initialize router
	r := gin.Default()

	// Configure trusted proxies
	if len(cfg.Server.TrustedProxies) > 0 {
		r.SetTrustedProxies(cfg.Server.TrustedProxies)
	} else {
		// Default trusted proxies if none configured
		r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	// Serve static files
	r.StaticFile("/style.css", "./static/style.css")
	r.StaticFile("/script.js", "./static/script.js")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "PUT" {
			h.UploadFile(c)
		} else {
			c.JSON(404, gin.H{"error": "page not found"})
		}
	})

	r.GET("/download/:id", h.DownloadFile)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
