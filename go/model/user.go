package model

import (
	"log"
	"tiktok/go/config"
)

// User 对应数据库User表结构的结构体
type User struct {
	Id       int64  //主键
	Name     string //昵称
	Password string //密码
}

func CreateUserTable() {
	user := User{}
	Db := config.Init()
	Db.AutoMigrate(&user)
}

// UserInfo  最详细的信息 不与数据库模型对应
type UserInfo struct {
	Id             int64  `json:"id,omitempty"`      //主键
	Name           string `json:"name,omitempty"`    //昵称
	FollowCount    int64  `json:"follow_count"`      //关注总数
	FollowerCount  int64  `json:"follower_count"`    //粉丝总数
	IsFollow       bool   `json:"is_follow"`         //是否关注
	AvatarUrl      string `json:"avatar,omitempty""` //用户的url
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
}

type FriendUser struct {
	UserInfo
	Message string `json:"message"` //聊天信息
	MsgType int64  `json:"msgType"` //message信息的类型，0=>请求用户接受信息，1=>当前请求用户发送的信息
}

var db = config.Init()

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

	if err := db.Debug().Where("Id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	//Db.Close()
	return user, nil
}

// GetUserByName 根据用户名(唯一）查询用户
func GetUserByName(name string) (User, error) {
	user := User{}
	//Db := config.Init()
	// 查数据表
	if err := db.Debug().Where("name = ?", name).Limit(1).Find(&user).Error; err != nil {
		//log.Println(err.Error())
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
		return userInfo, err
	}
	// 查询粉丝总数
	fansCount, err := GetFansById(user.Id)
	if err != nil {
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

// PackageUserToUserInfoByUserId 根据id将user包装成userInfo
func PackageUserToUserInfoByUserId(id int64) (UserInfo, error) {
	userInfo := UserInfo{}

	return userInfo, nil
}
