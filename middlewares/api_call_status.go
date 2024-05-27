package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/service"
	"go.uber.org/zap"
	"time"
)

func APICallStatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Log the request
		logger.SugarLogger.Info("API call",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
		)

		rdb := service.NewRedis()
		// Store the stats in Redis
		key := fmt.Sprintf("api_stats:%s:%s", method, path)
		rdb.HIncrBy(key, "count", 1)
		rdb.HIncrBy(key, "total_duration", duration.Milliseconds())
	}
}
