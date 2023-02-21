package model

import (
	"tiktok/go/config"
	"tiktok/go/util"
	"time"
)

type Message struct {
	ID         int64  `json:"id"`
	UserId     int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	IsWithdraw int8   `json:"is_withdraw,omitempty"`
	CreateTime int64  `json:"create_time" gorm:"column:createTime"` // 创建时间
}

//  插入数据

func InsertMessage(userId int64, toUserId int64, content string) (bool, error) {
	messageInfo := Message{
		UserId:     userId,
		ToUserId:   toUserId,
		Content:    content,
		CreateTime: time.Now().Unix(),
	}
	// INSERT INTO `messages` (`user_id`,`to_user_id`,`content`,`is_withdraw`,`createTime`) VALUES (5,1,'111',0,'2023-02-08 19:21:15.017')
	result := db.Create(&messageInfo)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return false, result.Error
	}
	return true, nil
}

//  根据用户Id查询聊天记录

func QueryMessageByUserId(userId int64) ([]Message, error) {
	var messages []Message
	//SELECT * FROM `messages` WHERE user_id = 1 AND is_withdraw = 0 LIMIT 10
	result := db.Where("user_id = ? AND is_withdraw = ?", userId, 0).Limit(config.MessageCount).Find(&messages)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return nil, result.Error
	}
	return messages, nil
}

//  通过用户Id查询最新的聊天记录 0-接受 1-发送 有点问题

func QueryNewestMessageByUserId(userId int64) (string, int8, error) {
	message := Message{}
	// SELECT * FROM `messages` WHERE (user_id = 1 Or to_user_id = 1) AND is_withdraw = 0 ORDER BY createTime desc LIMIT 1
	result := db.Where("(user_id = ? Or to_user_id = ?) AND is_withdraw = ?", userId, userId, 0).Order("createTime desc").Limit(1).Find(&message)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return "", -1, result.Error
	}
	if userId == message.UserId {
		return message.Content, 1, nil
	} else {
		return message.Content, 0, nil
	}
}

//  通过两者的用户Id查询最新最新的两者之间的聊天记录 0-接受 1-发送

func QueryNewestMessageByUserIdAndToUserID(userId int64, toUserId int64) (string, int8, error) {
	message := Message{}
	// SELECT `content`,`createTime`,`user_id`,`to_user_id` FROM `messages` WHERE (user_id = 2 AND to_user_id = 7 AND is_withdraw = 0) OR (user_id = 7 AND to_user_id = 2 AND is_withdraw = 0) ORDER BY createTime desc LIMIT 1
	result := db.Where("user_id = ? AND to_user_id = ? AND is_withdraw = ?", userId, toUserId, 0).Or("user_id = ? AND to_user_id = ? AND is_withdraw = ?", toUserId, userId, 0).Order("createTime desc").Limit(1).Select("content", "createTime", "user_id", "to_user_id").Find(&message)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return "", -1, result.Error
	}
	if userId == message.UserId {
		return message.Content, 1, nil
	} else {
		return message.Content, 0, nil
	}
}

//  查询两者的全部聊天记录

func QueryMessageByUserIdAndToUserId(userId int64, toUserId int64) ([]Message, error) {
	var messages []Message
	// SELECT * FROM `messages` WHERE (user_id = 1 AND to_user_id = 2 AND is_withdraw = 0) OR (user_id = 2 AND to_user_id = 1 AND is_withdraw = 0) ORDER BY createTime asc
	result := db.Where("user_id = ? AND to_user_id = ? AND is_withdraw = ?", userId, toUserId, 0).Or("user_id = ? AND to_user_id = ? AND is_withdraw = ?", toUserId, userId, 0).Order("createTime asc").Find(&messages)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return messages, result.Error
	}
	//for _, message := range messages {
	//	message.CreateTime = message.CreateTime.Unix()
	//}
	return messages, nil
}

//  查询消息记录的最大值

func QueryMessageMaxCount(userId int64, toUserId int64) (int64, error) {
	var count int64
	// SELECT count(*) FROM `messages` WHERE (user_id = 1 AND to_user_id = 2 ) OR (user_id = 2 AND to_user_id = 1 )
	result := db.Model(&Message{}).Where("user_id = ? AND to_user_id = ? ", userId, toUserId).Or("user_id = ? AND to_user_id = ? ", toUserId, userId).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return -1, result.Error
	}
	return count, nil
}
