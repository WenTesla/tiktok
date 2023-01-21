package model

import (
	"log"
	"tiktok/go/config"
)

// User 对应数据库User表结构的结构体
type User struct {
	Id       int64
	Name     string
	Password string
}

func CreateUser() {
	user := User{}
	Db := config.Init()
	Db.AutoMigrate(&user)

}

// 插入用户

func InsertUser(name string, password string) (int, error) {
	Db := config.Init()
	CreateUser()
	user := User{
		Name:     name,
		Password: password,
	}
	result := Db.Create(&user)
	return int(user.Id), result.Error
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
	Db := config.Init()
	// 查数据表
	if err := Db.Debug().Where("name = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}
