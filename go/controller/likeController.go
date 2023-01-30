package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/go/model"
	"tiktok/go/service"
)

type userFavoriteListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	VideoList  []model.Video `json:"video_list"`
}

// 点赞行为:  1-点赞，2-取消点赞
func LikeVideoByUserID(c *gin.Context) {
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	// 转类型
	id, _ := strconv.ParseInt(video_id, 10, 64)
	Type, _ := strconv.ParseInt(action_type, 10, 64)
	// 用于同步数据库
	switch Type {
	case 1:
		//设置为0
		Type = 0
	case 2:
		Type = 1
	default:
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  "非法参数",
		})
		return
	}
	// 提取用户Id -- todo 根据中间件提取
	var userId int64 = 1
	flag, err := service.LikeVideoByUserIDService(userId, id, Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if flag {
		c.JSON(http.StatusOK, model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "成功" + action_type,
		})
	}
	fmt.Println(video_id, action_type)

}

// 用户点赞列表
func UserFavoriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	userId, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, userFavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  "请求失败",
			VideoList:  nil,
		})
	}
	userFavoriteList, err := service.UserFavoriteListService(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, userFavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
			VideoList:  nil,
		})
	} else {
		c.JSON(http.StatusOK, userFavoriteListResponse{
			StatusCode: 0,
			StatusMsg:  "成功",
			VideoList:  userFavoriteList,
		})
	}
}
