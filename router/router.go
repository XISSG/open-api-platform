package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xissg/open-api-platform/controller"
)

func Router(router *gin.Engine) {
	//TODO:日志监控
	//TODO:鉴权(api调用次数监控),跨域

	//设置swagger api文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	interfaceInfo := router.Group("/api/interface")
	{
		interfaceController := controller.NewInterfaceController()
		interfaceInfo.GET("/get", interfaceController.GetInterfaceDetail)
		interfaceInfo.POST("/get_list", interfaceController.GetInterfaceList)
		interfaceInfo.POST("/add_list", interfaceController.AddInterfaceInfo)
		interfaceInfo.POST("/update", interfaceController.UpdateInterfaceInfo)
		interfaceInfo.POST("/delete", interfaceController.DeleteInterfaceInfo)
	}
	user := router.Group("/api/user")
	{
		userController := controller.NewUserController()
		user.POST("/login", userController.Login)
		user.POST("/register", userController.Register)
		user.POST("/logout", userController.Logout)
		user.POST("/get_info", userController.GetUserDetail)
		user.GET("/get_info/:id", userController.GetUserDetail)
		user.POST("/update_info", userController.UpdateUserInfo)
		user.POST("/delete", userController.DeleteUserInfo)
		user.POST("/get_list", userController.GetUserList)
	}

}
