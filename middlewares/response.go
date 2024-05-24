package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/logger"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "Success",
		Data:    data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.SugarLogger.Panic(err)
				c.JSON(http.StatusInternalServerError, ErrorResponse(500, "Internal Server Error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理请求后的响应
		if c.Writer.Status() != http.StatusOK {
			// 如果状态码不是 200，返回错误响应
			logger.SugarLogger.Errorf("Error request %v", c.Writer.Status())
			c.JSON(c.Writer.Status(), ErrorResponse(c.Writer.Status(), c.Errors.ByType(gin.ErrorTypePrivate).String()))
		} else {
			// 返回成功响应
			logger.SugarLogger.Infof("Success request %v", c.Writer.Status())
			c.JSON(http.StatusOK, SuccessResponse(c.Keys["data"]))
		}
	}
}
