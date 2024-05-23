package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"github.com/xissg/open-api-platform/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	AccessKey = "access_key"
	SecretKey = "secret_key"
)

type InvokeController struct {
	mysql *service.Mysql
	redis *service.Redis
}

func NewInvokeController() *InvokeController {
	return &InvokeController{
		mysql: service.NewMysqlService(),
		redis: service.NewRedis(),
	}
}

// 调用次数限制
func (c *InvokeController) Invoke(ctx *gin.Context) {
	var request models.InvokeRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		return
	}

	//TODO:数据校验
	data, err := c.send(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.Set("data", data)
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
