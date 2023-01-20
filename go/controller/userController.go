package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tiktok/go/service"
)

type userRegisterResponse struct {
	status_code int32
	status_msg  string
	user_id     int64
	token       string
}
type userLoginResponse struct {
	status_code int32
	status_msg  string
	user_id     int64
	token       string
}
type user struct {
}

// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//gin.Logger()
	fmt.Sprintln(username, password)
	//先校验参数
	pass, err := service.RegisterService(username, password)
	if !pass {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			status_code: 1,
			status_msg:  err.Error(),
			user_id:     0,
			token:       "",
		})
	}
	log.Println(username, password)
	response := userRegisterResponse{
		status_code: 0,
		status_msg:  "success",
		user_id:     0,
		token:       "1111",
	}
	c.JSON(http.StatusOK, response)
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println(username, password)

}

// 用户信息
func UserInfo(c *gin.Context) {

}
