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
	model.BaseResponse
	CommentList []model.CommentInfo `json:"comment_list"` // 评论列表
}
type CommentResponse struct {
	model.BaseResponse
	CommentInfo model.CommentInfo `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

//  用户评论 登录用户对视频进行评论

func CommentVideo(c *gin.Context) {
	// 获取登录用户
	user_id, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
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
				BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
				CommentInfo:  model.CommentInfo{},
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentResponse{
				BaseResponse: model.BaseResponseInstance.Success(),
				CommentInfo:  commentInfo,
			})
			return
		}
	case "2":
		log.Printf("在%d的视频上删除评论", video_id)
		comment_id := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(comment_id, 10, 64)
		isDelete, err := service.CancelCommentService(commentId, userId, videoId)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
			return
		} else {
			if !isDelete {
				c.JSON(http.StatusBadRequest, model.BaseResponseInstance.Fail())
				return
			} else {
				c.JSON(http.StatusOK, model.BaseResponseInstance.Success())
				return
			}
		}
	}

	println(userId)
}

//  评论列表

func CommentList(c *gin.Context) {
	// 获取视频的id
	video_id := c.Query("video_id")
	// 判空
	if video_id == "" {
		c.JSON(http.StatusOK, CommentListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull),
			CommentList:  nil,
		})
		return
	}
	// 转换
	videoId, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, CommentListResponse{
			BaseResponse: model.BaseResponseInstance.Fail(),
			CommentList:  nil,
		})
		return
	}
	// 调用服务
	commentInfos, err := service.CommentListService(videoId)
	if err != nil {
		c.JSON(http.StatusBadRequest, CommentListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			CommentList:  nil,
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		BaseResponse: model.BaseResponseInstance.Success(),
		CommentList:  commentInfos,
	})
	return
}
