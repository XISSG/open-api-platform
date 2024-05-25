package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

const (
	AccessKey = "e51d9c31e3ac0065ebe3f1ee6261c817"
	SecretKey = "864a73225dc963b4ecd72acdab2c8e501472926e75c56b4fe8816dcd7852f0d2"
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

	//生成签名
	timestamp := time.Now().Unix()
	signature := utils.GenerateSignature(SecretKey, timestamp)

	req.Header.Set("X-Access-Key", AccessKey)
	req.Header.Set("X-Signature", signature)
	req.Header.Set("X-TimeStamp", strconv.FormatInt(timestamp, 10))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data, _ := io.ReadAll(response.Body)

	return data, nil
}
