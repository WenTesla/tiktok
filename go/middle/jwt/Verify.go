package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/middle"
	"tiktok/go/model"
	"time"
)

// 核实token，用于中间件的验证
func VerifyToken(c *gin.Context) {
	token := c.Query("token")
	fmt.Printf("%v \t \n", token)
	if len(token) == 0 {
		////错误 直接
		//c.Abort()
		////返回json
		//c.JSON(http.StatusBadRequest, BasicResponse{
		//	StatusCode: -1,
		//	StatusMsg:  "未携带token",
		//})

		// 未登录开启IP限速器
		remoteIP := c.RemoteIP()
		result, _ := middle.FlowLimitRedisDbByIp.Get(remoteIP).Int()
		if result == 0 {
			// 设置
			middle.FlowLimitRedisDbByIp.Set(remoteIP, 1, time.Minute)
		} else if result <= 60 {
			// +1
			middle.FlowLimitRedisDbByIp.Incr(remoteIP)
		} else {
			c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestTooFast))
			c.Abort()
		}
		c.Next()
		return
	}
	Id, err := ParseToken(token)
	if err != nil {
		// 解析错误
		c.Abort()
		// 返回json
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
	} else {
		// 开启用户限速器
		//temp := int64(Id.(float64))
		userId := strconv.FormatFloat(Id.(float64), 'f', 0, 64)
		// 登录开启Id限速器
		result, _ := middle.FlowLimitRedisDbByIp.Get(userId).Int()
		if result == 0 {
			// 设置
			middle.FlowLimitRedisDbByIp.Set(userId, 1, time.Minute)
		} else if result <= 60 {
			// +1
			middle.FlowLimitRedisDbByIp.Incr(userId)
		} else {
			c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestTooFast))
			c.Abort()
		}

		// 解析签发时间
		tokenTime, err := ParseTokenTime(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenParseErr))
		}
		// 判断时间
		if GetDays(int64(tokenTime.(float64)), time.Now().Unix()) > 30 {
			c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenIsExpire))
		}
		// 解析正确
		//str := strconv.FormatFloat(Id, 'E', -1, 64)
		//strconv.ParseInt(str, 10, 64)
		c.Set("Id", Id)
		c.Next()
	}
}

//  通过post请求获取token

func VerifyTokenByPost(c *gin.Context) {
	token := c.PostForm("token")
	fmt.Printf("%v \t \n", token)
	if len(token) == 0 {
		//错误 直接
		c.Abort()
		//返回json
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
		return
	}
	Id, err := ParseToken(token)
	if err != nil {
		// 解析错误
		c.Abort()
		// 返回json
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenParseErr))
	} else {
		// 解析正确
		//str := strconv.FormatFloat(Id, 'E', -1, 64)
		//strconv.ParseInt(str, 10, 64)
		c.Set("Id", Id)
		c.Next()
	}
}

func GetDays(start, end int64) int {
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	sub := int(endTime.Sub(startTime).Hours())
	days := sub / 24
	if (sub % 24) > 0 {
		days = days + 1
	}
	return days
}
