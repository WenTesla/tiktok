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
	// 注册接口
	apiRouter.POST("/user/register/", controller.Register)
	// 登录接口
	apiRouter.POST("/user/login/", controller.Login)
	// 用户信息接口
	apiRouter.GET("/user/", jwt.VerifyToken, controller.UserInfo)
	// 视频流接口
	apiRouter.GET("/feed/", controller.VideoStream)
	// 发布接口（视频上传）
	apiRouter.POST("/publish/action/", jwt.VerifyToken, controller.VideoPublish)
	// 发布列表接口
	apiRouter.GET("/publish/list/", controller.VideoList)

	// extra apis - I todo
	// 用户点赞接口
	apiRouter.POST("/favorite/action/", controller.LikeVideoByUserID)
	// 喜欢列表接口
	apiRouter.GET("/favorite/list/", controller.UserFavoriteList)
	// 用户评论接口
	apiRouter.POST("/comment/action/", jwt.VerifyToken, controller.CommentVideo)
	// 评论接口
	apiRouter.GET("/comment/list/", jwt.VerifyToken, controller.CommentList)
	// extra apis - II todo

	// test
	//apiRouter.GET("/test/token",jwt.SignToken)
}
