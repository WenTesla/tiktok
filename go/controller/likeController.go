package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/service"
)

type userFavoriteListResponse struct {
	model.BaseResponse
	VideoList []model.Video `json:"video_list"`
}

// 点赞行为:  1-点赞，2-取消点赞

func LikeVideoByUserID(c *gin.Context) {
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	// 判断是否为空
	if video_id == "" || action_type == "" {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull))
		return
	}
	// 转类型
	id, _ := strconv.ParseInt(video_id, 10, 64)
	Type, _ := strconv.ParseInt(action_type, 10, 64)
	// 用于同步数据库
	switch Type {
	case 1:
		//设置为0
		Type = 0
	case 2:
		Type = 1
	default:
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestFail))
		return
	}
	// 提取用户Id
	user_id, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
		return
	}
	userId := int64(user_id.(float64))
	// 点赞Type为1 取消为2
	flag, err := service.LikeVideoByUserIDService(userId, id, Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	}
	if flag {
		c.JSON(http.StatusOK, model.BaseResponseInstance.Success())
		return
	}
}

// 用户点赞列表

func UserFavoriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 判空
	if user_id == "" {
		c.JSON(http.StatusBadRequest, userFavoriteListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull),
			VideoList:    []model.Video{},
		})
		return
	}
	userId, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, userFavoriteListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			VideoList:    []model.Video{},
		})
		return
	}
	userFavoriteList, err := service.UserFavoriteListService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, userFavoriteListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			VideoList:    []model.Video{},
		})
		return
	} else {
		c.JSON(http.StatusOK, userFavoriteListResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			VideoList:    userFavoriteList,
		})
		return
	}
}
