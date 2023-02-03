package model

import (
	"time"
)

// Follow 用户关系结构，对应用户关系表。
type Follow struct {
	Id         int64
	UserId     int64
	FollowerId int64
	Cancel     int8
	CreatedAt  time.Time `gorm:"-"`
}

// GetFollowingById 通过id获取自己关注的数目 同时筛选出cancel
func GetFollowingById(id int64) (int64, error) {
	var count int64
	result := db.Debug().Model(&Follow{}).Where("user_id = ? AND cancel = ?", id, 0).Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// GetFansById 通过id获取自己的粉丝数 同时筛选出cancel
func GetFansById(id int64) (int64, error) {
	var count int64
	result := db.Debug().Model(&Follow{}).Where("follower_id = ? AND cancel = ?", id, 0).Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// IsFollowingById 通过id查看是否关注
func IsFollowingById(id int64) (bool, error) {
	//var count int64

	return true, nil
}

// InsertFollow 插入关注
func InsertFollow(userId int64, toUserID int64) (bool, error) {
	follow := Follow{
		UserId:     userId,
		FollowerId: toUserID,
	}
	result := db.Debug().Create(&follow)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// QueryFollowByUserIdAndToUserID 根据用户id和关注的id查询是否存在
func QueryFollowByUserIdAndToUserID(userId int64, toUserID int64) (bool, error) {
	result := db.Debug().Where("user_id = ? AND follower_id = ?", userId, toUserID).Find(&Follow{})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// CancelFollow 取消关注
func CancelFollow(userId int64, toUserID int64) error {
	// 更新单列
	result := db.Debug().Model(&Follow{}).Where("user_id = ? AND follower_id = ?", userId, toUserID).Update("cancel", 1)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// RefocusUser 重新关注用户
func RefocusUser(userId int64, toUserID int64) error {
	// 更新单列
	result := db.Debug().Model(&Follow{}).Where("user_id = ? AND follower_id = ?", userId, toUserID).Update("cancel", 0)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
