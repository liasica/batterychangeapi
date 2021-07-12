package model

import (
	"battery/app/model/internal"
)

const DistrictsDefaultCityCode = "010" //默认北京

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
