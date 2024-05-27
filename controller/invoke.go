package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/middlewares"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"github.com/xissg/open-api-platform/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

type InvokeController struct {
	mysql     *service.Mysql
	redis     *service.Redis
	validator *validator.Validate
}

func NewInvokeController(mysql *service.Mysql, redis *service.Redis) *InvokeController {
	return &InvokeController{
		mysql:     mysql,
		redis:     redis,
		validator: validator.New(),
	}
}

// Invoke
// @Summary Invoke interface
// @Description Invoke interface
// @Tags Invoke
// @Accept json
// @Produce json
// @Param invokeRequest body  models.InvokeRequest true "invoke request"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /api/invoke [post]
// 调用次数限制
func (c *InvokeController) Invoke(ctx *gin.Context) {
	var request models.InvokeRequest
	err := ctx.ShouldBindJSON(&request)
	err = c.validator.Struct(request)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	data, err := c.send(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.Writer.Write(data)
}

func (c *InvokeController) send(request models.InvokeRequest) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		return nil, err
	}

	user, err := c.mysql.GetUserByName(constant.Admin)
	//生成签名
	timestamp := time.Now().Unix()
	signature := utils.GenerateSignature(user.SecretKey, timestamp)

	req.Header.Set(constant.AUTH_ACCESS_KEY, user.AccessKey)
	req.Header.Set(constant.AUTH_ACCESS_SIGNATURE, signature)
	req.Header.Set(constant.AUTH_ACCESS_TIMESTAMP, strconv.FormatInt(timestamp, 10))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data, _ := io.ReadAll(response.Body)

	return data, nil
}

// GetInvokeStatus
// @Summary Get invoke information
// @Description Get invoke information
// @Tags Invoke
// @Accept json
// @Produce json
// @Param invokeRequest body  models.GetInvokeRequest true "invoke request"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/invoke_info/status [post]
func (c *InvokeController) GetInvokeStatus(ctx *gin.Context) {
	var request models.GetInvokeRequest
	err := ctx.ShouldBindJSON(&request)
	err = c.validator.Struct(request)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	stats, err := c.getAPIStats(request.Method, request.Path)
	if err != nil {
		logger.SugarLogger.Errorf("Get API stats error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Get API stats error"))
		return
	}

	if stats == nil {
		logger.SugarLogger.Errorf("API Status Not Found")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "API Status Not Found"))
		return
	}
	ctx.Set(constant.RESPONSE_DATA_KEY, stats)
}

func (c *InvokeController) getAPIStats(method, path string) (map[string]string, error) {
	key := fmt.Sprintf("api_stats:%s:%s", method, path)
	return c.redis.HGetAll(key)
}
