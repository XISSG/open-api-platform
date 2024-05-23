package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/service"
	"github.com/xissg/open-api-platform/utils"
	"net/http"
	"time"
)

// 从Authorization字段中取出jwt tokenString,进行解析，没有该字段则中断返回，若已过期则中断返回，没过期则从redis中查询该用户是否存在，若存在则中断，不存在则通过认证继续
// validateJWT 验证 JWT 并返回解析结果
func validateJWT(c *gin.Context) (*utils.JWTResult, bool) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Please login first"))
		c.Abort()
		return nil, false
	}

	result := utils.ParseJWT(tokenString, []byte(utils.SecretJWTKey))
	if result.Err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid jwt token"))
		c.Abort()
		return nil, false
	}

	if result.Exp.Unix() < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "jwt token expired"))
		c.Abort()
		return nil, false
	}

	rdb := service.NewRedis()
	exists := rdb.Exists(result.User)
	if exists {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "User already logout"))
		c.Abort()
		return nil, false
	}

	return result, true
}

func Auth(c *gin.Context) {
	result, valid := validateJWT(c)
	if !valid {
		c.Keys = map[string]any{}
		c.Abort()
		return
	}
	c.Set("user", result)
	c.Next()
}

func IsAdmin(c *gin.Context) {
	result, valid := validateJWT(c)
	if !valid {
		c.Abort()
		return
	}

	//判断是否是admin
	if result.Role != constant.Admin {
		c.Abort()
		return
	}
	c.Set("user", result)
	c.Next()
}
