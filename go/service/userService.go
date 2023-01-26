package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"tiktok/go/middle/jwt"

	// "log"
	"regexp"
	"tiktok/go/model"
)

// 盐值
const SALT = "TikTok"

// email verify
func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// 根据name生成token
func GenerateTokenByName(username string) (string, error) {
	user, _ := model.GetUserByName(username)
	token := jwt.SignToken(user)
	return token, nil
}

// md5加盐加密
func Encryption(password string) string {
	password += SALT
	hash := md5.New()
	hash.Write([]byte(password))
	hash_password := hex.EncodeToString(hash.Sum(nil))
	return hash_password
}

// 注册服务
func RegisterService(username string, password string) (int64, error) {
	log.Println(username, "---", password)

	// 查表，是否存在id
	user, err := model.GetUserByName(username)
	if err != nil {
		return 0, errors.New("数据库查询错误")
	}
	if username == user.Name {
		return 0, errors.New("用户名已经存在")
	}
	// 插入
	user, _ = model.InsertUser(username, password)
	return user.Id, nil
}

// 登录服务
func LoginService(username string, password string) (int64, error) {

	user, err := model.GetUserByName(username)
	if err != nil {
		return 0, errors.New("数据库查询错误")
	}
	if username != user.Name {
		return 0, errors.New("用户名不存在")
	}
	if password != user.Password {
		return 0, errors.New("用户密码不正确")
	}

	return user.Id, nil
}

// 用户服务 先封装小的，再封装大的
func UserService(Id int64) (model.UserInfo, error) {
	user, err := model.GetUserById(Id)
	if err != nil {
		return model.UserInfo{}, err
	}
	//user.Password = ""
	// 查询自己的关注数目
	followingCount, _ := model.GetFollowingById(Id)
	log.Printf("关注的数量%d", followingCount)
	// 查询自己的粉丝
	fanCount, err := model.GetFansById(Id)
	log.Printf("粉丝数目%d", fanCount)
	// todo 查询是否关注
	userInfo := model.UserInfo{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   followingCount,
		FollowerCount: fanCount,
		IsFollow:      false,
	}
	return userInfo, nil
}
