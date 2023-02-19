package controller

import (
	"fmt"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/model"

	"github.com/gin-gonic/gin"

	"net/http"
	"tiktok/go/service"
)

// 要大写，不然无法传入json成功
type userRegisterResponse struct {
	model.BaseResponse
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type userLoginResponse struct {
	model.BaseResponse
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type douyinUserResponse struct {
	model.BaseResponse
	UserInfo *model.UserInfo `json:"user"`
}

// 用户注册

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 先判空
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg("账号密码为空"),
			UserId:       -1,
			Token:        "",
		})
		return
	}
	// 先校验参数长度
	if len(password) > 32 || len(password) <= 5 || len(username) > 32 {
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg("参数长度不正确"),
			UserId:       -1,
			Token:        "",
		})
		return
	}
	password = service.Encryption(password)
	//先校验合法性
	Id, err := service.RegisterService(username, password)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			UserId:       -1,
			Token:        "",
		})
		return
	} else {
		// 颁发token
		token, _ := service.GenerateTokenByName(username)
		// 成功返回
		c.JSON(http.StatusOK, userRegisterResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserId:       Id,
			Token:        token,
		})
		return
	}
}

// 用户登录

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 先判空
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg("账号密码为空"),
			UserId:       -1,
			Token:        "",
		})
		return
	}

	// 先校验参数
	if len(username) > 32 || len(password) > 32 || len(password) <= 5 {
		c.JSON(http.StatusBadRequest, userLoginResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg("参数长度不正确"),
			UserId:       -1,
			Token:        "",
		})
		return
	}
	password = service.Encryption(password)
	fmt.Println(username, password)
	Id, err := service.LoginService(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, userLoginResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			UserId:       -1,
			Token:        "",
		})
	} else {
		// 颁发token
		token, _ := service.GenerateTokenByName(username)
		c.JSON(http.StatusOK, userLoginResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserId:       Id,
			Token:        token,
		})
	}

}

// 用户信息

func UserInfo(c *gin.Context) {
	user_Id := c.Query("user_id")
	// 转换Id的类型
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	Id, exists := c.Get("Id")
	if exists != true {
		c.JSON(http.StatusNotFound, douyinUserResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.TokenIsNotExist),
		})
		return
	}
	if userId != int64(Id.(float64)) {
		c.JSON(http.StatusNotFound, douyinUserResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.TokenIsNotMatchUserId),
		})
		return
	}
	userInfo, err := service.UserService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, douyinUserResponse{
			BaseResponse: model.BaseResponseInstance.Success(),
			UserInfo:     &userInfo,
		})
	}

}
