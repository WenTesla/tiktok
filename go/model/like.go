package model

import (
	"tiktok/go/config"
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
	result := db.Debug().Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("is_cancel", action)
	return result.Error
}

// 插入数据
func InsertLikeData(userId int64, videoId int64) (bool, error) {
	like := Like{
		UserId:  userId,
		VideoId: videoId,
	}
	result := db.Debug().Create(&like)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// QueryDuplicateLikeData 查询是否有重复数据
func QueryDuplicateLikeData(userId int64, videoId int64) (bool, error) {
	like := Like{}
	result := db.Debug().Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Find(&like)
	if result.Error != nil {
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
	//	id IN ( SELECT video_id FROM likes WHERE user_id =?);
	result := db.Debug().Where("id IN (?)", db.Where("user_id = ?", userId).Select("video_id").Find(&Like{})).Find(&tableVideos)
	//log.Printf("%v", tableVideos)
	if result.Error != nil {
		return nil, result.Error
	}
	return tableVideos, nil
}

// QueryLikeByVideoId 根据id获取视频被点赞的总数
func QueryLikeByVideoId(videoId int64) (int64, error) {
	var count int64
	result := db.Debug().Model(&Like{}).Where("video_id = ?", videoId).Count(&count)
	if result.Error != nil {
		return -1, result.Error
	}
	return count, nil
}
