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
	"strconv"
	"strings"
	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/util"
	"time"
)

var nowTime string

var newFileName string

// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间

func VideoStreamService(lastTime time.Time, userId int64) ([]model.Video, error) {
	tableVideos, err := model.GetVideoByLastTime(lastTime)
	if err != nil {
		log.Printf("失败 %v", err)
		util.LogError(err.Error())
		return nil, dataSourceErr
	}
	log.Printf("获取成功")
	videos, err := packageVideos(tableVideos, userId)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
func VideoInfoByUserId(id int) ([]model.Video, error) {
	tableVideos, err := model.GetVideoByUserId(id)
	if err != nil {
		log.Printf("失败%v", err)
		util.LogError(err.Error())
		return nil, dataSourceErr
	}
	videos, err := packageVideos(tableVideos, -1)
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
func packageVideos(tableVideos []model.TableVideo, userId int64) ([]model.Video, error) {
	// 创建video模型
	videos := make([]model.Video, 0, config.VideoCount)
	if userId == -1 {
		// 填入author
		for _, tableVideo := range tableVideos {
			video, err := packageVideo(&tableVideo)
			if err != nil {
				return nil, err
			}
			videos = append(videos, video)
		}
	} else {
		// 填入author
		for _, tableVideo := range tableVideos {
			video, err := packageVideoWithUserId(&tableVideo, userId)
			if err != nil {
				return nil, err
			}
			videos = append(videos, video)
		}
	}

	return videos, nil
}

// 包装简单的视频列表

func packageSimpleVideos(tableVideos []model.TableVideo, userId int64) ([]model.Video, error) {
	// 创建video模型
	videos := make([]model.Video, 0, config.VideoCount)
	if userId == -1 {
		// 填入author
		for _, tableVideo := range tableVideos {
			video, err := PackSimpleVideoService(&tableVideo)
			if err != nil {
				return nil, err
			}
			videos = append(videos, video)
		}
	} else {
		// 填入author
		for _, tableVideo := range tableVideos {
			video, err := packageVideoWithUserId(&tableVideo, userId)
			if err != nil {
				return nil, err
			}
			videos = append(videos, video)
		}
	}

	return videos, nil
}

//  包装单个视频，不返回是否关注的信息-即未登录状态的信息

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
	// 先查询redis
	favoriteCount, err := likeRedisDb.SCard(strconv.FormatInt(video.ID, 10)).Result()
	if err != nil {
		//// 出错查询数据库
		favoriteCount, err = model.QueryLikeByVideoId(tableVideo.Id)
		if err != nil {
			return video, dataSourceErr
		}
	}
	video.FavoriteCount = favoriteCount
	// 获取"commentCount"
	commentCount, err := model.QueryCommentCountByVideoId(tableVideo.Id)
	if err != nil {
		return video, dataSourceErr
	}
	video.CommentCount = commentCount
	video.IsFavorite = false
	return video, nil
}

//  包装单个视频,返回是否关注的信息

func packageVideoWithUserId(tableVideo *model.TableVideo, id int64) (model.Video, error) {
	// 创建video单例
	video := model.Video{}
	// 获取作者信息
	userInfo, err := UserInfoService(tableVideo.AuthorId, id)
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
	// 先查询redis
	favoriteCount, err := likeRedisDb.SCard(strconv.FormatInt(video.ID, 10)).Result()
	if err != nil {
		//// 出错查询数据库
		favoriteCount, err = model.QueryLikeByVideoId(tableVideo.Id)
		if err != nil {
			return video, dataSourceErr
		}
	}
	video.FavoriteCount = favoriteCount
	// 获取"commentCount"
	commentCount, err := model.QueryCommentCountByVideoId(tableVideo.Id)
	if err != nil {
		return video, dataSourceErr
	}
	video.CommentCount = commentCount
	// 获取是否点赞 先查询redis
	is_favorite, err := likeRedisDb.SIsMember(strconv.FormatInt(video.ID, 10), id).Result()
	// redis查询失败查询数据库
	if err != nil {
		is_favorite, err = model.QueryIsLike(id, tableVideo.Id)
		if err != nil {
			return video, err
		}
	}
	video.IsFavorite = is_favorite
	return video, nil
}

// 包装最简单的视频信息 作者只包含

func PackSimpleVideoService(tableVideo *model.TableVideo) (model.Video, error) {
	// 创建video单例
	video := model.Video{}
	// 获取作者信息
	userInfo, err := SimpleUserService(tableVideo.AuthorId, -1)
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
	video.IsFavorite = false
	return video, nil
}

//  PublishVideo 可以优化

func PublishVideoService(file *multipart.FileHeader, userId int64, title string) error {
	src, err := file.Open()
	if err != nil {
		return errors.New("不能打开文件")
	}
	defer src.Close()
	// 获取视频文件名称
	nowTime = strconv.FormatInt(time.Now().Unix(), 10)
	// 转换视频名称
	newFileName, err = fileNameToTimeCurrentFileName(file.Filename, nowTime)
	// 上传到cos
	err = publishVideoByTencentCos(src, newFileName)
	if err != nil {
		return err
	}
	// 提取封面url
	play_cover, err := parseFileName(newFileName)
	if err != nil {
		return err
	}
	// 添加数据库
	err = model.InsertVideo(userId, config.CosUrl+"/"+newFileName, config.CosUrl+"/"+play_cover, title)
	if err != nil {
		return dataSourceErr
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
		return errors.New("文件上传失败")
	}
	//os.Open()
	return nil
}

// parseFileName parseFileName 解析文件名称，去除文件后缀并加上文件格式jpg
func parseFileName(fileName string) (string, error) {
	//
	lastIndex := strings.LastIndex(fileName, ".")
	if lastIndex == -1 {
		return "", errors.New("解析错误")
	}
	replaced := fileName[lastIndex:]
	// 判断文件后缀是否为要求的后缀
	if replaced != ".mp4" {
		return "", errors.New("文件格式不符合要求")
	}
	return strings.Replace(fileName, replaced, config.ReplaceSuffix, 1), nil
}

// fileNameToTimeCurrentFileName 将文件名称转化为时间戳并返回
func fileNameToTimeCurrentFileName(oldFileName string, newFileName string) (string, error) {
	// 提取文件后缀
	lastIndex := strings.LastIndex(oldFileName, ".")
	if lastIndex == -1 {
		return "", errors.New("文件名称转换错误")
	}
	replaced := oldFileName[:lastIndex]
	return strings.Replace(oldFileName, replaced, newFileName, 1), nil
}

// CheckFile -todo 检查文件内容的合法性
func CheckFile() {
	// 检查文件后缀
	//path.Ext()
}
