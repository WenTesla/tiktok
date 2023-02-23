package model

import (
	"errors"
	"tiktok/go/util"
	"time"
)

type Comment struct {
	ID         int64     // 评论id
	UserId     int64     // 用户Id
	VideoId    int64     //视频Id
	Text       string    // 评论内容
	IsCancel   int64     // 是否取消
	CreateTime time.Time `gorm:"column:createTime"` // 创建时间
}

// CommentInfo 评论详细信息
type CommentInfo struct {
	Content    string   `json:"content"`     // 评论内容
	CreateDate string   `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64    `json:"id"`          // 评论id
	UserInfo   UserInfo `json:"user"`        // 评论用户的具体信息
}

// QueryCommentByUserId 根据用户id获取评论
func QueryCommentByUserId(userId int64) ([]Comment, error) {
	var comments []Comment
	// 获取
	result := db.Debug().Where("user_id = ? AND is_cancel = ?", userId, 0).Order("createTime desc").Find(&comments)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return nil, result.Error
	}
	return comments, nil
}

// QueryCommentByVideoId 根据视频id获取评论
func QueryCommentByVideoId(videoId int64) ([]Comment, error) {
	var comments []Comment
	// 获取
	// SELECT * FROM `comments` WHERE video_id = 43 AND is_cancel = 0 ORDER BY createTime desc
	result := db.Where("video_id = ? AND is_cancel = ?", videoId, 0).Order("createTime desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil

}

// QueryCommentCountByVideoId 根据视频的id获取视频的评论数
func QueryCommentCountByVideoId(videoId int64) (int64, error) {
	var count int64
	//  SELECT count(*) FROM `comments` WHERE video_id = 43 AND is_cancel = 0
	result := db.Model(&Comment{}).Where("video_id = ? AND is_cancel = ?", videoId, 0).Count(&count)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return -1, result.Error
	}
	return count, nil
}

// InsertComment 插入评论 这有问题 传入的参数
func InsertComment(userId int64, videoId int64, content string) (Comment, error) {
	comment := Comment{
		UserId:     userId,
		VideoId:    videoId,
		Text:       content,
		CreateTime: time.Now(),
	}
	result := db.Debug().Select("user_id", "video_id", "text").Create(&comment)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return comment, result.Error
	}
	return comment, nil
}

// DeleteComment 删除评论
func DeleteComment(id int64) (bool, error) {
	result := db.Debug().Delete(&Comment{}, id)
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, errors.New("该评论不存在")
	}
	return true, nil
}

//  取消评论 gorm有bug还是设计问题，不能通过主键直接更新,必须先查询数据

func CancelComment(id int64, userId int64, videoId int64) (bool, error) {
	comment := Comment{
		ID: id,
	}
	//UPDATE `comments` SET `is_cancel`=1 WHERE (user_id = 1 AND video_id = 15) AND `id` = 48
	result := db.Debug().Model(&comment).Where("user_id = ? AND video_id = ?", userId, videoId).Update("is_cancel", 1)
	//// SELECT * FROM `comments` WHERE `comments`.`id` = 10 ORDER BY `comments`.`id` LIMIT 1
	//db.First(&comment, id)
	//comment.IsCancel = 1
	//result := db.Debug().Save(&comment)
	//if result.Error != nil {
	//	util.LogError(result.Error.Error())
	//	return false, result.Error
	//}
	//if result.RowsAffected == 0 {
	//	return false, errors.New("该评论不存在")
	//}
	if result.Error != nil {
		util.LogError(result.Error.Error())
		return false, result.Error
	}
	return true, nil
}
