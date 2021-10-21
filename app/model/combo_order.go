package model

import (
    "battery/app/model/internal"
    "github.com/gogf/gf/os/gtime"
    "github.com/gogf/gf/util/gmeta"
)

type ComboOrder internal.ComboOrder

const PayTypeAliPay = 1
const PayTypeWechat = 2

const ComboTypeNew = 1
const ComboTypeRenewal = 2
const ComboTypePenalty = 3

const PayStateWait = 1    // 待支付
const PayStateSuccess = 2 // 已支付

// ShopOrderListReq 门店订单列表
type ShopOrderListReq struct {
    Page
    Keywords string `json:"keywords"` // 搜索关键字
    Type     uint   `json:"type"`     // 1 新签， 2 续费
    Month    uint   `json:"month"`    // 月份
}

// ShopOrderListItem 门店订单列表
type ShopOrderListItem struct {
    Id         uint64      `json:"id"`         // 订单ID
    ComboName  string      `json:"comboName"`  // 套餐名称
    Amount     float64     `json:"amount"`     // 金额
    UserName   string      `json:"userName"`   // 用户姓名
    UserMobile string      `json:"userMobile"` // 用户电话
    OrderNo    string      `json:"orderNo"`    // 订单编号
    Type       uint        `json:"type"`       // 1 新签， 2 续费
    PayAt      *gtime.Time `json:"payAt"`      // 支付时间
}

type ShopOrderTotalReq struct {
    Type  uint `json:"type"`                                   // 1 新签， 2 续费
    Month uint `validate:"required" v:"required" json:"month"` // 月份
}

// ShopOrderTotalRep 门店订单列表
type ShopOrderTotalRep struct {
    Cnt    int64   `json:"cnt"`    // 总记录数
    Amount float64 `json:"amount"` // 总额
}

// UserCurrentComboOrder 个签骑手当前套餐信息
type UserCurrentComboOrder struct {
    ComboName    string      `validate:"required" json:"comboName"`    // 套餐名称
    ComboAmount  float64     `validate:"required" json:"comboAmount"`  // 套餐金额
    ExpirationAt *gtime.Time `validate:"required" json:"expirationAt"` // 截止日期
    CityName     string      `validate:"required" json:"cityName"`     // 可用城市
    OrderNo      string      `validate:"required" json:"orderNo"`      // 订单编号
    PayType      uint        `validate:"required" json:"payType"`      // 支付方式 1 支付宝 2 微信
    Amount       float64     `validate:"required" json:"amount"`       // 支付金额
    Deposit      float64     `validate:"required" json:"deposit"`      // 押金金额 0 为没有押金
    PayAt        *gtime.Time `validate:"required" json:"payAt"`        // 支付时间
    StartUseAt   *gtime.Time `validate:"required" json:"startUseAt"`   // 开通时间
    BatteryState uint        `validate:"required" json:"batteryState"` // 个签骑手换电状态：0 未开通， 1 新签未领 ，2 租借中，3 寄存中，4 已退租， 5 已逾期
}

// ShopManagerComboOrderScanDetailRep 门店认领订单，获取订单信息响应数据
type ShopManagerComboOrderScanDetailRep struct {
    UserType    uint    `validate:"required" json:"userType"`              // 1 个签 2 团签 3 团签BOOS
    UserName    string  `validate:"required" json:"userName"`              // 客户名称
    UserMobile  string  `validate:"required" json:"userMobile"`            // 客户手机号
    ComboName   string  `validate:"required" json:"comboName,omitempty"`   // 套餐名称，个签用户返回
    ComboAmount float64 `validate:"required" json:"comboAmount,omitempty"` // 套餐价格，个签用户返回
    BatteryType string  `validate:"required" json:"batteryType"`           // 电池型号 60 / 72

    GroupName string `validate:"required" json:"groupName,omitempty"` // 团签公司名称，团签签用户返回

    OrderNo string      `validate:"required" json:"orderNo,omitempty"` // 订单编号，个签用户返回
    Amount  float64     `validate:"required" json:"amount,omitempty"`  // 支付金额，个签用户返回
    Deposit float64     `validate:"required" json:"deposit"`           // 押金金额，个签用户返回
    PayType uint        `validate:"required" json:"payType,omitempty"` // 支付方式 1 支付宝 2 微信，个签用户返回
    PayAt   *gtime.Time `validate:"required" json:"payAt,omitempty"`   // 支付时间，个签用户返回

    ClaimState uint `validate:"required" json:"claimState"` // 1 未被认领， 2 已被认领
}

