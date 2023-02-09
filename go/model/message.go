package model

import (
	"tiktok/go/config"
	"time"
)

type Message struct {
	ID         int64     `json:"id"`
	UserId     int64     `json:"from_user_id"`
	ToUserId   int64     `json:"to_user_id"`
	Content    string    `json:"content"`
	IsWithdraw int8      `json:"is_withdraw,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"column:createTime"` // 创建时间
}

// InsertMessage 插入数据
func InsertMessage(userId int64, toUserId int64, content string) (bool, error) {
	messageInfo := Message{
		UserId:     userId,
		ToUserId:   toUserId,
		Content:    content,
		CreateTime: time.Now(),
	}
	// INSERT INTO `messages` (`user_id`,`to_user_id`,`content`,`is_withdraw`,`createTime`) VALUES (5,1,'111',0,'2023-02-08 19:21:15.017')
	result := db.Debug().Create(&messageInfo)
	if result.Error != nil {
		return false, result.Error
	}
	return false, nil
}

// QueryMessageByUserId 根据用户Id查询聊天记录
func QueryMessageByUserId(userId int64) ([]Message, error) {
	var messages []Message
	//SELECT * FROM `messages` WHERE user_id = 1 AND is_withdraw = 0 LIMIT 10
	result := db.Where("user_id = ? AND is_withdraw = ?", userId, 0).Limit(config.MessageCount).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

// QueryNewestMessageByUserId 通过用户Id查询最新的聊天记录 0-接受 1-发送 有点问题
func QueryNewestMessageByUserId(userId int64) (string, int8, error) {
	message := Message{}
	// SELECT * FROM `messages` WHERE (user_id = 1 Or to_user_id = 1) AND is_withdraw = 0 ORDER BY createTime desc LIMIT 1
	result := db.Debug().Where("(user_id = ? Or to_user_id = ?) AND is_withdraw = ?", userId, userId, 0).Order("createTime desc").Limit(1).Find(&message)
	if result.Error != nil {
		return "", -1, result.Error
	}
	if userId == message.UserId {
		return message.Content, 1, nil
	} else {
		return message.Content, 0, nil
	}
}

// QueryNewestMessageByUserIdAndToUserID 通过两者的用户Id查询最新最新的两者之间的聊天记录 0-接受 1-发送 有点问题
func QueryNewestMessageByUserIdAndToUserID(userId int64, toUserId int64) (string, int8, error) {
	message := Message{}
	// SELECT `content`,`createTime`,`user_id`,`to_user_id` FROM `messages` WHERE (user_id = 2 AND to_user_id = 7 AND is_withdraw = 0) OR (user_id = 7 AND to_user_id = 2 AND is_withdraw = 0) ORDER BY createTime desc LIMIT 1
	result := db.Debug().Where("user_id = ? AND to_user_id = ? AND is_withdraw = ?", userId, toUserId, 0).Or("user_id = ? AND to_user_id = ? AND is_withdraw = ?", toUserId, userId, 0).Order("createTime desc").Limit(1).Select("content", "createTime", "user_id", "to_user_id").Find(&message)
	if result.Error != nil {
		return "", -1, result.Error
	}
	if userId == message.UserId {
		return message.Content, 1, nil
	} else {
		return message.Content, 0, nil
	}
}

// QueryMessageByUserIdAndToUserId 查询两者的全部聊天记录
func QueryMessageByUserIdAndToUserId(userId int64, toUserId int64) ([]Message, error) {
	var messages []Message
	// SELECT * FROM `messages` WHERE (user_id = 1 AND to_user_id = 2 AND is_withdraw = 0) OR (user_id = 2 AND to_user_id = 1 AND is_withdraw = 0) ORDER BY createTime desc
	result := db.Debug().Where("user_id = ? AND to_user_id = ? AND is_withdraw = ?", userId, toUserId, 0).Or("user_id = ? AND to_user_id = ? AND is_withdraw = ?", toUserId, userId, 0).Order("createTime desc").Find(&messages)
	if result.Error != nil {
		return messages, result.Error
	}
	//for _, message := range messages {
	//	message.CreateTime = message.CreateTime.Format("2006-01-02 15:04:05")
	//}
	return messages, nil
}
