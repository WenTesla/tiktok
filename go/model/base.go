package model

import "tiktok/go/config"

// 数据库连接池
var db = config.InitDataSource()

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
