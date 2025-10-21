package main

import (
	"go_starter/internal/middleware"
	"go_starter/internal/model"
	_ "go_starter/internal/repository"
	"go_starter/internal/router"
	_ "go_starter/internal/service"
	"go_starter/internal/util"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := util.LoadENV()

	// Initialize logger
	logger := util.NewLogger()
	defer logger.Sync() // Flush any buffered log entries

	// Connect to database
	db := util.ConnectDB(cfg)
	model.AutoMigrate(db)

	// Initialize Gin without default middleware
	gin.SetMode(gin.ReleaseMode) // Set to release mode to use our custom logger
	r := gin.New()

	// Add recovery middleware (handles panics)
	r.Use(gin.Recovery())

	// Add custom logger middleware
	r.Use(middleware.Logger(logger))

	// Setup routes
	router.SetupRoutes(r, db, cfg, logger)

	// Log server start
	logger.Info("Starting server")

	// Run server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server")
	}
}
