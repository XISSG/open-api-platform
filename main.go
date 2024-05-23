package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/open-api-platform/router"
)

func main() {
	r := gin.New()
	router.Router(r)
	_ = r.Run(":8082")
}
