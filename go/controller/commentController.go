package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/service"
)

type CommentListResponse struct {
	CommentList []model.CommentInfo `json:"comment_list"` // 评论列表
	StatusCode  int64               `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string              `json:"status_msg"`   // 返回状态描述
}
type CommentResponse struct {
	CommentInfo model.CommentInfo `json:"comment"`     // 评论成功返回评论内容，不需要重新拉取整个列表
	StatusCode  int64             `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg   string            `json:"status_msg"`  // 返回状态描述
}

// CommentVideo 用户评论 登录用户对视频进行评论
func CommentVideo(c *gin.Context) {
	// 获取登录用户
	user_id, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  config.TokenIsNotMatchUserId,
		})
		return
	}
	//
	userId := int64(user_id.(float64))
	video_id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(video_id, 10, 64)
	action_type := c.Query("action_type")
	switch action_type {
	case "1":
		log.Printf("在%d的视频上发布评论", video_id)
		content := c.Query("comment_text")
		commentInfo, err := service.CreateCommentService(userId, videoId, content)
		if err != nil {
			c.JSON(http.StatusBadRequest, CommentResponse{
				CommentInfo: model.CommentInfo{},
				StatusCode:  -1,
				StatusMsg:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentResponse{
				CommentInfo: commentInfo,
				StatusCode:  0,
				StatusMsg:   config.Success,
			})
			return
		}
	case "2":
		log.Printf("在%d的视频上删除评论", video_id)
		comment_id := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(comment_id, 10, 64)
		isDelete, err := service.DeleteCommentService(commentId)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResponse{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		} else {
			if !isDelete {
				c.JSON(http.StatusBadRequest, model.BaseResponse{
					StatusCode: -1,
					StatusMsg:  config.Fail,
				})
				return
			} else {
				c.JSON(http.StatusOK, model.BaseResponse{
					StatusCode: 0,
					StatusMsg:  config.Success,
				})
				return
			}
		}
	}

	println(userId)
}

// CommentList 评论列表
func CommentList(c *gin.Context) {
	//// 获取登录用户
	//user_id, exists := c.Get("Id")
	//if !exists {
	//	c.JSON(http.StatusBadRequest, model.BaseResponse{
	//		StatusCode: -1,
	//		StatusMsg:  "失败",
	//	})
	//	return
	//}
	//userId := int64(user_id.(float64)) //userId, err := strconv.ParseInt(user_id, 10, 64)
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
	// 调用服务
	commentInfos, err := service.CommentListService(videoId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		CommentList: commentInfos,
		StatusCode:  0,
		StatusMsg:   "成功",
	})
}
