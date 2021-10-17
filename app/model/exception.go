package model

import (
    "battery/app/model/internal"
    "github.com/gogf/gf/util/gmeta"
    "time"
)

const (
    ExceptionStatePending   = iota // 待处理
    ExceptionStateProcessed        // 已处理
)

const (
    ExceptionReasonOther   = iota // 其他
    ExceptionReasonPlug           // 插头故障
    ExceptionReasonVoltage        // 电压故障
)

const (
    ExceptionTypeLost  = 1 + iota // 遗失
    ExceptionTypeFault            // 故障
)

const (
    ExceptionDiscoverDriver = 1 + iota // 骑手
    ExceptionDiscoverShop              // 店长
)

type Exception struct {
    Exception internal.Exception
    Img       ArrayString
}

type ExceptionEntity struct {
    gmeta.Meta `orm:"table:exception"`

    internal.Exception
    //
    // Id          uint64      `orm:"id,primary"  json:"id"`          //
    // ShopId      uint        `orm:"shopId"      json:"shopId"`      // 门店
    // State       uint        `orm:"state"       json:"state"`       // 状态 0未解决 1已解决
    // Type        uint        `orm:"type"        json:"type"`        // 1 遗失  2 故障
    // BatteryType uint        `orm:"batteryType" json:"batteryType"` // 电池型号 60 / 72
    // Discoverer  uint        `orm:"discoverer"  json:"discoverer"`  // 发现人 1 用户 2 店长
    // Detail      string      `orm:"detail"      json:"detail"`      // 详细说明
    // Img         string      `orm:"img"         json:"img"`         // 图片链接
    // Reason      int         `orm:"reason"      json:"reason"`      // 故障原因 0 其它 1 插头故障 2 无电压
    // CreatedAt   *gtime.Time `orm:"createdAt"   json:"createdAt"`   //
    // UpdatedAt   *gtime.Time `orm:"updatedAt"   json:"updatedAt"`   //

    ShopDetail *Shop `orm:"with:id=shopId"`
}

// ExceptionReportReq 异常上报请求
type ExceptionReportReq struct {
    ShopId      uint     `json:"shopId"`                                                     // 店铺ID
    Type        uint     `validate:"required" v:"required" json:"type" enums:"1,2"`          // 异常类别: 1遗失 2故障
    BatteryType uint     `validate:"required" v:"required" json:"batteryType" enums:"60,72"` // 电池型号: 60 / 72
    Discoverer  uint     `validate:"required" v:"required" json:"discoverer" enums:"1,2"`    // 发现人: 1店长 2骑手
    Detail      string   `validate:"required" v:"required" json:"detail"`                    // 详细说明
    Img         []string `validate:"required" v:"required" json:"img"`                       // 图片链接数组
    Reason      int      `validate:"required" v:"required|in:0,1,2" json:"reason"`           // 故障原因 0 其它 1 插头故障 2 无电压
}

// ExceptionListReq 电池异常记录请求
type ExceptionListReq struct {
    Page
    ShopId    uint      `json:"shopId"`    // 店铺ID
    StartTime time.Time `json:"startTime"` // 开始时间
    EndTime   time.Time `json:"endTime"`   // 结束时间
}

// ExceptionListItem 异常列表项
type ExceptionListItem struct {
    gmeta.Meta `orm:"table:exception"`

    Id          uint        `json:"id"`                        // ID
    Type        uint        `json:"type" enums:"1,2"`          // 异常类别: 1遗失 2故障
    ShopId      uint        `json:"shopId"`                    // 门店ID
    ShopName    string      `json:"shopName"`                  // 门店名称
    CityId      uint        `json:"cityId"`                    // 城市ID
    CityName    string      `json:"cityName"`                  // 城市
    Name        string      `json:"name"`                      // 名称
    State       uint        `json:"state" enums:"0,1"`         // 状态: 0待处理 1已处理
    Reason      uint        `json:"reason"`                    // 故障原因: 0其它 1插头故障 2电压故障
    BatteryType uint        `json:"batteryType" enums:"60,72"` // 产品类型: 60/72
    Discoverer  uint        `json:"discoverer" enums:"1,2"`    // 发现人: 1店长 2骑手
    Detail      string      `json:"detail"`                    // 详细说明
    Img         ArrayString `json:"img"`                       // 图片

    ShopDetail *Shop `json:"-" orm:"with:id=shopId"`
}
