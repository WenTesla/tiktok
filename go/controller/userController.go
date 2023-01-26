package controller

import (
	"fmt"
	"tiktok/go/model"

	"github.com/gin-gonic/gin"

	"net/http"
	"tiktok/go/service"
)

// 要大写，不然无法传入json成功
type userRegisterResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}
type userLoginResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}
type douyinUserResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	model.UserInfo
}
type user struct {
}

// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 先校验参数
	if len(password) > 32 || len(password) <= 5 || len(username) > 32 {
		c.JSON(http.StatusBadRequest, userLoginResponse{
			StatusCode: 1,
			StatusMsg:  "用户名和密码长度错误",
			UserId:     0,
			Token:      "",
		})
		return
	}
	password = service.Encryption(password)
	//先校验合法性
	Id, err := service.RegisterService(username, password)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
			UserId:     0,
			Token:      "",
		})
	} else {
		// 颁发token
		token, _ := service.GenerateTokenByName(username)
		// 成功返回
		c.JSON(http.StatusOK, userRegisterResponse{
			StatusCode: 0,
			StatusMsg:  "成功",
			UserId:     int64(Id),
			Token:      token,
		})
	}
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 先校验参数
	if len(username) > 32 || len(password) > 32 || len(password) <= 5 {
		c.JSON(http.StatusBadRequest, userLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名和密码长度错误",
			UserId:     0,
			Token:      "",
		})
		return
	}
	password = service.Encryption(password)
	fmt.Println(username, password)
	Id, err := service.LoginService(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, userLoginResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
			UserId:     0,
			Token:      "",
		})
	} else {
		// 颁发token
		token, _ := service.GenerateTokenByName(username)
		c.JSON(http.StatusOK, userLoginResponse{
			StatusCode: 0,
			StatusMsg:  "",
			UserId:     Id,
			Token:      token,
		})
	}

}

// 用户信息 todo
func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	token := c.Query("token")
	fmt.Println(userId, token)
	// 判断token
	// mock
	c.JSON(http.StatusOK, douyinUserResponse{
		StatusCode: 0,
		StatusMsg:  "成功",
		UserInfo: model.UserInfo{
			Id:             0,
			Name:           "",
			FollowCount:    1,
			FollowerCount:  2,
			IsFollow:       false,
			TotalFavorited: 3,
			FavoriteCount:  4,
		},
	})
}
