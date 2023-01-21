package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"

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

// md5加盐加密
func encryption(password string) string {
	password += SALT
	hash := md5.New()
	hash.Write([]byte(password))
	hash_password := hex.EncodeToString(hash.Sum(nil))
	return hash_password
}

// 根据用户名查表
func findUser(username string) (bool, error) {

	return false, nil
}

// 注册服务
func RegisterService(username string, password string) (int, error) {
	log.Println(username, "---", password)
	// 先校验参数
	if len(username) > 32 || (len(password) > 32 || len(password) <= 5) {
		return 0, errors.New("参数错误")
	}
	//if !VerifyEmailFormat(username) {
	//	return false, errors.New("邮箱格式错误")
	//}
	// 查表，是否存在id
	// pass, error := model.GetUserByName(username)

	// 插入
	Id, _ := model.InsertUser(username, password)

	return Id, nil
}

// 登录服务
func LoginService(username string, password string) (bool, error) {
	// 先校验参数
	if len(username) > 32 || (len(password) > 32 || len(password) <= 5) {
		return false, errors.New("参数错误")
	}

	return true, nil

}
