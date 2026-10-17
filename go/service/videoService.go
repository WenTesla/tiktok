package service

import (
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/util"
)

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

// 包装视频流，批量填充作者信息/点赞数/评论数/是否点赞，消除 N+1
//
//	author
//	"favorite_count"
//	"comment_count"
//	"is_favorite"
func packageVideos(tableVideos []model.TableVideo, userId int64) ([]model.Video, error) {
	videos := make([]model.Video, 0, len(tableVideos))
	if len(tableVideos) == 0 {
		return videos, nil
	}

	// 收集去重 authorIds 与 videoIds
	authorIds := make([]int64, 0, len(tableVideos))
	videoIds := make([]int64, 0, len(tableVideos))
	seenAuthor := make(map[int64]struct{})
	for i := range tableVideos {
		videoIds = append(videoIds, tableVideos[i].Id)
		if _, ok := seenAuthor[tableVideos[i].AuthorId]; !ok {
			seenAuthor[tableVideos[i].AuthorId] = struct{}{}
			authorIds = append(authorIds, tableVideos[i].AuthorId)
		}
	}

	// 1) 批量取作者信息（含关注/粉丝/作品/获赞数、是否关注 viewer）
	authors, err := UserServiceBatch(authorIds, userId)
	if err != nil {
		return nil, err
	}
	// 2) 批量取视频点赞数
	likeCounts, err := model.QueryLikeCountsByVideoIds(videoIds)
	if err != nil {
		return nil, dataSourceErr
	}
	// 3) 批量取视频评论数
	commentCounts, err := model.QueryCommentCountsByVideoIds(videoIds)
	if err != nil {
		return nil, dataSourceErr
	}
	// 4) 批量取是否点赞（仅登录态）
	var isLikeMap map[int64]bool
	if userId != -1 {
		isLikeMap, err = model.QueryIsLikeBatch(userId, videoIds)
		if err != nil {
			return nil, err
		}
	}

	// 循环内仅做 O(N) map 查找，不再有 DB/Redis 往返
	// 注意：通过索引取址，修复 Go<1.22 循环变量取址错位 bug
	for i := range tableVideos {
		tv := &tableVideos[i]
		videos = append(videos, model.Video{
			Author:        authors[tv.AuthorId],
			ID:            tv.Id,
			PlayURL:       tv.PlayUrl,
			CoverURL:      tv.CoverUrl,
			Title:         tv.Title,
			FavoriteCount: likeCounts[tv.Id],
			CommentCount:  commentCounts[tv.Id],
			IsFavorite:    userId != -1 && isLikeMap[tv.Id],
		})
	}
	return videos, nil
}

// PublishVideoService 发布视频：落盘到本地存储 + 服务端抽首帧生成封面 + 入库。
// 不依赖腾讯云 COS（COS 鉴权有 bug 且 scratch 镜像缺 CA 证书，实际不可用）。
func PublishVideoService(file *multipart.FileHeader, userId int64, title string) error {
	// 上传前校验文件内容（仅允许 mp4，按文件头魔数判断）
	if err := CheckFile(file); err != nil {
		return err
	}

	key := genVideoKey(userId)
	videoRel := "videos/" + key + ".mp4"
	coverRel := "covers/" + key + ".jpg"

	// 1) 视频落盘到本地存储（宿主卷挂载，对外通过 /static 静态路由可达）
	src, err := file.Open()
	if err != nil {
		return errors.New("不能打开文件")
	}
	defer src.Close()
	playURL, err := config.Storage.Put(videoRel, src)
	if err != nil {
		return errors.New("视频存储失败")
	}

	// 2) 抽首帧生成封面缩略图（服务端内嵌 ffmpeg，缩放至 480px 宽）
	videoLocal := config.Storage.LocalPath(videoRel)
	coverLocal := config.Storage.LocalPath(coverRel)
	coverURL := ""
	if err := util.ExtractFirstFrame(videoLocal, coverLocal, 480); err != nil {
		// 抽帧失败不阻断视频发布：封面留空，前端可降级处理
		util.LogError("封面抽帧失败，视频仍正常发布: " + err.Error())
	} else {
		coverURL = strings.TrimSuffix(config.Storage.PublicBaseURL, "/") + "/" + coverRel
	}

	// 3) 入库
	if err := model.InsertVideo(userId, playURL, coverURL, title); err != nil {
		return dataSourceErr
	}
	return nil
}

// genVideoKey 生成唯一视频 key，避免文件名碰撞（用户ID + 纳秒时间戳）。
func genVideoKey(userId int64) string {
	return strconv.FormatInt(userId, 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

// CheckFile 校验上传文件内容合法性：按文件头魔数判断，仅允许 mp4
func CheckFile(fh *multipart.FileHeader) error {
	src, err := fh.Open()
	if err != nil {
		return errors.New("不能打开文件")
	}
	defer src.Close()
	head := make([]byte, 512)
	n, _ := src.Read(head)
	if util.GetFileType(head[:n]) != "mp4" {
		return errors.New("仅支持 mp4 格式视频")
	}
	return nil
}
