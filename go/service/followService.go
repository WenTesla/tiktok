package service

// FollowUserService
func FollowUserService(userId int64, toUserId int64, actionType bool) (bool, error) {
	if actionType {
		// 插入数据

	} else {
		// 取消关注

	}
	return false, nil
}
