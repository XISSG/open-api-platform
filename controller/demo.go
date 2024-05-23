package controller

import "github.com/gin-gonic/gin"

type Demo struct{}

func NewDemo() *Demo {
	return &Demo{}
}

func (d *Demo) Hello(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello world",
	})
}
