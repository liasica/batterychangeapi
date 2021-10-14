package model

import (
    "battery/app/model/internal"
)

const DistrictsDefaultCityCode = "010" // 默认北京

const (
    DistrictLevelProvince = "province"
    DistrictLevelCity     = "city"
    DistrictLevelDistrict = "district"
    DistrictLevelStreet   = "street"
)

type Districts internal.Districts

// DistrictsChildRep 获取行政地区下级返回信息
type DistrictsChildRep struct {
    Id   uint   `validate:"required" json:"id"`
    Name string `validate:"required" json:"name"`
}

// DistrictsCurrentCityReq 获取当前城市请求
type DistrictsCurrentCityReq struct {
    Lng float64 `validate:"required" json:"lng" v:"required"`
    Lat float64 `validate:"required" json:"lat" v:"required"`
}

// DistrictsCurrentCityRep 获取当前城市响应
type DistrictsCurrentCityRep struct {
    Id   uint   `validate:"required" json:"id"`
    Name string `validate:"required" json:"name"`
}

// OpenCityListRepItem 获取开城市列表返回
type OpenCityListRepItem struct {
    Id     uint    `validate:"required" json:"id"`     // ID
    Name   string  `validate:"required" json:"name"`   // 名称
    AdCode uint    `validate:"required" json:"adCode"` // 行政区代码
    Lng    float64 `validate:"required" json:"lng"`    // 经度
    Lat    float64 `validate:"required" json:"lat"`    // 纬度
}
