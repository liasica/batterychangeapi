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
    InTotal  uint `json:"inTotal"`  // 入库总数
    OutTotal uint `json:"outTotal"` // 出库总数
}

// ShopBatteryRecordListReq 店长电池管理明细请求
type ShopBatteryRecordListReq struct {
    Page
    Type      uint        `json:"type"`      // "1 入库 2 出库"
    StartTime *gtime.Time `json:"startTime"` // 查询范围-开始时间
    EndTime   *gtime.Time `json:"endTime"`   // 查询范围-结束时间
}

// ShopBatteryRecordListRep 店长电池管理明细
type ShopBatteryRecordListRep struct {
    BizType     uint        `json:"bizType" enums:"0, 1, 2, 3, 4, 5"` // 0平台调拨, 1新签, 2换电池, 3寄存, 4恢复计费, 5退租
    UserName    string      `json:"userName"`                         // 操作员 平台调拨为空
    Num         uint        `json:"num"`                              // 数量
    BatteryType uint        `json:"batteryType"`                      // 60 / 72
    At          *gtime.Time `json:"At"`                               // 时间
    DayCnt      uint        `json:"dayCnt"`                           // 当天总数
}

// BatteryRecordListReq 电池日志请求
type BatteryRecordListReq struct {
    ShopBatteryRecordListReq
    ShopId uint `json:"shopId"` // 门店ID
    Page
}

// BatteryRecordListItem 电池日志项
type BatteryRecordListItem struct {
    Id          uint        `json:"id"`                               // ID
    ShopId      uint        `json:"shopId"`                           // 门店ID
    BizType     uint        `json:"bizType" enums:"0, 1, 2, 3, 4, 5"` // 操作类别: 0平台调拨, 1新签, 2换电池, 3寄存, 4恢复计费, 5退租
    UserName    string      `json:"userName"`                         // 操作员: 平台调拨为0
    Num         uint        `json:"num"`                              // 数量
    BatteryType uint        `json:"batteryType" enums:"60,72"`        // 电池型号: 60 / 72
    CreatedAt   *gtime.Time `json:"createdAt"`                        // 操作时间
}
