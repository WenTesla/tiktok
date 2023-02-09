package service

import "tiktok/go/model"

func FriendListService(userId int64) ([]model.FriendUser, error) {
	var FriendUsers []model.FriendUser
	var err error
	FriendUsers, err = PackageFriendLists(userId)
	if err != nil {
		return nil, err
	}
	return FriendUsers, nil
}

func MessageService() {

}

// PackageFriendLists 通过userId查询粉丝的数据，再包装加入消息
func PackageFriendLists(userId int64) ([]model.FriendUser, error) {
	var FriendLists []model.FriendUser
	var message string
	var msgType int8
	userInfos, err := FollowerListService(userId)
	if err != nil {
		return nil, err
	}
	for _, userInfo := range userInfos {
		// 查询Message和MsgType
		message, msgType, err = model.QueryNewestMessageByUserIdAndToUserID(userId, userInfo.Id)
		if err != nil {
			return nil, err
		}
		userInfo.AvatarUrl = "https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/20230206133334.png"
		FriendLists = append(FriendLists, model.FriendUser{
			UserInfo: userInfo,
			Message:  message,
			MsgType:  msgType,
		})
	}

	return FriendLists, nil
}

// PackageFriendList 包装单个请求
func PackageFriendList(userInfo model.FriendUser) (model.UserInfo, error) {

	// test
	userInfo.AvatarUrl = "https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/20230206171653.png"
	return model.UserInfo{}, nil
}

func MessageChatService(userId int64, toUserId int64) ([]model.Message, error) {
	// 查询userid和toUserId的表
	messages, err := model.QueryMessageByUserIdAndToUserId(userId, toUserId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
