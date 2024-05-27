package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/middlewares"
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
	mysql     *service.Mysql
	redis     *service.Redis
	validator *validator.Validate
}

func NewInterfaceController(mysql *service.Mysql, redis *service.Redis) InterfaceInfo {
	return &InterfaceController{
		mysql:     mysql,
		redis:     redis,
		validator: validator.New(),
	}
}

// AddInterfaceInfo
// @Summary Create interface information
// @Description Create interface information
// @Tags Interface information
// @Accept json
// @Produce json
// @Param addRequest body models.AddInfoRequest true "create interface request message"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/interface/add_list [post]
func (c *InterfaceController) AddInterfaceInfo(ctx *gin.Context) {
	var addRequest models.AddInfoRequest
	err := ctx.ShouldBindJSON(&addRequest)
	err = c.validator.Struct(addRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	interfaceInfo := models.AddInfoRequestToInterfaceInfo(addRequest)
	err = c.mysql.CreateInterfaceInfo(&interfaceInfo)
	if err != nil {
		logger.SugarLogger.Errorf("Craete interface info error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Create interface info error"))
		return
	}

	logger.SugarLogger.Info("Create interface info success")
	ctx.Set(constant.RESPONSE_DATA_KEY, "Create interface info success")
}

// GetInterfaceDetail
// @Summary Get interface information by id
// @Description Get interface information by id
// @Tags Interface information
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /api/interface_info/get_info/{id} [get]
func (c *InterfaceController) GetInterfaceDetail(ctx *gin.Context) {
	str := ctx.Param("id")
	if str == "" {
		logger.SugarLogger.Info("Invalid interface id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Id is empty"))
		return
	}

	id, _ := strconv.ParseInt(str, 10, 64)
	interfaceInfo, err := c.mysql.GetInterfaceInfoById(id)
	if err != nil {
		logger.SugarLogger.Errorf("Query interface info error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
		return
	}

	logger.SugarLogger.Info("Query successful")
	ctx.Set(constant.RESPONSE_DATA_KEY, interfaceInfo)
}

// GetInterfaceList
// @Summary Get interface information list
// @Description Get interface information list
// @Tags Interface information
// @Accept json
// @Produce json
// @Param queryRequest body models.QueryInfoRequest true "get interface request message"
// @Success 200 {object} []models.InfoResponse "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /api/interface_info/get_list [post]
func (c *InterfaceController) GetInterfaceList(ctx *gin.Context) {
	//使用page和pageSize生成唯一的key，将其缓存进redis
	var queryRequest models.QueryInfoRequest
	err := ctx.ShouldBindJSON(&queryRequest)
	err = c.validator.Struct(queryRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}
	var interfaceInfoList []*model.InterfaceInfo
	var results []models.InfoResponse

	cacheKey := fmt.Sprintf("interface:page=%v:size=%v", queryRequest.Page, queryRequest.PageSize)
	val, err := c.redis.Get(cacheKey)
	if err == nil && val != "" {
		err = json.Unmarshal([]byte(val), &interfaceInfoList)

	} else if err == redis.Nil || val == "" {
		interfaceInfoList, err = c.mysql.GetInterfaceInfoList(queryRequest)
		if err != nil || len(interfaceInfoList) == 0 {
			logger.SugarLogger.Errorf("Query interface info error %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
			return
		}

		data, _ := json.Marshal(interfaceInfoList)
		c.redis.Set(cacheKey, string(data))

	}

	for i := range interfaceInfoList {
		results = append(results, models.InterfaceInfoToInfoResponse(*interfaceInfoList[i]))
	}

	logger.SugarLogger.Info("Query successful")
	ctx.Set(constant.RESPONSE_DATA_KEY, results)
}

// UpdateInterfaceInfo
// @Summary Update interface information
// @Description Update interface information
// @Tags Interface information
// @Accept json
// @Produce json
// @Param updateRequest body models.UpdateInfoRequest true "update interface request message"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/interface/update [post]
func (c *InterfaceController) UpdateInterfaceInfo(ctx *gin.Context) {
	var updateRequest models.UpdateInfoRequest
	err := ctx.ShouldBindJSON(&updateRequest)
	err = c.validator.Struct(updateRequest)
	if err != nil {
		logger.SugarLogger.Infof("Data check error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Data format error"))
		return
	}

	if updateRequest.ID <= 0 || len(strconv.FormatInt(updateRequest.ID, 10)) > 256 {
		logger.SugarLogger.Infof("Authentication failed")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Invalid id"))
		return
	}

	if updateRequest.URL != "" {
		err = c.validator.Var(updateRequest.URL, "url")
		if err != nil {
			logger.SugarLogger.Infof("Data check error %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Invalid url"))
			return
		}
	}

	if updateRequest.Method != "" {
		err = c.validator.Var(updateRequest.Method, "oneof=GET POST PUT DELETE")
		if err != nil {
			logger.SugarLogger.Infof("Data check error %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Invalid method"))
			return
		}
	}

	if updateRequest.Status != 0 {
		err = c.validator.Var(updateRequest.Status, "oneof=0 1")
		if err != nil {
			logger.SugarLogger.Infof("Data check error %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Invalid status"))
			return
		}
	}

	err = c.mysql.UpdateInterfaceInfo(updateRequest)
	if err != nil {
		logger.SugarLogger.Errorf("Update interface info error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Update interface info error"))
		return
	}

	logger.SugarLogger.Info("Update interface info success")
	ctx.Set(constant.RESPONSE_DATA_KEY, "Update interface info success")
}

// DeleteInterfaceInfo
// @Summary Delete interface information by id
// @Description Delete interface information by id
// @Tags Interface information
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} middlewares.Response "ok"
// @Failure 400 {object} middlewares.Response "bad request"
// @Failure 500 {object} middlewares.Response "Internal Server Error"
// @Router /admin/interface/delete/{id} [get]
func (c *InterfaceController) DeleteInterfaceInfo(ctx *gin.Context) {
	str := ctx.Param("id")
	id, err := strconv.ParseInt(str, 10, 64)

	if len(str) > 256 || id <= 0 || err != nil {
		logger.SugarLogger.Info("Invalid interface id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Invalid interface id"))
		return
	}

	res, _ := c.mysql.GetInterfaceInfoById(id)
	if res == nil {
		logger.SugarLogger.Info("Not Found")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Not Found"))
		return
	}

	err = c.mysql.DeleteInterfaceInfo(id)
	if err != nil {
		logger.SugarLogger.Errorf("Delete interface info error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Delete interface info error"))
		return
	}

	logger.SugarLogger.Info("Delete interface info success")
	ctx.Set(constant.RESPONSE_DATA_KEY, "Delete interface info success")
}
