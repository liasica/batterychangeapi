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
    Id             uint    `json:"id"`                        //
    Name           string  `json:"name"`                      // 名称
    Days           uint    `json:"days"`                      // 套餐时长天数
    BatteryType    uint    `json:"batteryType" enums:"60,72"` // 60 / 72 电池类型
    Amount         float64 `json:"amount"`                    // 套餐价格(包含押金)
    Earnest        float64 `json:"earnest"`                   // 押金
    UsableCityName string  `json:"usableCityName"`            // 可用城市名称
    UsableShopCnt  int     `json:"usableShopCnt"`             // 可用门店数量
    Desc           string  `json:"desc"`                      // 套餐介绍
}

// PackageListItem 套餐列表项
type PackageListItem struct {
    Id          uint    `json:"id"`
    Name        string  `json:"name"`                      // 套餐名
    Amount      float64 `json:"amount"`                    // 套餐总价（包含押金）
    Price       float64 `json:"price"`                     // 套餐价格
    BatteryType uint    `json:"batteryType" enums:"60,72"` // 电池类型
    Days        uint    `json:"days"`                      // 套餐时长天数
    Earnest     float64 `json:"earnest"`                   // 押金
    ProvinceId  uint    `json:"provinceId"`                // 省份ID
    CityId      uint    `json:"cityId"`                    // 城市ID
    DeleteAt    string  `json:"deleteAt"`                  // 套餐停用时间
}

// PackageReq 套餐请求
type PackageReq struct {
    BatteryType uint    `v:"required|in:60,72" json:"batteryType" enums:"60,72"`                           // 电池类型
    Type        uint    `v:"required" json:"type"`                                                         // 套餐类型
    Name        string  `v:"required" json:"name"`                                                         // 名称
    Days        uint    `v:"required|integer|min:1" json:"days"`                                           // 套餐时长天数
    Price       float64 `v:"required|regex:'/(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})$)/'" json:"price"`    // 套餐价格
    Earnest     float64 `v:"required|regex:'/(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})?$)/'" json:"earnest"` // 押金
    ProvinceId  uint    `v:"required|integer|min:1" json:"provinceId"`                                     // 省份ID
    CityId      uint    `v:"required|integer|min:1" json:"cityId"`                                         // 城市ID
    Desc        string  `v:"required" json:"desc"`                                                         // 介绍
    Disable     uint8   `json:"disable,omitempty" enums:"0,1"`                                             // 0启用(默认) 1禁用
}
