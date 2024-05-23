package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/middlewares"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"github.com/xissg/open-api-platform/utils"
	"net/http"
	"strconv"
)

type UserInterface interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetUserList(ctx *gin.Context)
	GetUserDetail(ctx *gin.Context)
	UpdateUserInfo(ctx *gin.Context)
	DeleteUserInfo(ctx *gin.Context)
}

type UserController struct {
	mysql *service.Mysql
	redis *service.Redis
}

func NewUserController() UserInterface {
	return &UserController{
		mysql: service.NewMysqlService(),
		redis: service.NewRedis(),
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var addRequest models.AddUserRequest
	err := ctx.ShouldBindJSON(&addRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if c.mysql.IsUserExist(addRequest.UserName) {
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "User already exist!"))
		return
	}

	addRequest.UserPassword = utils.MD5Crypt(addRequest.UserPassword)
	user := models.AddUserRequestToUser(addRequest)
	err = c.mysql.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	ctx.Set("data", "register user successfully")
}

func (c *UserController) Login(ctx *gin.Context) {
	//登录成功返回一个jwt token
	var loginRequest models.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		return
	}

	user, err := c.mysql.GetUserByName(loginRequest.UserName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Please register first !"))
		return
	}

	loginRequest.UserPassword = utils.MD5Crypt(loginRequest.UserPassword)
	if user.UserPassword != loginRequest.UserPassword {
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Password is wrong!"))
		return
	}

	tokenString, _ := utils.GenerateJWT(user.UserName, user.UserRole, utils.JWTExpireTime, []byte(utils.SecretJWTKey))
	ctx.JSON(http.StatusOK, tokenString)
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
