package route

import (
	"github.com/gin-gonic/gin"
	"tiktok/go/controller"
)

// LoadRouter
// 路由分组/*
func LoadRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/feed/", nil)
	apiRouter.GET("/user/", nil)
	apiRouter.POST("/publish/action/", nil)
	apiRouter.GET("/publish/list/", nil)

	// extra apis - I todo

	// extra apis - II todo
}
