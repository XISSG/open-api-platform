package models

import (
	"github.com/xissg/open-api-platform/constant"
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/utils"
	"time"
)

type AddUserRequest struct {
	UserName     string `json:"user_name" binding:"required, max=256"`
	AvatarURL    string `json:"avatar_url" binding:"max=1024"`
	UserPassword string `json:"user_password" binding:"required, max=256"`
}

func AddUserRequestToUser(addRequest AddUserRequest) model.User {
	id := utils.Snowflake()
	return model.User{
		ID:           id,
		UserName:     addRequest.UserName,
		AvatarURL:    addRequest.AvatarURL,
		UserPassword: addRequest.UserPassword,
		CreateTime:   time.Now().UTC(),
		UserRole:     constant.User,
	}
}

type LoginRequest struct {
	UserName     string `json:"user_name" binding:"required, max=256"`
	UserPassword string `json:"user_password" binding:"required, max=256"`
}
type QueryUserRequest struct {
	Page     int `json:"page" binding:"required, max=256"`
	PageSize int `json:"page_size" binding:"required, max=64"`
}

type UpdateUserRequest struct {
	ID           int64  `json:"id" binding:"required, max=64"`
	UserName     string `json:"user_name" binding:" max=256"`
	AvatarURL    string `json:"avatar_url" binding:"max=1024"`
	UserPassword string `json:"user_password" binding:"max=256"`
	UserRole     string `json:"user_role" binding:"max=16"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	UserName  string `json:"user_name"`
	AvatarURL string `json:"avatar_url"`
	UserRole  string `json:"user_role"`
}

func UserToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		AvatarURL: user.AvatarURL,
		UserRole:  user.UserRole,
	}
}
