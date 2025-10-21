package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerConfig defines configuration for the logger middleware
type LoggerConfig struct {
	Logger          *zap.Logger
	LogRequestBody  bool
	SkipPaths       []string
	SensitiveFields []string // Fields to mask in request body (e.g., "password", "token")
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig(logger *zap.Logger) LoggerConfig {
	return LoggerConfig{
		Logger:          logger,
		LogRequestBody:  true,
		SkipPaths:       []string{"/health", "/ping"},
		SensitiveFields: []string{"password", "token", "secret"},
	}
}

// Logger creates a logging middleware with the given configuration
func Logger(logger *zap.Logger) gin.HandlerFunc {
	config := DefaultLoggerConfig(logger)
	return LoggerWithConfig(config)
}

// LoggerWithConfig creates a logging middleware with custom configuration
func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for certain paths
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Read and restore request body if needed
		var requestBody string
		if config.LogRequestBody && c.Request.Method != "GET" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the body for the actual handler
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Build log fields
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add request body if configured
		if config.LogRequestBody && requestBody != "" {
			fields = append(fields, zap.String("request_body", requestBody))
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log with appropriate level based on status code
		switch {
		case statusCode >= 500:
			config.Logger.Error("Server error", fields...)
		case statusCode >= 400:
			config.Logger.Warn("Client error", fields...)
		case statusCode >= 300:
			config.Logger.Info("Redirect", fields...)
		default:
			config.Logger.Info("Request completed", fields...)
		}
	}
}
