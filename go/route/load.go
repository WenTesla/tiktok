package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoadRouter
//路由分组/*
func LoadRouter(e *gin.Engine) {
	e.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "hello",
		})
	})
}
