package model

import (
	"tiktok/go/util"
)

// User 对应数据库User表结构的结构体
type User struct {
	Id       int64  //主键
	Name     string //昵称
	Password string //密码
}

// UserInfo  最详细的信息 不与数据库模型对应
type UserInfo struct {
	Id             int64  `json:"id,omitempty"`              //主键
	Name           string `json:"name,omitempty"`            //昵称
	FollowCount    int64  `json:"follow_count"`              //关注总数
	FollowerCount  int64  `json:"follower_count"`            //粉丝总数
	IsFollow       bool   `json:"is_follow"`                 //是否关注
	AvatarUrl      string `json:"avatar,omitempty"`          //用户的url
	TotalFavorited int64  `json:"total_favorited,omitempty"` //获赞数量
	WorkCount      int64  `json:"work_count,omitempty"`      //作品数量
	FavoriteCount  int64  `json:"favorite_count,omitempty"`  //点赞数量
}

type FriendUser struct {
	UserInfo
	Message string `json:"message"` //聊天信息
	MsgType int8   `json:"msgType"` //message信息的类型，0=>请求用户接受信息，1=>当前请求用户发送的信息
}

// InsertUser 插入用户
func InsertUser(name string, password string) (User, error) {
	//CreateUser()
	user := User{
		Name:     name,
		Password: password,
	}
	result := db.Create(&user)
	return user, result.Error
	// return true, nil
}

// GetUserById 根据id(主键）获取用户
func GetUserById(id int64) (User, error) {
	user := User{}
	// SELECT * FROM `users` WHERE Id = 9 ORDER BY `users`.`id` LIMIT 1
	if err := db.Where("Id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	//Db.Close()
	return user, nil
}

// GetUserByName 根据用户名(唯一）查询用户
func GetUserByName(name string) (User, error) {
	user := User{}
	//Db := config.InitDataSource()
	// 查数据表
	//  SELECT * FROM `users` WHERE name = '周子豪' LIMIT 1
	if err := db.Where("name = ?", name).Limit(1).Find(&user).Error; err != nil {
		//log.Println(err.Error())
		util.LogError(err.Error())
		return user, err
	}
	return user, nil
}

// PackageUserToUserInfo
//
//	 	FollowCount
//		FollowerCount
//		IsFollow
func PackageUserToUserInfo(user User) (UserInfo, error) {
	userInfo := UserInfo{}
	//var err error
	// 查询关注总数
	followingsCount, err := GetFollowingById(user.Id)
	if err != nil {
		util.LogError(err.Error())

		return userInfo, err
	}
	// 查询粉丝总数
	fansCount, err := GetFansById(user.Id)
	if err != nil {
		util.LogError(err.Error())
		return userInfo, err
	}
	// to-do 查询是否关注

	// 合并
	userInfo.Id = user.Id
	userInfo.Name = user.Name
	userInfo.FollowCount = followingsCount
	userInfo.FollowerCount = fansCount
	userInfo.IsFollow = false

	return userInfo, nil
}

//
//
//	IsFollow -todo 点赞和作品总数

func PackageUserToSimpleUserInfo(user User, userId int64) (UserInfo, error) {
	userInfo := UserInfo{}
	//  查询是否关注
	isFollow, err := QueryIsFollow(userId, user.Id)
	if err != nil {
		util.LogError(err.Error())
		return userInfo, err
	}
	// 合并
	userInfo.Id = user.Id
	userInfo.Name = user.Name
	userInfo.FollowCount = 0
	userInfo.FollowerCount = 0
	userInfo.IsFollow = isFollow
	return userInfo, nil
}

// PackageUserToUserInfoByUserId 根据id将user包装成userInfo
func PackageUserToUserInfoByUserId(id int64) (UserInfo, error) {
	userInfo := UserInfo{}

	return userInfo, nil
}

// 直接将user转换为userInfo,不查询数据库，只是包装

func PackageUserToDirectUserInfo(user User) (UserInfo, error) {
	userInfo := UserInfo{}
	userInfo.Id = user.Id
	userInfo.Name = user.Name
	userInfo.IsFollow = true
	return userInfo, nil
}
