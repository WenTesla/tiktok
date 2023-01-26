package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BasicResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 核实token，用于中间件
func VerifyToken(c *gin.Context) {
	//token := c.Request.PostFormValue("token")
	token := c.Query("token")
	fmt.Printf("%v \t \n", token)
	if len(token) == 0 {
		//错误 直接
		c.Abort()
		//返回json
		c.JSON(http.StatusBadRequest, BasicResponse{
			StatusCode: -1,
			StatusMsg:  "未携带token",
		})
		return
	}
	Id, err := ParseToken(token)
	if err != nil {
		// 解析错误
		c.Abort()
		// 返回json
		c.JSON(http.StatusBadRequest, BasicResponse{
			StatusCode: -1,
			StatusMsg:  "token错误",
		})
	} else {
		// 解析正确
		c.Set("Id", Id)
		c.Next()
	}
}
