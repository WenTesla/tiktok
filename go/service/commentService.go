package service

import "tiktok/go/model"

func CommentListService(userId int64, videoId int64) ([]model.CommentInfo, error) {

	comments, err := model.QueryCommentByUserId(userId)
	if err != nil {
		return nil, err
	}
	PackageComments(comments)

	return nil, nil
}

func PackageComments(comments []model.Comment) ([]model.CommentInfo, error) {

	for _, comment := range comments {
		PackageComment(comment)
	}
	return nil, nil
}

func PackageComment(comment model.Comment) (model.CommentInfo, error) {
	return model.CommentInfo{}, nil
}
