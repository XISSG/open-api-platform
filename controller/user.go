package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xissg/open-api-platform/logger"
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
	GetUserByName(ctx *gin.Context)
	UpdateUserInfo(ctx *gin.Context)
	DeleteUserInfo(ctx *gin.Context)
}

type UserController struct {
	mysql     *service.Mysql
	validator *validator.Validate
}

func NewUserController(mysql *service.Mysql) UserInterface {
	return &UserController{
		mysql:     mysql,
		validator: validator.New(),
	}
}

// Register
// @Summary Create a user
// @Description Create a user
// @Tags User
// @Accept json
// @Produce json
// @Param addRequest body models.AddUserRequest true "create user request message"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /api/user/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var addRequest models.AddUserRequest
	err := ctx.ShouldBindJSON(&addRequest)
	err = c.validator.Struct(addRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	if c.mysql.IsUserExist(addRequest.UserName) {
		logger.SugarLogger.Infof("User already exist!")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "User already exist!"))
		return
	}

	addRequest.UserPassword = utils.MD5Crypt(addRequest.UserPassword)
	user := models.AddUserRequestToUser(addRequest)
	err = c.mysql.CreateUser(&user)
	if err != nil {
		logger.SugarLogger.Errorf("Create user error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Create user error"))
		return
	}

	logger.SugarLogger.Info("Register user success")
	ctx.Set("data", "Register user success")
}

// Login
// @Summary Login
// @Description Login
// @Tags User
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "login user request message"
// @Success 200 {object} middlewares.Response "tokenString"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /api/user/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	//登录成功返回一个jwt token
	var loginRequest models.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	err = c.validator.Struct(loginRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	user, err := c.mysql.GetUserByName(loginRequest.UserName)
	if err != nil {
		logger.SugarLogger.Infof("Authentication failed: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Please register first !"))
		return
	}

	loginRequest.UserPassword = utils.MD5Crypt(loginRequest.UserPassword)
	if user.UserPassword != loginRequest.UserPassword {
		logger.SugarLogger.Infof("Authentication failed: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Password is wrong!"))
		return
	}

	tokenString, _ := utils.GenerateJWT(user.UserName, user.UserRole, utils.JWTExpireTime.Unix(), []byte(utils.SecretJWTKey))
	logger.SugarLogger.Infof("Login success")
	ctx.Set("data", tokenString)
}

// GetUserByName
// @Summary Get user by name
// @Description Get user by name
// @Tags User
// @Accept json
// @Produce json
// @Param name query string true "username"
// @Success 200 {object} models.UserResponse "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/user/get_info/{name} [get]
func (c *UserController) GetUserByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" || len(name) > 256 {
		logger.SugarLogger.Infof("Authentication failed")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Invalid id"))
		return
	}

	user, err := c.mysql.GetUserByName(name)
	if err != nil {
		logger.SugarLogger.Errorf("User find error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
		return
	}

	if user == nil {
		logger.SugarLogger.Infof("User not find")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
		return
	}

	responseUser := models.UserToUserResponse(*user)
	logger.SugarLogger.Info("Get user success")
	ctx.Set("data", responseUser)
}

// GetUserList
// @Summary Get user list
// @Description Get user list
// @Tags User
// @Accept json
// @Produce json
// @Param queryRequest body models.QueryUserRequest true "query user request message"
// @Success 200 {object} []models.UserResponse "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/user/get_list [post]
func (c *UserController) GetUserList(ctx *gin.Context) {
	var queryRequest models.QueryUserRequest
	err := ctx.ShouldBindJSON(&queryRequest)
	err = c.validator.Struct(queryRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	list, err := c.mysql.GetUserList(queryRequest)
	if err != nil {
		logger.SugarLogger.Errorf("User find error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
		return
	}
	if list == nil {
		logger.SugarLogger.Infof("User not find")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
		return
	}

	responseList := make([]models.UserResponse, 0)
	for _, user := range list {
		responseUser := models.UserToUserResponse(*user)
		responseList = append(responseList, responseUser)
	}

	logger.SugarLogger.Info("Get user success")
	ctx.Set("data", responseList)
}

// UpdateUserInfo
// @Summary Update
// @Description Update user information
// @Tags User
// @Accept json
// @Produce json
// @Param updateRequest body models.UpdateUserRequest true "update user request message"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/user/update_info [post]
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	var updateRequest models.UpdateUserRequest
	err := ctx.ShouldBindJSON(&updateRequest)
	err = c.validator.Struct(updateRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	err = c.mysql.UpdateUser(updateRequest)
	if err != nil {
		logger.SugarLogger.Errorf("Update user error %v", err)
		ctx.JSON(http.StatusInternalServerError, middlewares.ErrorResponse(http.StatusInternalServerError, "Update error"))
		return
	}

	logger.SugarLogger.Info("Update user success")
	ctx.Set("data", "Update user success")
}

// DeleteUserInfo
// @Summary Get user by name
// @Description Get user by name
// @Tags User
// @Accept json
// @Produce json
// @Param id query string true "user id"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/user/delete/{id} [get]
func (c *UserController) DeleteUserInfo(ctx *gin.Context) {
	str := ctx.Param("id")
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logger.SugarLogger.Infof("Data format error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	err = c.mysql.DeleteUser(id)
	if err != nil {
		logger.SugarLogger.Errorf("Delete user error %v", err)
		ctx.JSON(http.StatusInternalServerError, middlewares.ErrorResponse(http.StatusInternalServerError, "Delete error"))
		return
	}

	logger.SugarLogger.Info("Delete user success")
	ctx.Set("data", "Delete user success")
}
