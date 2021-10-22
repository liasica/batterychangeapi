package model

import (
    "battery/app/model/internal"
    "github.com/gogf/gf/os/gtime"
)

type ShopBatteryRecord internal.ShopBatteryRecord

const (
    ShopBatteryRecordTypeIn  = 1 // 入库
    ShopBatteryRecordTypeOut = 2 // 出库
)

// ShopBatteryRecordStatRep 门店资电池产管理统计
type ShopBatteryRecordStatRep struct {
    InTotal  int `json:"inTotal"`  // 入库总数
    OutTotal int `json:"outTotal"` // 出库总数
}

// ShopBatteryRecordListReq 门店电池管理明细请求
type ShopBatteryRecordListReq struct {
    Page
    Type      uint        `json:"type"`      // 1入库 2出库
    StartDate *gtime.Time `json:"startDate"` // 开始日期 eg: 2021-10-17
    EndDate   *gtime.Time `json:"endDate"`   // 结束日期 eg: 2021-10-17
}

// ShopBatteryRecordListWithDateGroup 按天统计电池出入库情况
type ShopBatteryRecordListWithDateGroup struct {
    Date     string                      `json:"date"`     // 日期
    InTotal  int                         `json:"inTotal"`  // 入库总数
    OutTotal int                         `json:"outTotal"` // 出库总数
    Items    []ShopBatteryRecordListItem `json:"items"`    // 详细
}

// ShopBatteryRecordListItem 电池出入库列表项
type ShopBatteryRecordListItem struct {
    BizType     uint   `json:"bizType" enums:"0,1,2,3,4,5"` // 0平台调拨, 1新签, 2换电池, 3寄存, 4恢复计费, 5退租
    UserName    string `json:"userName"`                    // 骑手名字
    Num         int    `json:"num"`                         // 数量
    BatteryType string `json:"batteryType"`                 // 60 / 72
    Date        string `json:"date"`                        // 调拨日期
}

// ShopBatteryRecordListRep 门店电池管理明细
type ShopBatteryRecordListRep struct {
    BizType     uint        `json:"bizType" enums:"0,1,2,3,4,5"` // 0平台调拨, 1新签, 2换电池, 3寄存, 4恢复计费, 5退租
    UserName    string      `json:"userName"`                    // 骑手名字
    Num         uint        `json:"num"`                         // 数量
    BatteryType string      `json:"batteryType"`                 // 60 / 72
    At          *gtime.Time `json:"at"`                          // 时间
    DayCnt      uint        `json:"dayCnt"`                      // 当天总数
}

// BatteryRecordListReq 电池日志请求
type BatteryRecordListReq struct {
    ShopBatteryRecordListReq
    ShopId uint `json:"shopId"` // 门店ID
    Page
}

// BatteryRecordListItem 电池日志项
type BatteryRecordListItem struct {
    Id          uint        `json:"id"`                          // ID
    ShopId      uint        `json:"shopId"`                      // 门店ID
    BizType     uint        `json:"bizType" enums:"0,1,2,3,4,5"` // 操作类别: 0平台调拨, 1新签, 2换电池, 3寄存, 4恢复计费, 5退租
    UserName    string      `json:"userName"`                    // 操作员: 平台调拨为0
    Num         uint        `json:"num"`                         // 数量
    BatteryType string      `json:"batteryType" enums:"60,72"`   // 电池型号: 60 / 72
    CreatedAt   *gtime.Time `json:"createdAt"`                   // 操作时间
}

// BatterTransferReq 电池调拨请求体
type BatterTransferReq struct {
    BatteryType string `json:"batteryType" enums:"60,72" v:"required|in:60,72"` // 电池型号
    From        uint   `json:"from" v:"required|integer|min:0"`                 // 调出自 0平台 其他店铺ID
    To          uint   `json:"to" v:"required|integer|min:0"`                   // 调入至 0平台 其他店铺ID
    Num         uint   `json:"num" v:"required|integer|between:1,9999"`         // 数量
}
