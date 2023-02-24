package service

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
	"tiktok/go/config"
	"tiktok/go/model"
	"tiktok/go/util"
	"time"
)

// 定义键值对维护消息记录 用户的Id->用户目前的消息记录索引
// var userCommentIndex = make(map[int64]int64)
var userCommentIndex sync.Map

// 定义键值对维护消息记录 用户的Id->用户目前的消息记录最大值
// var userMessageMaxIndex = make(map[int64]int64)
var userMessageMaxIndex sync.Map

var NowTime string

func FriendListService(userId int64) ([]model.FriendUser, error) {
	var FriendUsers []model.FriendUser
	var err error
	FriendUsers, err = PackageFriendLists(userId)
	if err != nil {
		return nil, err
	}
	return FriendUsers, nil
}

func MessageService() {

}

//  通过userId查询粉丝的数据，再包装加入消息

func PackageFriendLists(userId int64) ([]model.FriendUser, error) {
	var FriendLists []model.FriendUser
	var message string
	var msgType int8
	userInfos, err := MutualFollowListService(userId)
	if err != nil {
		return nil, err
	}
	for _, userInfo := range userInfos {
		// 查询Message和MsgType
		message, msgType, err = model.QueryNewestMessageByUserIdAndToUserID(userId, userInfo.Id)
		if err != nil {
			return nil, err
		}
		userInfo.AvatarUrl = config.MockAvatarUrl
		FriendLists = append(FriendLists, model.FriendUser{
			UserInfo: userInfo,
			Message:  message,
			MsgType:  msgType,
		})
	}

	return FriendLists, nil
}

//  包装单个请求

func PackageFriendList(userInfo model.FriendUser) (model.UserInfo, error) {
	// test
	userInfo.AvatarUrl = "https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/20230206171653.png"
	return model.UserInfo{}, nil
}

func MessageChatService(userId int64, toUserId int64, preMsgTime int64) ([]model.Message, error) {
	var messages []model.Message
	var err error
	// 第一次查询
	if preMsgTime == 0 {
		// 第一次使用
		// 查询userid和toUserId的表将全部内容返回
		messages, err = model.QueryMessageByUserIdAndToUserId(userId, toUserId)
		if err != nil {
			return nil, err
		}
		// 添加所有数据进入redis
		//AllMessageListAddRedis(userId, toUserId, messages)
		return messages, err
	}
	// 不是第一次查询 查询redis
	messages, err = ParseAllMessageListFromRedis(userId, toUserId, preMsgTime)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func MessageActionService(userId int64, toUserId int64, content string) (bool, error) {
	// 敏感词替换
	content, _ = util.SensitiveWordsFilter(content)

	// 添加数据进入redis
	MessageActionRedis(userId, toUserId, content)
	// 添加数据库
	pass, err := model.InsertMessage(userId, toUserId, content)
	if err != nil {
		return false, dataSourceErr
	}
	if !pass {
		return false, errors.New("发送消息失败")
	}
	return true, nil
}

func MessageActionRedis(userId int64, toUserId int64, content string) error {
	timeUnix := time.Now().Unix()
	message := model.Message{
		UserId:     userId,
		ToUserId:   toUserId,
		Content:    content,
		IsWithdraw: 0,
		CreateTime: timeUnix,
	}
	// 序列化
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	var messageRedisName = strconv.FormatInt(userId, 10) + "-" + strconv.FormatInt(toUserId, 10)
	// 添加到redis中hash
	messageRedisDb.RPush(messageRedisName, bytes)
	return nil
}

func AllMessageListAddRedis(userId int64, toUserId int64, messages []model.Message) {
	var messageRedisName = strconv.FormatInt(userId, 10) + "-" + strconv.FormatInt(toUserId, 10)
	// 添加到redis中hash
	for _, message := range messages {
		bytes, _ := json.Marshal(message)
		messageRedisDb.ZAdd(messageRedisName, redis.Z{
			Score:  float64(message.CreateTime),
			Member: bytes,
		})
	}
}

// 查询时间戳之前的记录 并删除

func ParseAllMessageListFromRedis(userId int64, toUserId int64, msgTime int64) ([]model.Message, error) {
	var messages []model.Message
	var messageRedisName = strconv.FormatInt(userId, 10) + "-" + strconv.FormatInt(toUserId, 10)
	for {
		bytes, _ := messageRedisDb.LPop(messageRedisName).Result()
		if bytes == "" {
			break
		}
		message := model.Message{}
		json.Unmarshal([]byte(bytes), &message)
		if message.IsWithdraw != 0 || message.CreateTime >= msgTime {
			continue
		}
		messages = append(messages, message)
	}
	return messages, nil
}
