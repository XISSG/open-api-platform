package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"
)

const (
	RateLimitDuration = time.Minute // 限制时间窗口
	MaxRequests       = 60          // 每个时间窗口内的最大请求数
)

// TODO:限流逻辑待优化使用，令牌桶？
func RateLimit(ctx *gin.Context) {}
