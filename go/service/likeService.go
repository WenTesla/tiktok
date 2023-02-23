package service

import (
	"errors"
	"strconv"
	"tiktok/go/model"
)

//  用户点赞服务

func LikeVideoByUserIDService(userId int64, videoId int64, actionType int64) (bool, error) {
	var err error = nil
	var isDuplicate bool
	videoIdString := strconv.FormatInt(videoId, 10)
	//// 点赞 -- 没有数据的话创建数据，存在数据的话插入数据

	//// 从redis判断是否点过赞
	//isMember, err := likeRedisDb.SIsMember(videoIdString, userId).Result()
	//// redis 正常
	//if isMember {
	//	// 查询到数据
	//	// 如果是点赞，直接返回
	//	if actionType == 0 {
	//		return true, nil
	//	} else if actionType == 1 {
	//		// 如果是取消点赞
	//		// 删除redis的操作
	//		likeRedisDb.SRem(videoIdString, userId)
	//		// 操作数据库
	//		go model.UpdateLikeVideoByUserId(userId, videoId, actionType)
	//		return true, nil
	//	}
	//} else {
	//	// 未查询到数据
	//	// 如果是点赞
	//	if actionType == 0 {
	//		// 添加数据到redis
	//		likeRedisDb.SAdd(videoIdString, userId)
	//		// 添加数据到数据库
	//		go model.InsertLikeData(userId, videoId)
	//		return true, nil
	//	} else if actionType == 1 {
	//		// 如果是取消点赞
	//		go model.UpdateLikeVideoByUserId(userId, videoId, actionType)
	//		return true, nil
	//	}
	//}
	// redis 操作异常 只操作数据库的做法
	// 检查是否存在当前的重复值
	isDuplicate, err = model.QueryDuplicateLikeData(userId, videoId)
	// 为查询相关数据或者数据查询错误
	if err != nil {
		return false, dataSourceErr
	}
	// 包含相关数据
	if isDuplicate {
		// 包含相关数据-更新数据
		err = model.UpdateLikeVideoByUserId(userId, videoId, actionType)
		if err != nil {
			return false, dataSourceErr
		}
		if actionType == 0 {
			// 添加redis
			likeRedisDb.SAdd(videoIdString, userId)
		} else if actionType == 1 {
			// 移除redis
			likeRedisDb.SRem(videoIdString, userId)
		}
	} else {
		if actionType == 0 {
			// 插入数据
			// 先检查视频的id是否存在
			IsExist, err := model.QueryIsExistVideoId(videoId)
			if IsExist != true {
				return false, errors.New("视频不存在")
			}
			_, err = model.InsertLikeData(userId, videoId)
			if err != nil {
				return false, dataSourceErr
			}
			// 添加到redis
			likeRedisDb.SAdd(videoIdString, userId)
		}
	}
	return true, err
}

//  用户喜欢列表服务

func UserFavoriteListService(userId int64) ([]model.Video, error) {
	// 根据用户查询点赞的视频
	tableVideos, err := model.QueryVideoByUserId(userId)
	if err != nil {
		return nil, dataSourceErr
	}
	videos, err := packageVideos(tableVideos, -1)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
