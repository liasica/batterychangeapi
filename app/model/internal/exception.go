// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/os/gtime"
)

// Exception is the golang structure for table exception.
type Exception struct {
	Id          uint64      `orm:"id,primary"  json:"id"`          //
	State       uint        `orm:"state"       json:"state"`       // 状态 0未解决 1已解决
	Type        uint        `orm:"type"        json:"type"`        // 1 遗失  2 故障
	BatteryType uint        `orm:"batteryType" json:"batteryType"` // 电池型号 60 / 72
	RecoverType uint        `orm:"recoverType" json:"recoverType"` // 1 用户 2 店长
	Detail      string      `orm:"detail"      json:"detail"`      // 详细说明
	Img         string      `orm:"img"         json:"img"`         // 图片链接
	Reason      int         `orm:"reason"      json:"reason"`      // 故障原因 0 其它 1 插头故障 2 无电压
	CreatedAt   *gtime.Time `orm:"createdAt"   json:"createdAt"`   //
	UpdatedAt   *gtime.Time `orm:"updatedAt"   json:"updatedAt"`   //
}
