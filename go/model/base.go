package model

import (
	"gorm.io/gorm"

	"tiktok/go/config"
)

// 数据库连接池，由 main 启动时通过 InitDB 注入；单测可注入 go-sqlmock。
var db *gorm.DB

// InitDB 注入数据库实例（启动时调用；单测中可注入 mock）
func InitDB(g *gorm.DB) {
	db = g
}

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var BaseResponseInstance = BaseResponse{}

func (baseResponse *BaseResponse) Success() BaseResponse {
	baseResponse.StatusCode = 0
	baseResponse.StatusMsg = config.Success
	return BaseResponseInstance
}

func (baseResponse *BaseResponse) Fail() BaseResponse {
	baseResponse.StatusCode = -1
	baseResponse.StatusMsg = config.Fail
	return BaseResponseInstance
}

func (baseResponse *BaseResponse) SuccessMsg(msg string) BaseResponse {
	baseResponse.StatusCode = 0
	baseResponse.StatusMsg = msg
	return BaseResponseInstance
}
func (baseResponse *BaseResponse) FailMsg(msg string) BaseResponse {
	baseResponse.StatusCode = -1
	baseResponse.StatusMsg = msg
	return BaseResponseInstance
}
