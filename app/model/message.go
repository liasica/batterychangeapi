package model

import "battery/app/model/internal"

type Message struct {
	internal.Message
	Detail interface{} `json:"detail"`
	IsRead bool        `json:"isRead"` //是否已读
}

const (
	MessageTypeSystem = 0 //系统公告，所有用户

	//骑手端100开始
	MessageTypeUserBizNewSuccess            = 100 //用户购买套餐成功
	MessageTypeUserBizPenaltySuccess        = 101 //用户支付违约金成功
	MessageTypeUserBizRenewalSuccess        = 102 //用户续约金成功
	MessageTypeUserBizBatteryRenewalSuccess = 103 //用户换电成功
	MessageTypeUserBizExitSuccess           = 104 //用户退租成功

	//商家端500开始
	MessageTypeShopManagerBatteryRenewal = 500 //骑手扫码换电
)

// MessageTypeUserBizNewSuccessDetail 骑手购买套餐成功消息详情
type MessageTypeUserBizNewSuccessDetail struct {
	PackagesName   string  `json:"packagesName"`   //套餐名称
	Type           string  `json:"type"`           //交易商品
	PackageOrderNo string  `json:"PackageOrderNo"` //订单编号
	PayType        uint    `json:"payType"`        //支付方式 1 支付宝 2 微信
	PayAt          uint    `json:"payAt"`          //支付时间
	Amount         float64 `json:"amount"`         //支付金额
}
