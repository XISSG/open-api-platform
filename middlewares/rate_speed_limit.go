package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/logger"
	"golang.org/x/time/rate"
	"net/http"
)

const (
	LIMIT_PER_SEC = 2
	BUCKET_SIZE   = 2
)

// 限速逻辑，https://pkg.go.dev/golang.org/x/time/rate
func RateLimit() gin.HandlerFunc {

	return func(c *gin.Context) {
		limiter := rate.NewLimiter(LIMIT_PER_SEC, BUCKET_SIZE)
		if !limiter.Allow() {
			logger.SugarLogger.Warnf("Too many requests")
			c.JSON(http.StatusTooManyRequests, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
