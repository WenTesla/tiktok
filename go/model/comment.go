package model

import "tiktok/go/config"

type Comment struct {
	ID       int64  // 评论id
	UserId   int64  // 用户Id
	VideoId  int64  //视频Id
	Text     string // 评论内容
	IsCancel int64  // 是否取消
}

// CommentInfo
type CommentInfo struct {
	Content    string   `json:"content"`     // 评论内容
	CreateDate string   `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64    `json:"id"`          // 评论id
	User       UserInfo `json:"user"`        // 评论用户信息
}

// type
// 根据用户id获取评论
func QueryCommentByUserId(userId int64) ([]Comment, error) {
	comments := make([]Comment, config.CommentCount)
	// 获取
	result := db.Debug().Where("user_id = ?", userId).Order("createTime desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
