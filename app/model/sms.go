package model

import (
	"battery/app/model/internal"
)

type Sms internal.Sms

// SmsSendReq 获取手机验证码请求数据
type SmsSendReq struct {
	Mobile    string `validate:"required" v:"required|phone-loose"` //手机号码
	Sig       string `validate:"required" v:"required"`
	SessionId string `validate:"required" v:"required"`
	Token     string `validate:"required" v:"required"`
}

// SmsVerifyReq 手机验证码验证请求数据
type SmsVerifyReq struct {
	Mobile string `validate:"required" v:"required|phone-loose"`
	Code   string `validate:"required" v:"required|length:6"`
}
