package controller

import (
	"fmt"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"testing"
)

func Test_send(t *testing.T) {
	invoke := NewInvokeController(service.NewMysql(), service.NewRedis())
	var request models.InvokeRequest
	request = models.InvokeRequest{
		Method: "GET",
		URL:    "https://www.baidu.com/",
	}

	res, _ := invoke.send(request)
	fmt.Printf("%v", string(res))
}
