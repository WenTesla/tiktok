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

func CreateUser() {
	user := User{}
	Db := config.Init()
	Db.AutoMigrate(&user)
}

// User 最详细的信息
type UserInfo struct {
	Id             int64  `json:"id,omitempty"`   //主键
	Name           string `json:"name,omitempty"` //昵称
	FollowCount    int64  `json:"follow_count"`   //关注总数
	FollowerCount  int64  `json:"follower_count"` //粉丝总数
	IsFollow       bool   `json:"is_follow"`      //是否关注
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
}

var db = config.Init()

// 插入用户

func InsertUser(name string, password string) (User, error) {
	Db := config.Init()
	//CreateUser()
	user := User{
		Name:     name,
		Password: password,
	}
	result := Db.Create(&user)
	return user, result.Error
	// return true, nil
}

// 根据id(主键）获取用户
func GetUserById(id int) (User, error) {
	user := User{}
	Db := config.Init()

	if err := Db.Debug().Where("Id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}

	return user, nil
}

// 根据用户名(唯一）查询用户
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
