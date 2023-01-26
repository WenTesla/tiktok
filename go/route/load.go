package route

import (
	"github.com/gin-gonic/gin"
	"tiktok/go/controller"
	"tiktok/go/middle/jwt"
)

// LoadRouter
// 路由分组/*
func LoadRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/user/", jwt.VerifyToken, controller.UserInfo)
	apiRouter.GET("/feed/", nil)
	apiRouter.POST("/publish/action/", nil)
	apiRouter.GET("/publish/list/", nil)

	// extra apis - I todo

	// extra apis - II todo

	// test
	//apiRouter.GET("/test/token",jwt.SignToken)
}
