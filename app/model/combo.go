package model

import (
    "battery/app/model/internal"
)

// Combo is the golang structure for table combo.
type Combo internal.Combo

// ComboListUserReq 用户获取套餐列表请求
type ComboListUserReq struct {
    Page
    CityId uint64 `validate:"required" json:"cityId"` // 所在城市
}

// ComboListUserRep 用户获取套餐列表响应
type ComboListUserRep []ComboRiderListRepItem

type ComboRiderListRepItem struct {
    Id             uint    `json:"id"`                        //
    Name           string  `json:"name"`                      // 名称
    Days           uint    `json:"days"`                      // 套餐时长天数
    BatteryType    string  `json:"batteryType" enums:"60,72"` // 60 / 72 电池类型
    Amount         float64 `json:"amount"`                    // 套餐价格(包含押金)
    Deposit        float64 `json:"deposit"`                   // 押金
    UsableCityName string  `json:"usableCityName"`            // 可用城市名称
    UsableShopCnt  int     `json:"usableShopCnt"`             // 可用门店数量
    Desc           string  `json:"desc"`                      // 套餐介绍
}

// ComboListItem 套餐列表项
type ComboListItem struct {
    Id          uint    `json:"id"`
    Name        string  `json:"name"`                      // 套餐名
    Type        uint    `json:"type" enums:"1,2"`          // 套餐类型: 1 个人 2 团体
    Amount      float64 `json:"amount"`                    // 套餐总价（包含押金）
    Price       float64 `json:"price"`                     // 套餐价格
    BatteryType string  `json:"batteryType" enums:"60,72"` // 电池类型
    Days        uint    `json:"days"`                      // 套餐时长天数
    Deposit     float64 `json:"deposit"`                   // 押金
    ProvinceId  uint    `json:"provinceId"`                // 省份ID
    CityId      uint    `json:"cityId"`                    // 城市ID
    DeleteAt    string  `json:"deleteAt"`                  // 套餐停用时间
    Desc        string  `json:"desc"`                      // 套餐详情
}

// ComboReq 套餐请求
type ComboReq struct {
    BatteryType string  `v:"required|in:60,72" json:"batteryType" enums:"60,72"`                           // 电池类型
    Type        uint    `v:"required" json:"type" enums:"1,2"`                                             // 套餐类型: 1个人 2团体
    Name        string  `v:"required" json:"name"`                                                         // 名称
    Days        uint    `v:"required|integer|min:1" json:"days"`                                           // 套餐时长天数
    Price       float64 `v:"required|regex:'/(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})$)/'" json:"price"`    // 套餐价格
    Deposit     float64 `v:"required|regex:'/(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})?$)/'" json:"deposit"` // 押金
    ProvinceId  uint    `v:"required|integer|min:1" json:"provinceId"`                                     // 省份ID
    CityId      uint    `v:"required|integer|min:1" json:"cityId"`                                         // 城市ID
    Desc        string  `v:"required" json:"desc"`                                                         // 介绍
    Disable     bool    `json:"disable,omitempty"`                                                         // 是否禁用，当值为true的时候其他参数忽略 false启用(默认) true禁用
}
