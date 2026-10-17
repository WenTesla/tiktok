package service

import (
	"errors"
	"log"
	"regexp"

	"tiktok/go/config"
	"tiktok/go/middle/jwt"
	"tiktok/go/model"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对密码加盐哈希（替代原先的 MD5+静态盐）
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 比对明文密码与 bcrypt 哈希
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

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

// GenerateTokenByName 根据name生成token
func GenerateTokenByName(username string) (string, error) {
	user, _ := model.GetUserByName(username)
	token := jwt.SignToken(user)
	return token, nil
}

// RegisterService 注册服务
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
	// 插入（密码经 bcrypt 哈希后入库）
	hashed, err := HashPassword(password)
	if err != nil {
		return 0, errors.New("密码加密失败")
	}
	user, _ = model.InsertUser(username, hashed)
	return user.Id, nil
}

// LoginService 登录服务
func LoginService(username string, password string) (int64, error) {

	user, err := model.GetUserByName(username)
	if err != nil {
		return 0, dataSourceErr
	}
	if username != user.Name {
		return 0, errors.New("用户名不存在")
	}
	if !CheckPassword(user.Password, password) {
		return 0, errors.New("用户或密码不正确")
	}

	return user.Id, nil
}

//  用户服务 先封装小的，再封装大的 根据用户Id查询用户的具体信息

func UserService(Id int64) (model.UserInfo, error) {
	user, err := model.GetUserById(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 脱密
	user.Password = ""
	// 查询自己的关注数目
	followingCount, _ := model.GetFollowingById(Id)
	// 查询自己的粉丝
	fanCount, err := model.GetFansById(Id)
	// 查询点赞数量
	favoriteCount, err := model.QueryFavoriteCountByUserId(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 查询作品的数量
	workCount, err := model.QueryWorkCountByUserId(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 查询获赞数量
	totalFavorited, err := model.QueryTotalFavorited(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 以下为假数据 -avator

	// 关注一定为false
	userInfo := model.UserInfo{
		Id:             user.Id,
		Name:           user.Name,
		FollowCount:    followingCount,
		FollowerCount:  fanCount,
		IsFollow:       false,
		AvatarUrl:      config.MockAvatarUrl,
		TotalFavorited: totalFavorited,
		WorkCount:      workCount,
		FavoriteCount:  favoriteCount,
	}
	return userInfo, nil
}

// UserServiceBatch 批量封装作者信息（消除 N+1），返回 authorId->UserInfo
func UserServiceBatch(authorIds []int64, viewerId int64) (map[int64]model.UserInfo, error) {
	result := make(map[int64]model.UserInfo, len(authorIds))
	if len(authorIds) == 0 {
		return result, nil
	}
	users, err := model.GetUserByIds(authorIds)
	if err != nil {
		return nil, err
	}
	following, fans, err := model.GetFollowCountsByUserIds(authorIds)
	if err != nil {
		return nil, err
	}
	favorite, totalFavorited, err := model.GetLikeCountsByUserIds(authorIds)
	if err != nil {
		return nil, err
	}
	workCounts, err := model.GetWorkCountsByUserIds(authorIds)
	if err != nil {
		return nil, err
	}
	var isFollow map[int64]bool
	if viewerId != -1 {
		isFollow, err = model.QueryIsFollowBatch(viewerId, authorIds)
		if err != nil {
			return nil, err
		}
	}
	for _, u := range users {
		result[u.Id] = model.UserInfo{
			Id:             u.Id,
			Name:           u.Name,
			FollowCount:    following[u.Id],
			FollowerCount:  fans[u.Id],
			IsFollow:       viewerId != -1 && isFollow[u.Id],
			AvatarUrl:      config.MockAvatarUrl,
			TotalFavorited: totalFavorited[u.Id],
			WorkCount:      workCounts[u.Id],
			FavoriteCount:  favorite[u.Id],
		}
	}
	return result, nil
}

//  用户服务 先封装小的，再封装大的 此时为登录状态，需要查询是否登录

func UserInfoService(Id int64, userId int64) (model.UserInfo, error) {
	user, err := model.GetUserById(Id)
	if err != nil {
		return model.UserInfo{}, err
	}
	user.Password = ""
	// 查询关注数目
	followingCount, _ := model.GetFollowingById(Id)
	// 查询粉丝
	fanCount, err := model.GetFansById(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 查询是否关注
	isFollow, err := model.QueryIsFollow(userId, Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 查询作品数量
	favoriteCount, err := model.QueryFavoriteCountByUserId(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 查询点赞数量
	workCount, err := model.QueryWorkCountByUserId(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	// 查询获赞数量
	totalFavorited, err := model.QueryTotalFavorited(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	userInfo := model.UserInfo{
		Id:             user.Id,
		Name:           user.Name,
		FollowCount:    followingCount,
		FollowerCount:  fanCount,
		IsFollow:       isFollow,
		AvatarUrl:      config.CosUrl,
		TotalFavorited: totalFavorited,
		WorkCount:      workCount,
		FavoriteCount:  favoriteCount,
	}
	return userInfo, nil
}

//  脱密后的信息

func SimpleUserService(Id int64, userId int64) (model.UserInfo, error) {
	user, err := model.GetUserById(Id)
	if err != nil {
		return model.UserInfo{}, dataSourceErr
	}
	user.Password = ""
	// 关注一定为true
	userInfo := model.UserInfo{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return userInfo, nil
}
