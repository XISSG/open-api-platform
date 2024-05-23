package controller

import (
	"fmt"
	"github.com/xissg/open-api-platform/models"
	"testing"
)

func Test_send(t *testing.T) {
	invoke := NewInvokeController()
	var request models.InvokeRequest
	request = models.InvokeRequest{
		Method: "GET",
		URL:    "https://www.baidu.com/",
	}

	res, _ := invoke.send(request)
	fmt.Printf("%v", string(res))
}
