package model

import (
	"tiktok/go/util"
	"time"
)

// Follow 用户关系结构，对应用户关系表。
type Follow struct {
	Id         int64 `gorm:"primaryKey;column:id"`
	UserId     int64 `gorm:"index"`
	FollowerId int64 `gorm:"index"`
	Cancel     int8
	CreatedAt  time.Time `gorm:"-"`
}

// GetFollowingById 通过id获取自己关注的数目 同时筛选出cancel
func GetFollowingById(id int64) (int64, error) {
	var count int64
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)
	result := db.Model(&Follow{}).Where("user_id = ? AND cancel = ?", id, 0).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetFansById 通过id获取自己的粉丝数 同时筛选出cancel
func GetFansById(id int64) (int64, error) {
	var count int64
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)
	result := db.Model(&Follow{}).Where("follower_id = ? AND cancel = ?", id, 0).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
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
	// INSERT INTO `follows` (`user_id`,`follower_id`,`cancel`) VALUES (7,9,0)
	result := db.Create(&follow)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return false, result.Error
	}
	return true, nil
}

// QueryFollowByUserIdAndToUserID 根据用户id和关注的id查询是否存在
func QueryFollowByUserIdAndToUserID(userId int64, toUserID int64) (bool, error) {
	//SELECT * FROM `follows` WHERE user_id = 1 AND follower_id = 5
	result := db.Where("user_id = ? AND follower_id = ?", userId, toUserID).Find(&Follow{})
	if result.Error != nil {
		util.LogError(result.Error.Error())
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
	//UPDATE `follows` SET `cancel`=1 WHERE user_id = 1 AND follower_id = 5
	result := db.Debug().Model(&Follow{}).Where("user_id = ? AND follower_id = ?", userId, toUserID).Update("cancel", 1)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return result.Error
	}
	return nil
}

// RefocusUser 重新关注用户
func RefocusUser(userId int64, toUserID int64) error {
	// 更新单列
	// UPDATE `follows` SET `cancel`=0 WHERE user_id = 1 AND follower_id = 5
	result := db.Model(&Follow{}).Where("user_id = ? AND follower_id = ?", userId, toUserID).Update("cancel", 0)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return result.Error
	}
	return nil
}

//  根据用户id查询用户关注的用户id切片 比较绕

func QueryFollowUsersByUserId(userId int64) ([]User, error) {
	var users []User
	//SELECT * FROM `users` WHERE id IN (SELECT `follower_id` FROM `follows` WHERE user_id = 1 AND cancel = 0)
	result := db.Where("id IN (?)", db.Where("user_id = ? AND cancel = ?", userId, 0).Select("follower_id").Find(&Follow{})).Find(&users)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return nil, result.Error
	}
	return users, nil
}

//  根据用户id查询当前用户的粉丝id切片 比较绕

func QueryFansUsersByUserId(userId int64) ([]User, error) {
	var users []User
	//SELECT * FROM `users` WHERE id IN (SELECT `user_id` FROM `follows` WHERE follower_id = 1 AND cancel = 0)
	result := db.Where("id IN (?)", db.Where("follower_id = ? AND cancel = ?", userId, 0).Select("user_id").Find(&Follow{})).Find(&users)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return nil, result.Error
	}
	return users, nil
}

/*
查询互相关注的用户列表
*/

func QueryMutualFollowListByUserId(userId int64) ([]User, error) {

	// inner join
	//db.InnerJoins("Company").Find(&users)
	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` INNER JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`

	var MutualFollowList []User
	// SELECT a.follower_id FROM follows a join follows b on a.user_id=b.follower_id AND a.follower_id = b.user_id WHERE a.cancel = 0 AND b.cancel = 0 AND a.user_id = ? AND b.follower_id = ?
	//result := db.Debug().Table("follows a").Select("a.follower_id").Joins("join follows b on a.user_id=b.follower_id AND a.follower_id = b.user_id").Where("a.cancel = ? AND b.cancel = ? AND a.user_id = ? AND b.follower_id = ?", 0, 0, userId, userId).Find(&MutualFollowList)
	// SELECT `id` `name` FROM `users` WHERE id in (SELECT a.follower_id FROM follows a join follows b on a.user_id=b.follower_id AND a.follower_id = b.user_id WHERE a.cancel = 0 AND b.cancel = 0 AND a.user_id = 1 AND b.follower_id = 1)
	result := db.Where("id in (?)", db.Table("follows a").Select("a.follower_id").Joins("join follows b on a.user_id=b.follower_id AND a.follower_id = b.user_id").Where("a.cancel = ? AND b.cancel = ? AND a.user_id = ? AND b.follower_id = ?", 0, 0, userId, userId).Find(&Follow{})).Select("id", "name").Find(&MutualFollowList)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return nil, result.Error
	}
	return MutualFollowList, nil
}

//  查询是否关注 第一个参数为当前用户的id，第二个参数为要关注的用户Id

func QueryIsFollow(userId int64, toUserId int64) (bool, error) {
	// 自己不能关注自己
	if userId == toUserId {
		return true, nil
	}
	var count int64
	//SELECT count(*) FROM `follows` WHERE user_id = ? AND follower_id = ? AND cancel = 0
	result := db.Model(&Follow{}).Where("user_id = ? AND follower_id = ? AND cancel = ?", userId, toUserId, 0).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return false, result.Error
	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}

// GetFollowCountsByUserIds 批量取用户的关注数/粉丝数
func GetFollowCountsByUserIds(ids []int64) (following, fans map[int64]int64, err error) {
	following = map[int64]int64{}
	fans = map[int64]int64{}
	if len(ids) == 0 {
		return following, fans, nil
	}
	type row struct {
		Uid int64 `gorm:"column:uid"`
		Cnt int64 `gorm:"column:cnt"`
	}
	var frows []row
	if err := db.Model(&Follow{}).Select("user_id as uid, count(*) as cnt").Where("user_id IN (?) AND cancel = ?", ids, 0).Group("user_id").Scan(&frows).Error; err != nil {
		util.LogError(err.Error())
		return following, fans, err
	}
	for _, r := range frows {
		following[r.Uid] = r.Cnt
	}
	var arows []row
	if err := db.Model(&Follow{}).Select("follower_id as uid, count(*) as cnt").Where("follower_id IN (?) AND cancel = ?", ids, 0).Group("follower_id").Scan(&arows).Error; err != nil {
		util.LogError(err.Error())
		return following, fans, err
	}
	for _, r := range arows {
		fans[r.Uid] = r.Cnt
	}
	return following, fans, nil
}

// QueryIsFollowBatch 批量查询 viewer 是否关注给定作者列表，返回 authorId->true
func QueryIsFollowBatch(viewerId int64, authorIds []int64) (map[int64]bool, error) {
	m := map[int64]bool{}
	if len(authorIds) == 0 {
		return m, nil
	}
	type fid struct {
		FollowerId int64 `gorm:"column:follower_id"`
	}
	var rows []fid
	if err := db.Model(&Follow{}).Select("follower_id").Where("user_id = ? AND follower_id IN (?) AND cancel = ?", viewerId, authorIds, 0).Find(&rows).Error; err != nil {
		util.LogError(err.Error())
		return m, err
	}
	for _, r := range rows {
		m[r.FollowerId] = true
	}
	return m, nil
}

//  查询粉丝的Id

func QueryFansId(userId int64) ([]int64, error) {

	return nil, nil
}
