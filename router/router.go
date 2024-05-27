package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xissg/open-api-platform/apis"
	"github.com/xissg/open-api-platform/controller"
	"github.com/xissg/open-api-platform/middlewares"
	"github.com/xissg/open-api-platform/service"
)

func Router(router *gin.Engine) {

	//用于测试接口是否能够调用成功
	demoController := apis.NewDemo()
	invokeAPI := router.Group("/demo")
	invokeAPI.Use(middlewares.InvokeAuth())
	{
		invokeAPI.GET("/", demoController.Hello)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//全局中间件
	router.Use(middlewares.CORS())
	router.Use(middlewares.RateLimit())
	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.ResponseMiddleware())
	router.Use(middlewares.APICallStatsMiddleware())

	//初始化依赖
	mysql := service.NewMysql()
	redis := service.NewRedis()
	userController := controller.NewUserController(mysql)
	invokeController := controller.NewInvokeController(mysql, redis)
	interfaceController := controller.NewInterfaceController(mysql, redis)

	router.POST("/api/user/register", userController.Register)
	router.POST("/api/user/login", userController.Login)

	authorize := router.Group("/api")
	authorize.Use(middlewares.Auth())

	interfaceInfo := authorize.Group("/interface_info")
	{
		interfaceInfo.GET("/get/:id", interfaceController.GetInterfaceDetail)
		interfaceInfo.POST("/get_list", interfaceController.GetInterfaceList)
	}

	admin := router.Group("/admin")
	admin.Use(middlewares.IsAdmin())
	{
		adminUser := admin.Group("/user")
		adminUser.GET("/get_info/:name", userController.GetUserByName)
		adminUser.POST("/get_list", userController.GetUserList)
		adminUser.POST("/update_info", userController.UpdateUserInfo)
		adminUser.GET("/delete/:id", userController.DeleteUserInfo)

		interInfo := admin.Group("/interface")
		interInfo.POST("/add_list", interfaceController.AddInterfaceInfo)
		interInfo.POST("/update", interfaceController.UpdateInterfaceInfo)
		interInfo.GET("/delete/:id", interfaceController.DeleteInterfaceInfo)

		invokeInfo := admin.Group("/invoke_info")
		invokeInfo.POST("/status", invokeController.GetInvokeStatus)
	}

	invoke := router.Group("/invoke")
	{
		invoke.POST("/", invokeController.Invoke)
	}

}
