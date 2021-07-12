package model

import (
	"battery/app/model/internal"
)

type ShopManager internal.ShopManager

// ShopManagerLoginReq 店长登录请求数据
type ShopManagerLoginReq struct {
	Mobile string `validate:"required" v:"required|phone-loose"` //手机号
	Sms    string `validate:"required" v:"required|length:6,6"`  //短信验证码
}

// ShopManagerLoginRep 店长登录返回数据
type ShopManagerLoginRep struct {
	AccessToken string // 请求 token
}

// ShopManagerShopStateReq 店长登录返回数据
type ShopManagerShopStateReq struct {
	State string // 请求 token
}

// ShopManagerResetMobileReq 切换手机号请求
type ShopManagerResetMobileReq struct {
	Mobile string `validate:"required" v:"required|length:11|phone-loose"` //手机号
	Sms    string `validate:"required" v:"required|length:6,16"`           //短信验证码
}
