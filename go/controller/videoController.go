package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tiktok/go/model"
	"tiktok/go/service"
	"time"
)

//type ApifoxModel struct {
//	NextTime   *int64  `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
//	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
//	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
//	VideoList  []Video `json:"video_list"`  // 视频列表
//}

// // 视频作者信息
// //
// // User
//
//	type User struct {
//		FollowCount   int64  `json:"follow_count"`  // 关注总数
//		FollowerCount int64  `json:"follower_count"`// 粉丝总数
//		ID            int64  `json:"id"`            // 用户id
//		IsFollow      bool   `json:"is_follow"`     // true-已关注，false-未关注
//		Name          string `json:"name"`          // 用户名称
//	}

type Stream struct {
	model.BaseResponse
}

// 视频流接口
func VideoStream(c *gin.Context) {
	input_time := c.Query("last_time")
	log.Printf("获取的参数 %s", input_time)
	var last_time time.Time
	if len(input_time) != 0 {

	} else {
		last_time = time.Now()
	}
	log.Printf("获取的时间戳 %v", last_time)
	//userId := 20053
	videos, err := service.VideoStreamService(last_time, 1)
	if err != nil {

	}
	c.JSON(http.StatusOK, videos)
}

// 视频发布接口
func VideoPublish(c *gin.Context) {

}

// 发布列表
func VideoList(c *gin.Context) {

}
