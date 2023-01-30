package model

import (
	"tiktok/go/config"
	"time"
)

type ApifoxModel struct {
	NextTime   *int64  `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 视频列表
}

// 对应sql表格的实体 更加简单的配置
type TableVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"` //视频名，5.23添加
}

// Video
type Video struct {
	Author        UserInfo `json:"author"`         // 视频作者信息
	CommentCount  int64    `json:"comment_count"`  // 视频的评论总数
	CoverURL      string   `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"` // 视频的点赞总数
	ID            int64    `json:"id"`             // 视频唯一标识
	IsFavorite    bool     `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string   `json:"play_url"`       // 视频播放地址
	Title         string   `json:"title"`          // 视频标题
}

func (TableVideo) TableName() string {
	return "videos"
}

//var db = config.Init()

// 获取时间戳之前的视频
func GetVideoByLastTime(lastTime time.Time) ([]TableVideo, error) {
	tableVideos := make([]TableVideo, config.VideoCount)
	result := db.Debug().Where("publish_time <= ?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&tableVideos)
	if result.Error != nil {
		return nil, result.Error
	}
	return tableVideos, nil
}

// 通过用户id获取视频
func GetVideoByUserId(userId int) ([]TableVideo, error) {
	tableVideos := make([]TableVideo, config.VideoMaxCount)
	db.Debug().Where("author_id = ?", userId).Order("publish_time desc").Limit(config.VideoMaxCount).Find(&tableVideos)
	return tableVideos, nil
}

// 获取发布最早的视频的时间戳，作为下次请求的时间戳
func GetVideoNextTime(lastTime time.Time) (time.Time, error) {
	tableVideo := TableVideo{}
	result := db.Debug().Where("publish_time <= ?", lastTime).Order("publish_time asc").Limit(1).Select("publish_time").Find(&tableVideo)
	if result.Error != nil {
		return time.Time{}, result.Error
	}
	return tableVideo.PublishTime, nil
}
