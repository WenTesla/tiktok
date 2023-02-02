package service

import (
	"errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/net/context"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	// 获取"commentCount"
	commentCount, err := model.QueryCommentCountByVideoId(tableVideo.Id)
	if err != nil {
		return video, err
	}
	video.CommentCount = commentCount
	return video, nil
}

// PublishVideo
func PublishVideoService(file *multipart.FileHeader, userId int64, title string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// 获取视频文件名称
	// 上传到cos
	err = publishVideoByTencentCos(src, file.Filename)
	if err != nil {
		return err
	}
	// 提取封面url
	play_cover, err := parseFileName(file.Filename)
	if err != nil {
		return err
	}
	// 添加数据库
	err = model.InsertVideo(userId, config.CosUrl+"/"+file.Filename, config.CosUrl+"/"+play_cover, title)
	if err != nil {
		return err
	}
	return nil
}

// publishVideoByTencentCos
func publishVideoByTencentCos(file multipart.File, fileName string) error {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为真实的信息
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse(config.CosUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(config.SecretId),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: os.Getenv(config.SecretKey), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	_, err := c.Object.Put(context.Background(), fileName, file, nil)
	if err != nil {
		return err
	}

	//os.Open()
	return nil
}

// parseFileName 解析文件名称，去除加上文件格式jpg
func parseFileName(fileName string) (string, error) {
	//
	lastIndex := strings.LastIndex(fileName, ".")
	if lastIndex == -1 {
		return "", errors.New("解析错误")
	}
	replaced := fileName[lastIndex:]
	return strings.Replace(fileName, replaced, config.Replace, 1), nil
}
