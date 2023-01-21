package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

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

// 用户注册 -todo
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//gin.Logger()
	fmt.Sprintln(username, password)
	//先校验参数
	Id, err := service.RegisterService(username, password)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, userRegisterResponse{
			status_code: 1,
			status_msg:  err.Error(),
			user_id:     0,
			token:       "",
		})
	}

	// 返回
	c.JSON(http.StatusOK, userRegisterResponse{
		status_code: 1,
		status_msg:  "",
		user_id:     int64(Id),
		token:       "",
	})
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
