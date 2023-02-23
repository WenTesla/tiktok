package service

import (
	"errors"
	"tiktok/go/model"
	"tiktok/go/util"
)

func CommentListService(videoId int64) ([]model.CommentInfo, error) {
	comments, err := model.QueryCommentByVideoId(videoId)
	if err != nil {
		return nil, dataSourceErr
	}
	commentInfos, err := PackageComments(comments)
	return commentInfos, nil
}

//  包装切片

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

//  包装一个结构体 传入userInfo值

func PackageComment(comment model.Comment) (model.CommentInfo, error) {
	// 定义
	commentInfo := model.CommentInfo{}
	// 填入Id
	commentInfo.ID = comment.ID
	// 填入评论内容
	commentInfo.Content = comment.Text
	// 填入创建时间 "2006-01-02 15:04:05.999999999 -0700 MST"
	commentInfo.CreateDate = comment.CreateTime.Format("01-02")
	// 根据用户id查询
	user, err := model.GetUserById(comment.UserId)
	if err != nil {
		return commentInfo, dataSourceErr
	}
	userInfo, err := model.PackageUserToUserInfo(user)
	if err != nil {
		return commentInfo, err
	}
	commentInfo.UserInfo = userInfo
	return commentInfo, nil
}

// 创建评论

func CreateCommentService(userId int64, videoId int64, content string) (model.CommentInfo, error) {
	// 敏感词处理
	content, err := replaceSensitive(content)
	commentInfo := model.CommentInfo{}
	isExistVideoId, err := model.QueryIsExistVideoId(videoId)
	if isExistVideoId != true {
		return commentInfo, errors.New("视频不存在")
	}
	// 插入评论
	comment, err := model.InsertComment(userId, videoId, content)
	if err != nil {
		return commentInfo, dataSourceErr
	}
	commentInfo, err = PackageComment(comment)
	if err != nil {
		return commentInfo, dataSourceErr
	}
	return commentInfo, nil
}

//  删除评论

func CancelCommentService(id int64, userId int64, videoId int64) (bool, error) {

	isDelete, err := model.CancelComment(id, userId, videoId)
	if err != nil {
		return false, dataSourceErr
	}
	return isDelete, nil
}

// checkSensitive 检验敏感词
func checkSensitive(content string) bool {

	return false
}

//  替换敏感词 优化-》传递指针

func replaceSensitive(content string) (string, error) {
	content, err := util.SensitiveWordsFilter(content)
	return content, err
}
