package model

import "tiktok/go/config"

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// 响应状态
func (BaseResponse) Success() BaseResponse {
	return BaseResponse{
		StatusCode: 0,
		StatusMsg:  config.Success,
	}
}

func (baseResponse *BaseResponse) Fail() BaseResponse {
	return BaseResponse{
		StatusCode: -1,
		StatusMsg:  config.Fail,
	}
}

func (BaseResponse) SuccessMsg(msg string) BaseResponse {
	return BaseResponse{
		StatusCode: 0,
		StatusMsg:  msg,
	}
}
func (BaseResponse) FailMsg(msg string) BaseResponse {
	return BaseResponse{
		StatusCode: -1,
		StatusMsg:  msg,
	}
}
