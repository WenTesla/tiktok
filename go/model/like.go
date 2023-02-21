package model

import (
	"tiktok/go/config"
	"tiktok/go/util"
)

// Like 表的结构,不需要json化
type Like struct {
	Id       int64 //自增主键
	UserId   int64 //点赞用户的id
	VideoId  int64 //视频的id
	IsCancel int8  //是否点赞，0为点赞，1为取消赞
}

// 更新点赞数据
func UpdateLikeVideoByUserId(userId int64, videoId int64, action int64) error {
	//更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	// UPDATE `likes` SET `is_cancel`=0 WHERE `user_id` = 9 AND `video_id` = 39
	result := db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("is_cancel", action)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return result.Error
	}
	return nil
}

// 插入数据
func InsertLikeData(userId int64, videoId int64) (bool, error) {
	like := Like{
		UserId:  userId,
		VideoId: videoId,
	}
	result := db.Create(&like)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return false, result.Error
	}
	return true, nil
}

// QueryDuplicateLikeData 查询是否有重复数据
func QueryDuplicateLikeData(userId int64, videoId int64) (bool, error) {
	like := Like{}
	// SELECT * FROM `likes` WHERE `user_id` = 9 AND `video_id` = 39
	result := db.Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Find(&like)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		// 查询错误
		return false, result.Error
	}
	if like.Id == 0 {
		return false, nil
	}
	return true, nil
}

// QueryVideoByUserId 根据用户id查询视频的信息
func QueryVideoByUserId(userId int64) ([]TableVideo, error) {
	tableVideos := make([]TableVideo, config.VideoMaxCount) //
	// SELECT
	//	*
	//FROM
	//	videos
	//WHERE
	//	id IN ( SELECT video_id FROM likes WHERE user_id =? AND is_cancel = 0);
	result := db.Where("id IN (?)", db.Where("user_id = ? AND is_cancel = ?", userId, 0).Select("video_id").Find(&Like{})).Find(&tableVideos)
	//log.Printf("%v", tableVideos)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return nil, result.Error
	}
	return tableVideos, nil
}

// QueryLikeByVideoId 根据id获取视频被点赞的总数
func QueryLikeByVideoId(videoId int64) (int64, error) {
	var count int64
	//  SELECT count(*) FROM `likes` WHERE video_id = ? AND is_cancel = 0
	result := db.Model(&Like{}).Where("video_id = ? AND is_cancel = ?", videoId, 0).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return -1, result.Error
	}
	return count, nil
}

//  查询是否有点赞

func QueryIsLike(userId int64, videoId int64) (bool, error) {
	like := Like{}
	// SELECT * FROM `likes` WHERE `user_id` = 9 AND `video_id` = 39 AND is_cancel = 0
	result := db.Where(map[string]interface{}{"user_id": userId, "video_id": videoId, "is_cancel": 0}).Find(&like)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		// 查询错误
		return false, result.Error
	}
	if like.Id == 0 {
		return false, nil
	}
	return true, nil
}

// 查询点赞数量

func QueryFavoriteCountByUserId(userId int64) (int64, error) {
	var count int64
	// SELECT count(*) FROM `likes` WHERE user_id = 1 AND is_cancel = 0
	result := db.Model(&Like{}).Where("user_id = ? AND is_cancel = ?", userId, 0).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return -1, result.Error
	}
	return count, nil
}

// 获取获赞数量

func QueryTotalFavorited(userId int64) (int64, error) {
	var count int64
	// SELECT count(*) FROM `likes` WHERE video_id in (SELECT `id` FROM `videos` WHERE author_id = 1) AND is_cancel = 0
	db.Model(&Like{}).Where("video_id in (?) AND is_cancel = ?", db.Model(&TableVideo{}).Where("author_id = ?", userId).Select("id"), 0).Count(&count)
	return count, nil
}
