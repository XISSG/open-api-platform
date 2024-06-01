package models

import (
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/utils"
	"time"
)

type AddUserRequest struct {
	UserName     string `json:"user_name" validate:"required,max=256"`
	AvatarURL    string `json:"avatar_url" validate:"max=1024,omitempty"`
	UserPassword string `json:"user_password" validate:"required,max=256"`
}

func AddUserRequestToUser(addRequest AddUserRequest) model.User {
	id := utils.Snowflake()
	accessKey, _ := utils.AccessKeyGenerator()
	secretKey, _ := utils.SecretKeyGenerator()
	return model.User{
		ID:           id,
		UserName:     addRequest.UserName,
		AvatarURL:    addRequest.AvatarURL,
		UserPassword: addRequest.UserPassword,
		CreateTime:   time.Now().UTC(),
		UserRole:     constant.User,
		AccessKey:    accessKey,
		SecretKey:    secretKey,
	}
}

type LoginRequest struct {
	UserName     string `json:"user_name" validate:"required,max=256"`
	UserPassword string `json:"user_password" validate:"required,max=256"`
}
type QueryUserRequest struct {
	Page     int `json:"page" validate:"required,max=256"`
	PageSize int `json:"page_size" validate:"required,max=64"`
}

type UpdateUserRequest struct {
	ID           int64  `json:"id" validate:"required"`
	UserName     string `json:"user_name" validate:"max=256,omitempty"`
	AvatarURL    string `json:"avatar_url" validate:"max=1024,omitempty"`
	UserPassword string `json:"user_password" validate:"max=256,omitempty"`
	UserRole     string `json:"user_role" validate:"max=16,omitempty"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	UserName  string `json:"user_name"`
	AvatarURL string `json:"avatar_url"`
	UserRole  string `json:"user_role"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func UserToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		AvatarURL: user.AvatarURL,
		UserRole:  user.UserRole,
		AccessKey: user.AccessKey,
		SecretKey: user.SecretKey,
	}
}
