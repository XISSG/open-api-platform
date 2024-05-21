package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/router"
)

func main() {
	r := gin.New()
	router.Router(r)
	r.Run(":8081")
}
