package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/docs"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/router"
)

// @author		xissg
func main() {

	initDocs()
	logger.InitLogger()
	r := gin.New()
	router.Router(r)
	_ = r.Run(":8082")
}

func initDocs() {
	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8082"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
