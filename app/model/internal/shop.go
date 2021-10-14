// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/os/gtime"
)

// Shop is the golang structure for table shop.
type Shop struct {
	Id              uint        `orm:"id,primary"      json:"id"`              //
	State           uint        `orm:"state"           json:"state"`           // 门店状态 0 休息总，1 营业中，2 外出中
	ManagerName     string      `orm:"managerName"     json:"managerName"`     //
	Name            string      `orm:"name,unique"     json:"name"`            // 门店名称
	Mobile          string      `orm:"mobile,unique"   json:"mobile"`          // 手机号
	ReturnAt        *gtime.Time `orm:"returnAt"        json:"returnAt"`        // 外出大致返回时间
	BatteryOutCnt60 uint        `orm:"batteryOutCnt60" json:"batteryOutCnt60"` // 60伏电池出库数量
	BatteryInCnt60  uint        `orm:"batteryInCnt60"  json:"batteryInCnt60"`  // 60伏电池入库数量
	BatteryInCnt72  uint        `orm:"batteryInCnt72"  json:"batteryInCnt72"`  // 72伏电池入库数量
	BatteryOutCnt72 uint        `orm:"batteryOutCnt72" json:"batteryOutCnt72"` // 72伏电池出库数量
	ChargerInCnt    uint        `orm:"chargerInCnt"    json:"chargerInCnt"`    // 充电器入库数量
	ChargerOutCnt   uint        `orm:"chargerOutCnt"   json:"chargerOutCnt"`   // 充电器出库数量
	BatteryCnt72    int         `orm:"batteryCnt72"    json:"batteryCnt72"`    // 72伏电池数量
	BatteryCnt60    int         `orm:"batteryCnt60"    json:"batteryCnt60"`    // 60伏电池数量
	ChargerCnt      int         `orm:"chargerCnt"      json:"chargerCnt"`      // 充电器数量
	Lng             float64     `orm:"lng"             json:"lng"`             // 经度
	Lat             float64     `orm:"lat"             json:"lat"`             // 纬度
	Qr              string      `orm:"qr,unique"       json:"qr"`              // 二维码编号
	ProvinceId      uint        `orm:"provinceId"      json:"provinceId"`      // 省级行政编码
	CityId          uint        `orm:"cityId"          json:"cityId"`          // 市级行政编码
	DistrictId      uint        `orm:"districtId"      json:"districtId"`      // 区县行政编码
	Address         string      `orm:"address"         json:"address"`         // 详细地址
	CreatedAt       *gtime.Time `orm:"createdAt"       json:"createdAt"`       //
	UpdatedAt       *gtime.Time `orm:"updatedAt"       json:"updatedAt"`       //
}
