package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"net/http"
	"strconv"
)

type InterfaceInfo interface {
	AddInterfaceInfo(ctx *gin.Context)
	GetInterfaceList(ctx *gin.Context)
	GetInterfaceDetail(ctx *gin.Context)
	UpdateInterfaceInfo(ctx *gin.Context)
	DeleteInterfaceInfo(ctx *gin.Context)
}

type InterfaceController struct {
	service *service.Mysql
	redis   *service.Redis
}

func NewInterfaceController() InterfaceInfo {
	return &InterfaceController{
		service: service.NewMysqlService(),
		redis:   service.NewRedis(),
	}
}

func (c *InterfaceController) AddInterfaceInfo(ctx *gin.Context) {
	var addRequest models.AddInfoRequest
	err := ctx.ShouldBindJSON(&addRequest)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *InterfaceController) GetInterfaceDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.Keys = map[string]interface{}{
		"data": gin.H{
			"message": "pong",
		},
	}
	fmt.Println(id)
}

func (c *InterfaceController) GetInterfaceList(ctx *gin.Context) {
	var queryRequest models.QueryInfoRequest
	err := ctx.ShouldBindJSON(&queryRequest)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *InterfaceController) UpdateInterfaceInfo(ctx *gin.Context) {
	var updateRequest models.UpdateInfoRequest
	err := ctx.ShouldBindJSON(&updateRequest)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *InterfaceController) DeleteInterfaceInfo(ctx *gin.Context) {
	str := ctx.Param("id")
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	fmt.Println(id)
	ctx.JSON(http.StatusOK, nil)
}
