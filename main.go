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

	// Routes
	r.POST("/upload", h.UploadFile)
	r.GET("/download/:id", h.DownloadFile)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
