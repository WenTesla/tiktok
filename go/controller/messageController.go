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

type FriendListResponse struct {
	model.BaseResponse
	UserList []model.FriendUser `json:"user_list"` // 用户列表
}

type MessageListResponse struct {
	model.BaseResponse
	MessageList []model.Message `json:"message_list"`
}

// FriendList 好友列表
func FriendList(c *gin.Context) {
	user_id := c.Query("user_id")
	//userId, _ := strconv.ParseInt(user_id, 10, 64)
	log.Println(user_id)
	// 提取用户Id
	userid, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  config.TokenIsNotExist,
		})
		return
	}
	userId := int64(userid.(float64))
	//
	log.Printf("%v", userId)
	friendUsers, err := service.FriendListService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FriendListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  config.Success,
		},
		UserList: friendUsers,
	})
	return
}

// MessageChat 聊天记录 点进去才能看到
func MessageChat(c *gin.Context) {
	to_user_id := c.Query("to_user_id")
	toUserId, _ := strconv.ParseInt(to_user_id, 10, 64)
	log.Println(to_user_id)
	// 提取用户Id
	userid, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  config.TokenIsNotExist,
		})
		return
	}
	userId := int64(userid.(float64))
	messages, err := service.MessageChatService(userId, toUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, MessageListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  config.Success,
		},
		MessageList: messages,
	})
	return
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {

}
