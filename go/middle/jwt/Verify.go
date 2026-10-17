package jwt

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/model"
	"time"
)

// 核实token，用于中间件的验证
func VerifyToken(c *gin.Context) {
	ctx := context.Background()
	token := c.Query("token")
	fmt.Printf("%v \t \n", token)
	if len(token) == 0 {
		// 未登录开启IP限速器
		remoteIP := c.RemoteIP()
		result, _ := config.RedisClient(10).Get(ctx, remoteIP).Int()
		if result == 0 {
			// 设置
			config.RedisClient(10).Set(ctx, remoteIP, 1, time.Minute)
		} else if result <= 60 {
			// +1
			config.RedisClient(10).Incr(ctx, remoteIP)
		} else {
			c.Abort()
			c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestTooFast))
			return
		}
		c.Next()
		return
	}
	Id, err := ParseToken(token)
	if err != nil {
		// 解析错误
		c.Abort()
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
		return
	}
	// 开启用户限速器
	userId := strconv.FormatFloat(Id.(float64), 'f', 0, 64)
	result, _ := config.RedisClient(10).Get(ctx, userId).Int()
	if result == 0 {
		// 设置
		config.RedisClient(10).Set(ctx, userId, 1, time.Minute)
	} else if result <= 60 {
		// +1
		config.RedisClient(10).Incr(ctx, userId)
	} else {
		c.Abort()
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestTooFast))
		return
	}

	// token 携带 exp 声明，ParseToken 已在验签的同时完成过期校验（过期会返回 err 并走上面的解析失败分支）。
	c.Set("Id", Id)
	c.Next()
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
