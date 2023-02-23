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

//  关注操作

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
	// 判断为空
	if to_user_id == "" || user_id == "" || action_type == "" {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull))
		return
	}
	// 转换
	toUserId, _ := strconv.ParseInt(to_user_id, 10, 64)
	userId := int64(user_id.(float64))
	if userId == toUserId {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg("不能自己关注自己"))
		return
	}
	var actionType bool
	// 关注
	if action_type == "1" {
		actionType = true
	} else if action_type == "2" {
		actionType = false
	} else {
		c.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg(config.RequestFail))
		return
	}
	pass, err := service.FollowUserService(userId, toUserId, actionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	}
	if !pass {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.Fail())
		return
	} else {
		c.JSON(http.StatusOK, model.BaseResponseInstance.Success())
		return
	}
}

//  关注列表

func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")
	userId, _ := strconv.ParseInt(user_id, 10, 64)
	// 取token
	loginUserId, exists := c.Get("Id")
	// 判空
	if user_id == "" || loginUserId == "" {
		c.JSON(http.StatusBadRequest, FollowListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull),
			UserList:     nil,
		})
		return
	}
	if !exists {
		// 不存在即未登录
		userInfos, err := service.FollowListService(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, FollowListResponse{
				BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
				UserList:     nil,
			})
			return
		}
		c.JSON(http.StatusOK, FollowListResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserList:     userInfos,
		})
		return
	} else {
		loginUserId := int64(loginUserId.(float64))
		userInfos, err := service.FollowListServiceWithUserId(userId, loginUserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, FollowListResponse{
				BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
				UserList:     nil,
			})
			return
		}
		c.JSON(http.StatusOK, FollowListResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserList:     userInfos,
		})
		return
	}

}

//  粉丝列表

func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")
	userId, _ := strconv.ParseInt(user_id, 10, 64)
	if user_id == "" {
		c.JSON(http.StatusBadRequest, FollowerListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull),
			UserList:     nil,
		})
		return
	}
	// 取token
	loginUserId, exists := c.Get("Id")
	if !exists {
		userInfos, err := service.FollowerListService(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, FollowerListResponse{
				BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
				UserList:     nil,
			})
			return
		}
		c.JSON(http.StatusOK, FollowerListResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserList:     userInfos,
		})
		return
	} else {
		loginUserId := int64(loginUserId.(float64))
		userInfos, err := service.FollowerListServiceWithUserId(userId, loginUserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, FollowerListResponse{
				BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
				UserList:     nil,
			})
			return
		}
		c.JSON(http.StatusOK, FollowerListResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserList:     userInfos,
		})
		return
	}

}
