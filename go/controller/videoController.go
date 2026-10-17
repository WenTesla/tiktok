package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktok/go/config"
	"tiktok/go/middle/jwt"
	"tiktok/go/model"
	"tiktok/go/service"
	"time"
)

type VideoStreamModel struct {
	model.BaseResponse
	NextTime  int64         `json:"next_time,omitempty"` // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList []model.Video `json:"video_list"`          // 视频列表
}

type VideoPublishListResponse struct {
	model.BaseResponse
	VideoList []model.Video `json:"video_list"` // 用户发布的视频列表
}

type Stream struct {
	model.BaseResponse
}

// maxUploadBytes 上传视频大小上限（约 100MB）
const maxUploadBytes = 100 << 20

// 定义变量

//  视频流接口

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

	token := c.Query("token")
	userId := int64(-1)
	if len(token) != 0 {
		parseToken, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, VideoStreamModel{
				BaseResponse: model.BaseResponseInstance.FailMsg(config.TokenParseErr),
				VideoList:    []model.Video{},
			})
			return
		}
		userId = int64(parseToken.(float64))
	}

	videos, err := service.VideoStreamService(last_time, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, VideoStreamModel{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			VideoList:    []model.Video{},
		})
		return
	}

	// 空列表直接返回，避免 videos[len(videos)-1] 越界 panic（空库 / 没有更多视频）
	if len(videos) == 0 {
		c.JSON(http.StatusOK, VideoStreamModel{
			BaseResponse: model.BaseResponseInstance.Success(),
			NextTime:     0,
			VideoList:    []model.Video{},
		})
		return
	}

	nextTime, err := model.QueryNextTimeByVideoId(videos[len(videos)-1].ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, VideoStreamModel{
			BaseResponse: model.BaseResponseInstance.FailMsg(err.Error()),
			VideoList:    []model.Video{},
		})
		return
	}
	c.JSON(http.StatusOK, VideoStreamModel{
		NextTime:     nextTime.UnixNano() / 1e6,
		BaseResponse: model.BaseResponseInstance.Success(),
		VideoList:    videos,
	})
}

//  登录用户选择视频上传

func VideoPublish(c *gin.Context) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestFail))
		return
	}
	// 参数判断空
	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(config.RequestParameterIsNull))
		return
	}
	// 文件大小上限（约 100MB）
	if file.Size > maxUploadBytes {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg("文件过大，仅支持 100MB 以内的视频"))
		return
	}
	// 获取登录用户的id
	user_id, exists := c.Get("Id")
	if !exists {
		c.JSON(http.StatusNotFound, model.BaseResponseInstance.FailMsg(config.RequestFail))
		return
	}
	// 转换
	//user_id.(float64)
	// 获取标题
	title := c.PostForm("title")
	// 上传视频
	err = service.PublishVideoService(file, int64(user_id.(float64)), title)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponseInstance.FailMsg(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, model.BaseResponseInstance.Success())
		return
	}

}

//  发布列表 用户的视频发布列表，直接列出用户所有投稿过的视频

func VideoList(c *gin.Context) {
	user_id := c.Query("user_id")
	Id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, VideoPublishListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.RequestFail),
			VideoList:    nil,
		})
		return
	}
	videos, err := service.VideoInfoByUserId(int(Id))
	if err != nil {
		c.JSON(http.StatusNotFound, VideoPublishListResponse{
			BaseResponse: model.BaseResponseInstance.FailMsg(config.DatabaseError),
			VideoList:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, VideoPublishListResponse{
		BaseResponse: model.BaseResponseInstance.Success(),
		VideoList:    videos,
	})
}