// ShopManagerComboOrderListDetailRep 订单记录获取订单详情响应数据
type ShopManagerComboOrderListDetailRep struct {
    UserName    string `validate:"required" json:"userName"`    // 客户名称
    UserMobile  string `validate:"required" json:"userMobile"`  // 客户手机号
    ComboName   string `validate:"required" json:"comboName"`   // 套餐名称
    BatteryType string `validate:"required" json:"batteryType"` // 电池型号 60 / 72

    OrderNo string      `validate:"required" json:"orderNo"` // 订单编号
    Amount  float64     `validate:"required" json:"amount"`  // 支付金额
    Deposit float64     `validate:"required" json:"deposit"` // 押金金额
    PayType uint        `validate:"required" json:"payType"` // 支付方式 1 支付宝 2 微信
    PayAt   *gtime.Time `validate:"required" json:"payAt"`   // 支付时间
}

// ShopManagerComboOrderClaimReq 门店认领订单，请求数据
type ShopManagerComboOrderClaimReq struct {
    Code string `validate:"required" json:"code"` // 扫码获取
}

// Prepay 发起支付数据
type Prepay struct {
    Description string
    No          string
    Amount      float64
    NotifyUrl   string
}

// BizNewCdeReq 门店认领订单扫码获取订单请求数据
type BizNewCdeReq struct {
    Code string `validate:"required" v:"required" json:"code"` // code二维码
}

// OrderListReq 订单列表请求项
type OrderListReq struct {
    Page
    No        string      `json:"no"`                     // 订单编号
    ShopId    uint        `json:"shopId"`                 // 门店ID
    UserId    uint        `json:"userId"`                 // 骑手ID
    ComboId   uint        `json:"comboId"`                // 套餐ID
    CityId    uint        `json:"cityId"`                 // 城市ID
    RealName  string      `json:"realName"`               // 骑手姓名
    Mobile    string      `json:"mobile" v:"phone-loose"` // 手机号
    Type      uint        `json:"type"`                   // 订单类型: 1新签 2续签 3违约金
    StartDate *gtime.Time `json:"startDate"`              // 开始日期 eg: 2021-10-17
    EndDate   *gtime.Time `json:"endDate"`                // 结束日期 eg: 2021-10-19
}

// OrderEntity 订单实体
type OrderEntity struct {
    gmeta.Meta `orm:"table:combo_order" swaggerignore:"true"`

    Id         uint        `json:"id"`
    No         uint        `json:"no"`
    ShopId     uint        `json:"shopId"`
    UserId     uint        `json:"userId"`
    ComboId    uint        `json:"comboId"`
    CityId     uint        `json:"cityId"`
    Type       uint        `json:"type" enums:"1,2,3"`    // 订单类型: 1新签 2续签 3违约金
    Amount     float64     `json:"amount"`                // 支付金额
    Deposit    float64     `json:"deposit"`               // 押金
    PayType    uint        `json:"payType" enums:"0,1,2"` // 支付方式: 0未知 1支付宝 2微信
    PayState   uint        `json:"payState" enums:"1,2"`  // 支付状态: 1未支付 2已支付
    PayAt      *gtime.Time `json:"payAt"`                 // 支付时间
    FirstUseAt *gtime.Time `json:"firstUseAt"`            // 开始时间
    CreatedAt  *gtime.Time

    User        *User      `orm:"with:id=userId"`
    City        *Districts `orm:"with:id=cityId"`
    Shop        *Shop      `orm:"with:id=shopId"`
    ComboDetail *Combo     `orm:"with:id=comboId"`
}

// OrderListItem 订单列表返回
type OrderListItem struct {
    Id         uint        `json:"id"`
    No         uint        `json:"no"`                              // 订单编号
    ShopId     uint        `json:"shopId"`                          // 门店ID
    UserId     uint        `json:"userId"`                          // 骑手ID
    RealName   string      `json:"realName"`                        // 骑手姓名
    Mobile     string      `json:"mobile" v:"required|phone-loose"` // 手机号
    Type       uint        `json:"type" enums:"1,2,3"`              // 订单类型: 1新签 2续签 3违约金
    ComboName  string      `json:"comboName"`                       // 套餐名称
    ShopName   string      `json:"shopName"`                        // 门店名
    CityName   string      `json:"cityName"`                        // 城市
    Amount     float64     `json:"amount"`                          // 支付金额
    Deposit    float64     `json:"deposit"`                         // 押金
    PayType    uint        `json:"payType" enums:"0,1,2"`           // 支付方式: 0未知 1支付宝 2微信
    PayState   uint        `json:"payState" enums:"1,2"`            // 支付状态: 1未支付 2已支付
    PayAt      *gtime.Time `json:"payAt"`                           // 支付时间
    FirstUseAt *gtime.Time `json:"firstUseAt"`                      // 开始时间
    CreatedAt  *gtime.Time `json:"createdAt"`                       // 订单时间
}
