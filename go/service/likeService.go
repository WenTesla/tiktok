package service

import "tiktok/go/model"

// LikeVideoByUserIDService 用户点赞服务
func LikeVideoByUserIDService(userId int64, videoId int64, actionType int64) (bool, error) {
	// 点赞 -- 没有数据的话创建数据，存在数据的话插入数据
	// 检查是否存在当前的重复值
	isDuplicate, err := model.QueryDuplicateLikeData(userId, videoId)
	// 为查询相关数据或者数据查询错误
	if err != nil {
		return false, err
	}
	// 包含相关数据
	if isDuplicate {
		// 包含相关数据-更新数据
		err := model.UpdateLikeVideoByUserId(userId, videoId, actionType)
		if err != nil {
			return false, err
		}
	} else {
		// 添加到redis
		// 插入数据
		_, err := model.InsertLikeData(userId, videoId)
		if err != nil {
			return false, err
		}
	}
	return true, err
}

// UserFavoriteListService 用户喜欢列表服务
func UserFavoriteListService(userId int64) ([]model.Video, error) {
	// 根据用户查询点赞的视频
	tableVideos, err := model.QueryVideoByUserId(userId)
	if err != nil {
		return nil, err
	}
	videos, err := packageVideos(tableVideos, -1)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
