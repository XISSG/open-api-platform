package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"net/http"
	"strconv"
)

type UserInterface interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetUserList(ctx *gin.Context)
	GetUserDetail(ctx *gin.Context)
	UpdateUserInfo(ctx *gin.Context)
	DeleteUserInfo(ctx *gin.Context)
}

type UserController struct {
	service *service.Service
}

func NewUserController() UserInterface {
	return &UserController{
		service: service.NewService(),
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var addRequest models.AddUserRequest
	err := ctx.ShouldBindJSON(&addRequest)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) Login(ctx *gin.Context) {
	//登录成功返回一个jwt token
	var loginRequest models.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) Logout(ctx *gin.Context) {

}

func (c *UserController) GetUserDetail(ctx *gin.Context) {
	str := ctx.Param("id")
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	fmt.Println(id)
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) GetUserList(ctx *gin.Context) {
	var queryRequest models.QueryUserRequest
	err := ctx.ShouldBindJSON(&queryRequest)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	var updateRequest models.UpdateUserRequest
	err := ctx.ShouldBindJSON(&updateRequest)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) DeleteUserInfo(ctx *gin.Context) {
	str := ctx.Param("id")
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	fmt.Println(id)
	ctx.JSON(http.StatusOK, nil)
}
