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
	CreatedAt  time.Time
}

// 通过id获取自己关注的数目
func GetFollowingById(id int64) (int64, error) {
	var count int64
	result := db.Debug().Model(&Follow{}).Where("user_id = ?", id).Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// 通过id获取自己的粉丝数
func GetFansById(id int64) (int64, error) {
	var count int64
	result := db.Debug().Model(&Follow{}).Where("follower_id = ?", id).Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// 通过id查看是否关注
func IsFollowingById(id int64) (bool, error) {
	//var count int64

	return true, nil
}
