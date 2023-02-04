package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/service"
)

type FollowListResponse struct {
	model.BaseResponse
	UserList []model.UserInfo `json:"user_list"` // 用户信息列表
}

type FollowerListResponse struct {
	model.BaseResponse
	UserList []model.UserInfo `json:"user_list"` // 用户信息列表
}

// FollowUser 关注操作
func FollowUser(c *gin.Context) {
	// 对方用户的id
	to_user_id := c.Query("to_user_id")
	// 取token
	user_id, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  config.TokenIsNotExist,
		})
		return
	}
	// 取类型
	action_type := c.Query("action_type")
	// 转换
	toUserId, _ := strconv.ParseInt(to_user_id, 10, 64)
	userId := int64(user_id.(float64))
	if userId == toUserId {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "不能自己关注自己",
		})
		return
	}
	var actionType bool
	// 关注
	if action_type == "1" {
		actionType = true
	} else if action_type == "2" {
		actionType = false
	} else {
		c.JSON(http.StatusNotFound, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  config.RequestFail,
		})
		return
	}
	pass, err := service.FollowUserService(userId, toUserId, actionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if !pass {
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

// FollowList 关注列表
func FollowList(c *gin.Context) {
	//// 取token
	//user_id, exists := c.Get("Id")
	//if !exists {
	//	c.JSON(http.StatusBadRequest, model.BaseResponse.Fail)
	//	return
	//}
	//_ = int64(user_id.(float64))
	user_id := c.Query("user_id")
	userId, _ := strconv.ParseInt(user_id, 10, 64)
	userInfos, err := service.FollowListService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, FollowListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: -1,
				StatusMsg:  config.Fail,
			},
			UserList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, FollowListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  config.Success,
		},
		UserList: userInfos,
	})
	return
}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")
	userId, _ := strconv.ParseInt(user_id, 10, 64)
	userInfos, err := service.FollowerListService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  config.Fail,
		})
		return
	}
	c.JSON(http.StatusOK, FollowerListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  config.Success,
		},
		UserList: userInfos,
	})
	return
}
