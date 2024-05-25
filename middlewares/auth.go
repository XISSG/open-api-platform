package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/utils"
	"net/http"
	"time"
)

type Result struct {
	user string
	role string
	exp  int64
}

// 从Authorization字段中取出jwt tokenString,进行解析，没有该字段则中断返回，若已过期则中断返回，没过期则从redis中查询该用户是否存在，若存在则中断，不存在则通过认证继续
// validateJWT 验证 JWT 并返回解析结果
func validateJWT(c *gin.Context) (*Result, bool) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		logger.SugarLogger.Infof("Invalid token string")
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Please login first"))
		c.Abort()
		return nil, false
	}

	mapClaims, err := utils.ParseJWT(tokenString, []byte(utils.SecretJWTKey))
	if err != nil {
		logger.SugarLogger.Infof("Authentication failed: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid jwt token"))
		c.Abort()
		return nil, false
	}

	res := &Result{
		user: mapClaims.Id,
		role: mapClaims.Audience,
		exp:  mapClaims.ExpiresAt,
	}

	fmt.Println(res)
	if res.exp < time.Now().Unix() {
		logger.SugarLogger.Infof("Token expired")
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Token expired"))
		c.Abort()
		return nil, false
	}

	return res, true
}

func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {
		result, valid := validateJWT(c)
		if !valid {
			c.Keys = map[string]any{}
			c.Abort()
			return
		}

		logger.SugarLogger.Infof("Valid token string")
		c.Set("user", result)
		c.Next()
	}
}

func IsAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		result, valid := validateJWT(c)
		if !valid {
			c.Abort()
			return
		}

		//判断是否是admin
		if result.role != constant.Admin {
			logger.SugarLogger.Infof("Permission denied")
			c.Abort()
			return
		}

		logger.SugarLogger.Infof("Valid token string")
		c.Set("user", result)
		c.Next()
	}
}
