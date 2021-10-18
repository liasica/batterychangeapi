package model

import (
    "battery/app/model/internal"
    "github.com/gogf/gf/os/gtime"
)

type Message struct {
    internal.Message
    Detail MessageDetail `json:"detail,omitempty"` // 详情页信息  type 为 100，101， 104 时返回
    IsRead bool          `json:"isRead"`           // 是否已读
}

const (
    MessageTypeSystem = 0 // 系统公告，所有用户

    // 骑手端100开始
    MessageTypeUserBizNewSuccess            = 100 // 用户购买套餐成功
    MessageTypeUserBizPenaltySuccess        = 101 // 用户支付违约金成功
    MessageTypeUserBizRenewalSuccess        = 102 // 用户续约金成功
    MessageTypeUserBizBatteryRenewalSuccess = 103 // 用户换电成功
    MessageTypeUserBizExitSuccess           = 104 // 用户退租成功

    // 商家端500开始
    MessageTypeShopManagerBatteryRenewal = 500 // 骑手扫码换电
)

// MessageDetail 订单详情
type MessageDetail struct {
    Type         string      `json:"type,omitempty"`         // 订单类型 type 为 100，102 时返回
    ComboName    string      `json:"comboName,omitempty"`    // 套餐名称 type 为 100 时返回
    CityName     string      `json:"cityName,omitempty"`     // 可用城市 type 为 100 时返回
    BatteryType  uint        `json:"batteryType,omitempty"`  // 可用电池 type 为 100 时返回
    ComboOrderNo string      `json:"comboOrderNo,omitempty"` // 订单编号 type 为 100，102 时返回
    PayType      uint        `json:"payType,omitempty"`      // 支付方式 1 支付宝 2 微信  type 为 100，102 时返回
    PayAt        *gtime.Time `json:"payAt,omitempty"`        // 支付时间 type 为 100，102 时返回
    Amount       float64     `json:"amount,omitempty"`       // 支付金额 type 为 100，102 时返回
    Deposit      float64     `json:"deposit,omitempty"`      // 押金金额 type 为 100 时返回

    ShopName    string      `json:"shopName,omitempty"`              // 门店名称 type 为 104 时返回
    ExitAt      *gtime.Time `json:"at,omitempty"`                    // 退租时间 type 为 104 时返回
    DepositDesc string      `json:"depositDesc,omitempty,omitempty"` // 押金说明 type 为 104 时返回
}
