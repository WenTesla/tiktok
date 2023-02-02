package service

import (
	"tiktok/go/model"
)

func CommentListService(userId int64, videoId int64) ([]model.CommentInfo, error) {
	comments, err := model.QueryCommentByVideoId(videoId)
	if err != nil {
		return nil, err
	}
	commentInfos, err := PackageComments(comments)

	return commentInfos, nil
}

// PackageComments 包装切片
func PackageComments(comments []model.Comment) ([]model.CommentInfo, error) {
	// 提前定义切片
	var commentInfos []model.CommentInfo
	for _, comment := range comments {
		commentInfo, err := PackageComment(comment)
		if err != nil {
			return nil, err
		}
		commentInfos = append(commentInfos, commentInfo)
	}
	return commentInfos, nil
}

// PackageComment 包装一个结构体 传入userInfo值
func PackageComment(comment model.Comment) (model.CommentInfo, error) {
	// 定义
	commentInfo := model.CommentInfo{}
	// 填入Id
	commentInfo.ID = comment.ID
	// 填入评论内容
	commentInfo.Content = comment.Text
	// 填入创建时间 "2006-01-02 15:04:05.999999999 -0700 MST"
	commentInfo.CreateDate = comment.CreateTime.Format("2006-01-02 15:04:05")
	// 根据用户id查询
	user, err := model.GetUserById(comment.UserId)
	if err != nil {
		return commentInfo, err
	}
	userInfo, err := model.PackageUserToUserInfo(user)
	if err != nil {
		return commentInfo, err
	}
	commentInfo.UserInfo = userInfo
	return commentInfo, nil
}

// CreateCommentService 创建评论
func CreateCommentService(userId int64, videoId int64, content string) (model.CommentInfo, error) {
	// 检验敏感词

	commentInfo := model.CommentInfo{}
	// 插入评论
	comment, err := model.InsertComment(userId, videoId, content)
	if err != nil {
		return commentInfo, err
	}
	commentInfo, err = PackageComment(comment)
	if err != nil {
		return commentInfo, err
	}
	return commentInfo, nil
}

// DeleteCommentService 删除评论
func DeleteCommentService(id int64) (bool, error) {
	isDelete, err := model.DeleteComment(id)
	if err != nil {
		return false, err
	}
	return isDelete, nil
}

// checkSensitive 检验敏感词
func checkSensitive(context string) bool {

	return false
}
