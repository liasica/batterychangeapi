package model

import (
	"battery/app/model/internal"
)

// Packages is the golang structure for table packages.
type Packages internal.Packages

// PackagesListUserReq 用户获取套餐列表请求
type PackagesListUserReq struct {
	Page
	CityId uint64 `validate:"required" json:"cityId"` // 所在城市
}

// PackagesListUserRep 用户获取套餐列表响应
type PackagesListUserRep []PackagesRiderListRepItem

type PackagesRiderListRepItem struct {
	Id             uint    `json:"id"`             //
	Name           string  `json:"name"`           // 名称
	Days           uint    `json:"days"`           // 套餐时长天数
	BatteryType    uint    `json:"batteryType"`    // 60 / 72 电池类型
	Amount         float64 `json:"amount"`         // 套餐价格(包含押金)
	Earnest        float64 `json:"earnest"`        // 押金
	UsableCityName string  `json:"usableCityName"` // 可用城市名称
	UsableShopCnt  int     `json:"usableShopCnt"`  // 可用店铺数量
	Desc           string  `json:"desc"`           // 套餐介绍
}
