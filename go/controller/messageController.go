package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tiktok/go/config"
	"tiktok/go/model"
)

type FriendListResponse struct {
	model.BaseResponse
	UserList []model.FriendUser `json:"user_list"` // 用户列表
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
	// test
	var FriendUsers []model.FriendUser
	FriendUsers = append(FriendUsers, model.FriendUser{
		UserInfo: model.UserInfo{
			Id:            1,
			Name:          "testuser",
			FollowCount:   1,
			FollowerCount: 2,
			IsFollow:      true,
			AvatorUrl:     "https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/20230206171653.png",
		},
		Message: "testmessage",
		MsgType: 0,
	})

	c.JSON(http.StatusOK, FriendListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  config.Success,
		},
		UserList: FriendUsers,
	})
}
