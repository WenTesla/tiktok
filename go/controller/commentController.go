package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktok/go/model"
	"tiktok/go/service"
)

type CommentListResponse struct {
	CommentList []service.CommentInfo `json:"comment_list"` // 评论列表
	StatusCode  int64                 `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   *string               `json:"status_msg"`   // 返回状态描述
}

// 用户评论 登录用户对视频进行评论
func CommentVideo(c *gin.Context) {
	// 获取登录用户
	userId, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "",
		})
		return
	}
	// 获取
	println(userId)
}

// 评论列表
func CommentList(c *gin.Context) {
	// 获取登录用户
	user_id, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "",
		})
		return
	}
	userId := user_id.(int64) //userId, err := strconv.ParseInt(user_id, 10, 64)
	// 获取视频的id
	video_id := c.Query("video_id")
	// 转换
	videoId, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Println(userId, videoId)
	// 调用服务
	service.CommentListService(userId, videoId)
	c.JSON(http.StatusOK, model.BaseResponse{
		StatusCode: 0,
		StatusMsg:  "成功",
	})
}
