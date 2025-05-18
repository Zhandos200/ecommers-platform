package middleware

import (
	"time"

	"api-gateway/logger"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.Log.WithFields(map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"clientIP": c.ClientIP(),
			"status":   c.Writer.Status(),
			"duration": duration,
		}).Info("📥 HTTP Request")
	}
}
