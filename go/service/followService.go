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

// FollowListService 关注列表服务
func FollowListService(userId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryFollowUsersByUserId(userId)
	if err != nil {
		return nil, err
	}
	// 定义userInfos 切片
	var userInfos []model.UserInfo
	// 循环
	for _, user := range users {
		userInfo, err := model.PackageUserToUserInfo(user)
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}

// FollowerListService 粉丝列表服务
func FollowerListService(userId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryFansUsersByUserId(userId)
	if err != nil {
		return nil, err
	}
	// 定义userInfos 切片
	var userInfos []model.UserInfo
	// 循环
	for _, user := range users {
		userInfo, err := model.PackageUserToUserInfo(user)
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}
