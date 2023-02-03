package service

import (
	"errors"
	"tiktok/go/model"
)

// FollowUserService 关注用户服务
func FollowUserService(userId int64, toUserId int64, actionType bool) (bool, error) {
	if actionType {
		// 先查询数据
		isExist, err := model.QueryFollowByUserIdAndToUserID(userId, toUserId)
		if err != nil {
			return false, err
		}
		if isExist {
			// 修改数据
			err := model.RefocusUser(userId, toUserId)
			if err != nil {
				return false, err
			}
			return true, nil
		} else {
			// 插入数据
			pass, err := model.InsertFollow(userId, toUserId)
			if err != nil {
				return false, err
			}
			if !pass {
				return false, errors.New("关注失败")
			}
			return true, nil
		}
	} else {
		// 取消关注
		err := model.CancelFollow(userId, toUserId)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
