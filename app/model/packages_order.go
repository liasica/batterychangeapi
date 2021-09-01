package model

import (
	"github.com/gogf/gf/os/gtime"

	"battery/app/model/internal"
)

type PackagesOrder internal.PackagesOrder

const PayTypeAliPay = 1
const PayTypeWechat = 2

const PackageTypeNew = 1
const PackageTypeRenewal = 2
const PackageTypePenalty = 3

const PayStateWait = 1    //待支付
const PayStateSuccess = 2 // 已支付

//ShopOrderListReq 店长订单列表
type ShopOrderListReq struct {
	Page
	Keywords string `json:"keywords"` // 搜索关键字
	Type     uint   `json:"type"`     // 1 新签， 2 续费
	Month    uint   `json:"month"`    // 月份
}

//ShopOrderListItem 店长订单列表
type ShopOrderListItem struct {
	Id          uint64      `json:"id"`          //订单ID
	PackageName string      `json:"packageName"` //套餐名称
	Amount      float64     `json:"amount"`      //金额
	UserName    string      `json:"userName"`    //用户姓名
	UserMobile  string      `json:"userMobile"`  //用户电话
	OrderNo     string      `json:"orderNo"`     //订单编号
	Type        uint        `json:"type"`        // 1 新签， 2 续费
	PayAt       *gtime.Time `json:"shopId"`      //支付时间
}

type ShopOrderTotalReq struct {
	Type  uint `json:"type"`                                   // 1 新签， 2 续费
	Month uint `validate:"required" v:"required" json:"month"` //月份
}

//ShopOrderTotalRep 店长订单列表
type ShopOrderTotalRep struct {
	Cnt    int64   `json:"cnt"`    //总记录数
	Amount float64 `json:"amount"` //总额
}

// UserCurrentPackageOrder 个签骑手当前套餐信息
type UserCurrentPackageOrder struct {
	PackageName  string      `validate:"required" json:"packageName"`  //套餐名称
	ExpirationAt *gtime.Time `validate:"required" json:"expirationAt"` //截止日期
	CityName     string      `validate:"required" json:"cityName"`     //可用城市
	OrderNo      string      `validate:"required" json:"orderNo"`      //订单编号
	PayType      uint        `validate:"required" json:"payType"`      //支付方式 1 支付宝 2 微信
	Amount       float64     `validate:"required" json:"amount"`       //支付金额
	Earnest      float64     `validate:"required" json:"earnest"`      //押金金额 0 为没有押金
	PayAt        *gtime.Time `validate:"required" json:"payAt"`        //支付时间
	StartUseAt   *gtime.Time `validate:"required" json:"startUseAt"`   //开通时间
}

// ShopManagerPackagesOrderScanDetailRep 店长认领订单，获取订单信息响应数据
type ShopManagerPackagesOrderScanDetailRep struct {
	UserType       uint    `validate:"required" json:"userType"`                 // 1 个签 2 团签 3 团签BOOS
	UserName       string  `validate:"required" json:"userName"`                 //客户名称
	UserMobile     string  `validate:"required" json:"userMobile"`               //客户手机号
	PackagesName   string  `validate:"required" json:"packagesName,omitempty"`   //套餐名称，个签用户返回
	PackagesAmount float64 `validate:"required" json:"packagesAmount,omitempty"` //套餐价格，个签用户返回
	BatteryType    uint    `validate:"required" json:"batteryType"`              //电池型号 60 / 72

	GroupName string `validate:"required" json:"groupName,omitempty"` //团签公司名称，团签签用户返回

	OrderNo string      `validate:"required" json:"orderNo,omitempty"` //订单编号，个签用户返回
	Amount  float64     `validate:"required" json:"amount,omitempty"`  //支付金额，个签用户返回
	Earnest float64     `validate:"required" json:"earnest"`           //押金金额，个签用户返回
	PayType uint        `validate:"required" json:"payType,omitempty"` //支付方式 1 支付宝 2 微信，个签用户返回
	PayAt   *gtime.Time `validate:"required" json:"payAt,omitempty"`   //支付时间，个签用户返回

	ClaimState uint `validate:"required" json:"claimState"` // 1 未被认领， 2 已被认领
}

// ShopManagerPackagesOrderListDetailRep 订单记录获取订单详情响应数据
type ShopManagerPackagesOrderListDetailRep struct {
	UserName     string `validate:"required" json:"userName"`     //客户名称
	UserMobile   string `validate:"required" json:"userMobile"`   //客户手机号
	PackagesName string `validate:"required" json:"packagesName"` //套餐名称
	BatteryType  uint   `validate:"required" json:"batteryType"`  //电池型号 60 / 72

	OrderNo string      `validate:"required" json:"orderNo"` //订单编号
	Amount  float64     `validate:"required" json:"amount"`  //支付金额
	Earnest float64     `validate:"required" json:"earnest"` //押金金额
	PayType uint        `validate:"required" json:"payType"` //支付方式 1 支付宝 2 微信
	PayAt   *gtime.Time `validate:"required" json:"payAt"`   //支付时间
}

// ShopManagerPackagesOrderClaimReq 店长认领订单，请求数据
type ShopManagerPackagesOrderClaimReq struct {
	Code string `validate:"required" json:"code"` //扫码获取
}

// Prepay 发起支付数据
type Prepay struct {
	Description string
	No          string
	Amount      float64
	NotifyUrl   string
}

//BizNewCdeReq 店长认领订单扫码获取订单请求数据
type BizNewCdeReq struct {
	Code string `validate:"required" v:"required" json:"code"` //code二维码
}
