package model

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
