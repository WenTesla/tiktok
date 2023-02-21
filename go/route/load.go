package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/go/controller"
	"tiktok/go/middle/jwt"
	"tiktok/go/model"
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
	apiRouter.POST("/publish/action/", jwt.VerifyTokenByPost, controller.VideoPublish)
	// 发布列表接口
	apiRouter.GET("/publish/list/", controller.VideoList)

	// extra apis - I
	// 用户点赞接口
	apiRouter.POST("/favorite/action/", jwt.VerifyToken, controller.LikeVideoByUserID)
	// 喜欢列表接口
	apiRouter.GET("/favorite/list/", controller.UserFavoriteList)
	// 用户评论接口
	apiRouter.POST("/comment/action/", jwt.VerifyToken, controller.CommentVideo)
	// 评论列表接口
	apiRouter.GET("/comment/list/", controller.CommentList)
	// extra apis - II
	// 关注操作
	apiRouter.POST("/relation/action/", jwt.VerifyToken, controller.FollowUser)
	// 关注列表
	apiRouter.GET("/relation/follow/list/", jwt.VerifyToken, controller.FollowList)
	// 粉丝列表
	apiRouter.GET("/relation/follower/list/", jwt.VerifyToken, controller.FollowerList)
	// 用户好友列表
	apiRouter.GET("/relation/friend/list/", jwt.VerifyToken, controller.FriendList)
	// 聊天记录
	apiRouter.GET("/message/chat/", jwt.VerifyToken, controller.MessageChat)
	// 发送消息
	apiRouter.POST("/message/action/", jwt.VerifyToken, controller.MessageAction)
	// test
	//apiRouter.GET("/test/token",jwt.SignToken)

	// other
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg("页面不存在"))
	})
}
