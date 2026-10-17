package model

import (
	"tiktok/go/config"
	"tiktok/go/util"
)

// Like 表的结构,不需要json化
type Like struct {
	Id       int64 `gorm:"primaryKey;column:id"` //自增主键
	UserId   int64 `gorm:"index"`
	VideoId  int64 `gorm:"index"`
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

// QueryLikeCountsByVideoIds 批量取视频点赞数
func QueryLikeCountsByVideoIds(ids []int64) (map[int64]int64, error) {
	m := map[int64]int64{}
	if len(ids) == 0 {
		return m, nil
	}
	type row struct {
		VideoId int64 `gorm:"column:video_id"`
		Cnt     int64 `gorm:"column:cnt"`
	}
	var rows []row
	if err := db.Model(&Like{}).Select("video_id, count(*) as cnt").Where("video_id IN (?) AND is_cancel = ?", ids, 0).Group("video_id").Scan(&rows).Error; err != nil {
		util.LogError(err.Error())
		return m, err
	}
	for _, r := range rows {
		m[r.VideoId] = r.Cnt
	}
	return m, nil
}

// QueryIsLikeBatch 批量查询用户是否点赞给定视频列表，返回 videoId->true
func QueryIsLikeBatch(userId int64, videoIds []int64) (map[int64]bool, error) {
	m := map[int64]bool{}
	if len(videoIds) == 0 {
		return m, nil
	}
	type vid struct {
		VideoId int64 `gorm:"column:video_id"`
	}
	var rows []vid
	if err := db.Model(&Like{}).Select("video_id").Where("user_id = ? AND video_id IN (?) AND is_cancel = ?", userId, videoIds, 0).Find(&rows).Error; err != nil {
		util.LogError(err.Error())
		return m, err
	}
	for _, r := range rows {
		m[r.VideoId] = true
	}
	return m, nil
}

// GetLikeCountsByUserIds 批量取用户的点赞数(点赞他人)与被赞总数(获赞)
func GetLikeCountsByUserIds(ids []int64) (favorite, totalFavorited map[int64]int64, err error) {
	favorite = map[int64]int64{}
	totalFavorited = map[int64]int64{}
	if len(ids) == 0 {
		return favorite, totalFavorited, nil
	}
	type row struct {
		Uid int64 `gorm:"column:uid"`
		Cnt int64 `gorm:"column:cnt"`
	}
	var frows []row
	if err := db.Model(&Like{}).Select("user_id as uid, count(*) as cnt").Where("user_id IN (?) AND is_cancel = ?", ids, 0).Group("user_id").Scan(&frows).Error; err != nil {
		util.LogError(err.Error())
		return favorite, totalFavorited, err
	}
	for _, r := range frows {
		favorite[r.Uid] = r.Cnt
	}
	var trows []row
	if err := db.Table("likes l").Select("v.author_id as uid, count(*) as cnt").
		Joins("JOIN videos v ON l.video_id = v.id").
		Where("l.is_cancel = ? AND v.author_id IN (?)", 0, ids).
		Group("v.author_id").Scan(&trows).Error; err != nil {
		util.LogError(err.Error())
		return favorite, totalFavorited, err
	}
	for _, r := range trows {
		totalFavorited[r.Uid] = r.Cnt
	}
	return favorite, totalFavorited, nil
}
