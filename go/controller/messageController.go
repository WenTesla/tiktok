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

//  好友列表

func FriendList(c *gin.Context) {
	user_id := c.Query("user_id")
	//userId, _ := strconv.ParseInt(user_id, 10, 64)
	if user_id == "" {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull))
		return
	}
	log.Println(user_id)
	// 提取用户Id
	userid, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
		return
	}
	userId := int64(userid.(float64))
	//
	log.Printf("%v", userId)
	friendUsers, err := service.FriendListService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, FriendListResponse{
		BaseResponse: model.BaseResponseInstance.Success(),
		UserList:     friendUsers,
	})
	return
}

//  聊天记录 点进去才能看到

func MessageChat(c *gin.Context) {
	to_user_id := c.Query("to_user_id")
	if to_user_id == "" {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull))
		return
	}
	toUserId, _ := strconv.ParseInt(to_user_id, 10, 64)
	log.Println(to_user_id)
	// 提取用户Id
	userid, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
		return
	}
	userId := int64(userid.(float64))
	// 提取上次最新消息的时间
	pre_msg_time := c.Query("pre_msg_time")
	preMsgTime, err := strconv.ParseInt(pre_msg_time, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	}

	messages, err := service.MessageChatService(userId, toUserId, preMsgTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, MessageListResponse{
		BaseResponse: model.BaseResponseInstance.Success(),
		MessageList:  messages,
	})
	return
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	// 提取用户Id
	userid, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
		return
	}
	userId := int64(userid.(float64))
	// 获取toUser
	to_user_id := c.Query("to_user_id")
	toUserId, err := strconv.ParseInt(to_user_id, 10, 64)
	content := c.Query("content")
	actionType := c.Query("action_type")
	// 参数错误
	if actionType != "1" {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestFail))
		return
	}
	pass, err := service.MessageActionService(userId, toUserId, content)
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
