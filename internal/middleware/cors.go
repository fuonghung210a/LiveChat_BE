package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// CORS returns a middleware that handles CORS
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Set Access-Control-Allow-Origin
		allowedOrigins := strings.Split(config.AllowedOrigins, ",")
		if config.AllowedOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set other CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Methods", config.AllowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", config.AllowedHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}