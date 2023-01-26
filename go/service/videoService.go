package service

import (
	"log"
	"tiktok/go/model"
	"time"
)

// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
func VideoStreamService(lastTime time.Time, userId int64) ([]model.Video, error) {
	tableVideos, err := model.GetVideoByLastTime(lastTime)
	if err != nil {
		log.Printf("失败 %v", err)
	}
	log.Printf("获取成功")
	videos, err := packageVideos(tableVideos)
	// 获取发布最早的时间
	nextTime, err := model.GetVideoNextTime(lastTime)
	if err != nil {
		return nil, err
	}
	log.Printf("%v", nextTime)
	return videos, nil
}

// 包装视频，填入内容
// user
func packageVideos(video []model.TableVideo) ([]model.Video, error) {

	return nil, nil
}
