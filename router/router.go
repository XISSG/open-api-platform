package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xissg/open-api-platform/controller"
	"github.com/xissg/open-api-platform/middlewares"
)

func Router(router *gin.Engine) {
	//TODO:跨域

	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.ResponseMiddleware())

	userController := controller.NewUserController()
	interfaceController := controller.NewInterfaceController()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/api/user/register", userController.Register)
	router.POST("/api/user/login", userController.Login)

	authorize := router.Group("/api")
	authorize.Use(middlewares.Auth)

	interfaceInfo := authorize.Group("/interface_info")
	interfaceInfo.Use(middlewares.RateLimit)
	{
		interfaceInfo.GET("/get/:id", interfaceController.GetInterfaceDetail)
		interfaceInfo.POST("/get_list", interfaceController.GetInterfaceList)
	}

	admin := router.Group("/admin")
	admin.Use(middlewares.IsAdmin)
	{
		adminUser := admin.Group("/user")
		adminUser.POST("/get_info", userController.GetUserDetail)
		adminUser.POST("/get_list", userController.GetUserList)
		adminUser.GET("/get_info/:id", userController.GetUserDetail)
		adminUser.POST("/update_info", userController.UpdateUserInfo)
		adminUser.POST("/delete", userController.DeleteUserInfo)

		interInfo := admin.Group("/interface")
		interInfo.POST("/add_list", interfaceController.AddInterfaceInfo)
		interInfo.POST("/update", interfaceController.UpdateInterfaceInfo)
		interInfo.POST("/delete", interfaceController.DeleteInterfaceInfo)
	}

	//提供给新用户测试调用的接口，无需验证，有次数调用限制
	invokeController := controller.NewInvokeController()
	invoke := router.Group("/invoke")
	{
		invoke.POST("/", invokeController.Invoke)
	}

	//用于测试接口是否能够调用成功
	demoController := controller.NewDemo()
	invokeAPI := router.Group("/demo")
	invokeAPI.Use(middlewares.InvokeAuth)
	{
		invokeAPI.POST("/", demoController.Hello)
	}
}
