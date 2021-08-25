package model

import (
	"github.com/gogf/gf/os/gtime"

	"battery/app/model/internal"
)

type ShopBatteryRecord internal.ShopBatteryRecord

const ShopBatteryRecordTypeIn = 1
const ShopBatteryRecordTypeOut = 2

// ShopBatteryRecordStatRep 店长资电池产管理统计
type ShopBatteryRecordStatRep struct {
	InTotal  uint `json:"inTotal"`  //入库总数
	OutTotal uint `json:"outTotal"` // 出库总数
}

// ShopBatteryRecordListReq 店长资电池管理明细请求
type ShopBatteryRecordListReq struct {
	Page
	Type      uint        `json:"type" v:"required|in:1,2"` //"1 入库 2 出库"
	StartTime *gtime.Time `json:"startTime"`
	EndTime   *gtime.Time `json:"endTime"`
}

// ShopBatteryRecordListRep 店长资电池管理明细
type ShopBatteryRecordListRep struct {
	BizType     uint        `json:"bizType"`     // 0 平台调拨,  1 新签,  2 换电池 ,3 寄存, 4 恢复计费 , 5 退租
	UserName    string      `json:"userName"`    //用户名  平台调拨为空
	Num         uint        `json:"num"`         //数量
	BatteryType uint        `json:"batteryType"` // 60 / 72
	At          *gtime.Time `json:"At"`          // 时间
	DayCnt      uint        `json:"dayCnt"`      //当天总数
}
