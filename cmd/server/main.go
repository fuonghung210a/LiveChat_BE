package main

import (
	"go_starter/internal/middleware"
	"go_starter/internal/model"
	_ "go_starter/internal/repository"
	"go_starter/internal/router"
	_ "go_starter/internal/service"
	"go_starter/internal/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := util.LoadENV()

	// Initialize logger
	logger := util.NewLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger) // Flush any buffered log entries

	// Connect to database
	db := util.ConnectDB(cfg)
	model.AutoMigrate(db)

	// Initialize Gin without default middleware
	gin.SetMode(gin.DebugMode) // Set to release mode to use our custom logger
	r := gin.New()

	// Add recovery middleware (handles panics)
	r.Use(gin.Recovery())

	// Add custom logger middleware
	r.Use(middleware.Logger(logger))

	// Add CORS middleware
	r.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		AllowedMethods: cfg.CORS.AllowedMethods,
		AllowedHeaders: cfg.CORS.AllowedHeaders,
	}))

	// Setup routes
	router.SetupRoutes(r, db, cfg, logger)

	// Log server start
	logger.Info("Starting server")

	// Run server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server")
	}
}
