package service

import (
	"errors"
	"tiktok/go/model"
)

//  关注用户服务

func FollowUserService(userId int64, toUserId int64, actionType bool) (bool, error) {
	if actionType {
		// 先查询数据
		isExist, err := model.QueryFollowByUserIdAndToUserID(userId, toUserId)
		if err != nil {
			return false, dataSourceErr
		}
		if isExist {
			// 修改数据
			err = model.RefocusUser(userId, toUserId)
			if err != nil {
				return false, dataSourceErr
			}
			return true, nil
		} else {
			// 插入数据
			pass, err := model.InsertFollow(userId, toUserId)
			if err != nil {
				return false, dataSourceErr
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
			return false, dataSourceErr
		}
		return true, nil
	}
	return false, nil
}

//  未登录状态下的关注列表服务

func FollowListService(userId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryFollowUsersByUserId(userId)
	if err != nil {
		return nil, dataSourceErr
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

//  未登录状态下的粉丝列表服务

func FollowerListService(userId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryFansUsersByUserId(userId)
	if err != nil {
		return nil, dataSourceErr
	}
	// 定义userInfos 切片
	var userInfos []model.UserInfo
	// 循环
	for _, user := range users {
		userInfo, err := model.PackageUserToUserInfo(user)
		userInfo.IsFollow = false
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}

//  登录状态下的关注列表 第一个参数为  第二个参数为登录用户的Id

func FollowListServiceWithUserId(userId int64, loginUserId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryFollowUsersByUserId(userId)
	if err != nil {
		return nil, dataSourceErr
	}
	// 定义userInfos 切片
	var userInfos []model.UserInfo
	// 循环
	for _, user := range users {
		// 判断对方是否关注

		userInfo, err := model.PackageUserToSimpleUserInfo(user, loginUserId)
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}

//  登录状态下的粉丝列表

func FollowerListServiceWithUserId(userId int64, loginUserId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryFansUsersByUserId(userId)
	if err != nil {
		return nil, dataSourceErr
	}
	// 定义userInfos 切片
	var userInfos []model.UserInfo
	// 循环
	for _, user := range users {
		userInfo, err := model.PackageUserToSimpleUserInfo(user, loginUserId)
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}

// 互相关注的好友列表

func MutualFollowListService(userId int64) ([]model.UserInfo, error) {
	// 先根据用户Id取用户关注
	users, err := model.QueryMutualFollowListByUserId(userId)
	if err != nil {
		return nil, dataSourceErr
	}
	// 定义userInfos 切片
	var userInfos []model.UserInfo
	// 循环
	for _, user := range users {
		userInfo, _ := model.PackageUserToDirectUserInfo(user)
		userInfos = append(userInfos, userInfo)
	}
	return userInfos, nil
}
