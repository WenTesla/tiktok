package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/service"
)

// followUser 关注操作
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
	switch action_type {
	// 关注
	case "1":
		{
			pass, err := service.FollowUserService(userId, toUserId, true)
			if !pass {
				c.JSON(http.StatusBadRequest, model.BaseResponse{
					StatusCode: -1,
					StatusMsg:  config.Fail,
				})
				return
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, model.BaseResponse{
					StatusCode: -1,
					StatusMsg:  err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, model.BaseResponse{
				StatusCode: 0,
				StatusMsg:  config.Success,
			})
		}
		//取消关注
	case "2":
		{
			service.FollowUserService(userId, toUserId, false)
		}
	}
}
