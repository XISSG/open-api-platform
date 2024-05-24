package models

import (
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/utils"
	"time"
)

type AddInfoRequest struct {
	Name           string `json:"name" validate:"max=256,omitempty"`                                      // 名称
	Description    string `json:"description" validate:"max=256,omitempty"`                               // 描述
	URL            string `json:"url" validate:"required,max=512,url"`                                    // 接口地址
	RequestHeader  string `json:"request_header" validate:"max=8192,omitempty"`                           // 请求头
	ResponseHeader string `json:"response_header" validate:"max=8192,omitempty"`                          // 响应头
	Status         int32  `json:"status" validate:"required,max=10,oneof=0 1,omitempty"`                  // 接口状态(0-关闭， 1-开启)
	Method         string `json:"method" validate:"required,max=256,oneof=GET POST PUT DELETE,omitempty"` // 请求类型
	UserID         int64  `json:"user_id" validate:"required,max=64,omitempty"`                           // 创建人
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
	Page     int `json:"page" validate:"required,max=256"`
	PageSize int `json:"page_size" validate:"required,max=64"`
}

type UpdateInfoRequest struct {
	ID             int64  `json:"id" validate:"required,max=64"`
	Name           string `json:"name" validate:"max=256,omitempty"`                             // 名称
	Description    string `json:"description" validate:"max=256,omitempty"`                      // 描述
	URL            string `json:"url" validate:"max=512,url,omitempty"`                          // 接口地址
	RequestHeader  string `json:"request_header" validate:"max=8192,omitempty"`                  // 请求头
	ResponseHeader string `json:"response_header" validate:"max=8192,omitempty"`                 // 响应头
	Status         int32  `json:"status" validate:"max=10,oneof=0 1,omitempty"`                  // 接口状态(0-关闭， 1-开启)
	Method         string `json:"method" validate:"max=256,oneof=GET POST PUT DELETE,omitempty"` // 请求类型
	UserID         int64  `json:"user_id" validate:"max=64,omitempty"`                           // 创建人
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

type InvokeRequest struct {
	URL    string `json:"url" validate:"required,url"`                 // 接口地址
	Method string `json:"method" validate:"oneof=GET POST PUT DELETE"` // 请求类型
}
