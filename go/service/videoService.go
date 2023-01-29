package service

import (
	"log"
	"tiktok/go/config"
	"tiktok/go/model"
	"time"
)

// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
func VideoStreamService(lastTime time.Time) ([]model.Video, error) {
	tableVideos, err := model.GetVideoByLastTime(lastTime)
	if err != nil {
		log.Printf("失败 %v", err)
	}
	log.Printf("获取成功")
	videos, err := packageVideos(tableVideos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
func VideoInfoByUserId(id int) ([]model.Video, error) {
	tableVideos, err := model.GetVideoByUserId(id)
	if err != nil {
		log.Printf("失败%v", err)
	}
	videos, err := packageVideos(tableVideos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 包装视频流，填入内容如下
//
//		author
//	"favorite_count": 0,
//	 "comment_count": 0,
//	 "is_favorite": true,
//
// user
func packageVideos(tableVideos []model.TableVideo) ([]model.Video, error) {
	// 创建video模型
	videos := make([]model.Video, 0, config.VideoCount)
	// 填入author
	for _, tableVideo := range tableVideos {
		video, err := packageVideo(&tableVideo)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// 包装单个视频
func packageVideo(tableVideo *model.TableVideo) (model.Video, error) {
	// 创建video单例
	video := model.Video{}
	// 获取作者信息
	userInfo, err := UserService(tableVideo.AuthorId)
	if err != nil {
		return model.Video{}, err
	}
	log.Printf("%v", userInfo)
	//video.Author=user
	video.Author = userInfo
	// 填充Videos的
	video.ID = tableVideo.Id
	video.PlayURL = tableVideo.PlayUrl
	video.CoverURL = tableVideo.CoverUrl
	video.Title = tableVideo.Title
	// 获取 favorite_count
	video.FavoriteCount = 10
	// 获取"comment_count"
	video.CommentCount = 10
	return video, nil
}
