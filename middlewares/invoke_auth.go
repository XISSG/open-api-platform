package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/service"
	"github.com/xissg/open-api-platform/utils"
	"strconv"
)

// API调用认证中间件
func InvokeAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessKey := ctx.GetHeader(constant.AUTH_ACCESS_KEY)
		signature := ctx.GetHeader(constant.AUTH_ACCESS_SIGNATURE)
		timestamp := ctx.GetHeader(constant.AUTH_ACCESS_TIMESTAMP)
		if accessKey == "" || signature == "" || timestamp == "" {
			logger.SugarLogger.Infof("Invalid access key")
			ctx.JSON(401, ErrorResponse(401, "Invalid access"))
			ctx.Abort()
			return
		}
		mysql := service.NewMysql()
		user, err := mysql.GetUserByAccessKey(accessKey)
		if err != nil {
			logger.SugarLogger.Errorf("Invalid access key", err)
			ctx.JSON(401, ErrorResponse(401, "Invalid access"))
			ctx.Abort()
			return
		}

		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		checkStr := utils.GenerateSignature(user.SecretKey, timestampInt)
		if checkStr != signature {
			ctx.JSON(401, ErrorResponse(401, "Invalid access"))
			ctx.Abort()
			return
		}

		logger.SugarLogger.Infof("Success access")
		ctx.Next()
	}

}
