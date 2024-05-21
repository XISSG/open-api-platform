package models

import (
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/utils"
	"time"
)

type AddInfoRequest struct {
	Name           string `json:"name" binding:"max=256"`             // 名称
	Description    string `json:"description" binding:"max=256"`      // 描述
	URL            string `json:"url" binding:"required, max=512"`    // 接口地址
	RequestHeader  string `json:"request_header" binding:"max=8192"`  // 请求头
	ResponseHeader string `json:"response_header" binding:"max=8192"` // 响应头
	Status         int32  `json:"status" binding:"required, max=10"`  // 接口状态(0-关闭， 1-开启)
	Method         string `json:"method" binding:"required, max=256"` // 请求类型
	UserID         int64  `json:"user_id" binding:"required, max=64"` // 创建人
}

func AddInfoRequestToInterfaceInfo(addRequest AddInfoRequest) model.InterfaceInfo {
	id := utils.Snowflake()
	return model.InterfaceInfo{
		ID:             id,
		Name:           addRequest.Name,
		Description:    addRequest.Description,
		URL:            addRequest.URL,
		RequestHeader:  addRequest.RequestHeader,
		ResponseHeader: addRequest.ResponseHeader,
		Status:         addRequest.Status,
		Method:         addRequest.Method,
		UserID:         addRequest.UserID,
		CreateTime:     time.Now().UTC(),
	}
}

type QueryInfoRequest struct {
	Page     int `json:"page" binding:"required, max=256"`
	PageSize int `json:"page_size" binding:"required, max=64"`
}

type UpdateInfoRequest struct {
	ID             int64  `json:"id" binding:"required, max=64"`
	Name           string `json:"name" binding:"max=256"`             // 名称
	Description    string `json:"description" binding:"max=256"`      // 描述
	URL            string `json:"url" binding:"max=512"`              // 接口地址
	RequestHeader  string `json:"request_header" binding:"max=8192"`  // 请求头
	ResponseHeader string `json:"response_header" binding:"max=8192"` // 响应头
	Status         int32  `json:"status" binding:"max=10"`            // 接口状态(0-关闭， 1-开启)
	Method         string `json:"method" binding:"max=256"`           // 请求类型
	UserID         int64  `json:"user_id" binding:"max=64"`           // 创建人
}

type InfoResponse struct {
	ID             int64  `json:"id"`              // 主键
	Name           string `json:"name"`            // 名称
	Description    string `json:"description"`     // 描述
	URL            string `json:"url"`             // 接口地址
	RequestHeader  string `json:"request_header"`  // 请求头
	ResponseHeader string `json:"response_header"` // 响应头
	Status         int32  `json:"status"`          // 接口状态(0-关闭， 1-开启)
	Method         string `json:"method"`          // 请求类型
	UserID         int64  `json:"user_id"`         // 创建人
}

func InterfaceInfoToInfoResponse(interfaceInfo model.InterfaceInfo) InfoResponse {
	return InfoResponse{
		ID:             interfaceInfo.ID,
		Name:           interfaceInfo.Name,
		Description:    interfaceInfo.Description,
		URL:            interfaceInfo.URL,
		RequestHeader:  interfaceInfo.RequestHeader,
		ResponseHeader: interfaceInfo.ResponseHeader,
		Status:         interfaceInfo.Status,
		Method:         interfaceInfo.Method,
		UserID:         interfaceInfo.UserID,
	}
}
