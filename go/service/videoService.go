package service

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"golang.org/x/net/context"
	"log"
	"mime/multipart"
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
	count, err := model.QueryLikeByVideoId(tableVideo.Id)
	if err != nil {
		return video, err
	}
	video.FavoriteCount = count
	// 获取"comment_count"
	video.CommentCount = 10
	return video, nil
}

// PublishVideo
func PublishVideoService(file *multipart.FileHeader, userId int64, title string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	//配置参数
	putPolicy := storage.PutPolicy{
		Scope: config.Bucket,
	}
	mac := qbox.NewMac(config.AccessKey, config.SecretKey)
	// 获取上传凭证 默认为
	upToken := putPolicy.UploadToken(mac)
	//
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数
	err = formUploader.Put(context.Background(), &ret, upToken, file.Filename, src, file.Size, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(ret.Key, ret.Hash)
	// 添加数据库 -todo cover选择
	err = model.InsertVideo(userId, config.ImgUrl+ret.Key, "https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/image-20230129170217818.png", title)
	if err != nil {
		return err
	}
	return nil
}
