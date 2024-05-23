package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/service"
	"github.com/xissg/open-api-platform/utils"
	"strconv"
)

// API调用认证中间件
func InvokeAuth(ctx *gin.Context) {
	accessKey := ctx.GetHeader("X-Access-Key")
	signature := ctx.GetHeader("X-Signature")
	timestamp := ctx.GetHeader("X-Timestamp")
	if accessKey == "" || signature == "" || timestamp == "" {
		ctx.JSON(401, ErrorResponse(401, "Invalid access"))
		ctx.Abort()
		return
	}
	mysql := service.NewMysqlService()
	user, err := mysql.GetUserByAccessKey(accessKey)
	if err != nil {
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

	ctx.Next()
}
