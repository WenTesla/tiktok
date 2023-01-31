package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktok/go/model"
	"tiktok/go/service"
	"time"
)

type VideoStreamModel struct {
	NextTime   int64         `json:"next_time,omitempty"` // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64         `json:"status_code"`         // 状态码，0-成功，其他值-失败
	StatusMsg  string        `json:"status_msg"`          // 返回状态描述
	VideoList  []model.Video `json:"video_list"`          // 视频列表
}

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
	// 传入的参数
	input_time := c.Query("latest_time")
	log.Printf("获取的参数 %s", input_time)
	var last_time time.Time
	if len(input_time) != 0 {
		// 处理传入的时间戳（这里是毫秒的）
		temp, _ := strconv.ParseInt(input_time, 10, 64)
		temp /= 1000
		last_time = time.Unix(temp, 0)
	} else {
		last_time = time.Now()
	}
	log.Printf("获取的时间戳 %v", last_time)
	//userId := 20053
	videos, err := service.VideoStreamService(last_time)
	if err != nil {
		c.JSON(http.StatusBadRequest, VideoStreamModel{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 获取发布最早的时间 作为下一条next参数 这里有问题
	nextTime, err := model.QueryNextTimeByVideoId(videos[len(videos)-1].ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, VideoStreamModel{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Printf("%v", videos)
	c.JSON(http.StatusOK, VideoStreamModel{
		NextTime:   nextTime.UnixNano() / 1e6,
		StatusCode: 0,
		StatusMsg:  "成功",
		VideoList:  videos,
	})
}

// 登录用户选择视频上传
func VideoPublish(c *gin.Context) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 获取登录用户的id
	userId, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "用户名不存在",
		})
		return
	}
	// 获取标题
	title := c.PostForm("title")
	log.Printf("%v %v %v", file, userId, title)
	c.JSON(http.StatusBadRequest, model.BaseResponse{
		StatusCode: -1,
		StatusMsg:  "失败",
	})
}

// 发布列表 用户的视频发布列表，直接列出用户所有投稿过的视频
func VideoList(c *gin.Context) {
	user_id := c.Query("user_id")
	log.Printf("%v", user_id)
	Id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "请求错误",
		})
		return
	}
	videos, err := service.VideoInfoByUserId(int(Id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "资源不存在",
		})
		return
	}
	c.JSON(http.StatusOK, VideoStreamModel{
		StatusCode: 0,
		StatusMsg:  "成功",
		VideoList:  videos,
	})
}
