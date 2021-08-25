package model

import (
	"battery/app/model/internal"
)

type Exception struct {
	Exception internal.Exception
	Img       []string
}

// ExceptionReportReq 异常上报请求
type ExceptionReportReq struct {
	Type        uint     `validate:"required" v:"required" json:"type"`            // 1 遗失  2 故障
	BatteryType uint     `validate:"required" v:"required" json:"batteryType"`     // 电池型号 60 / 72
	RecoverType uint     `validate:"required" v:"required" json:"recoverType"`     // 1 用户 2 店长
	Detail      string   `validate:"required" v:"required" json:"detail"`          // 详细说明
	Img         []string `validate:"required" v:"required" json:"img"`             // 图片链接数组
	Reason      int      `validate:"required" v:"required|in:0,1,2" json:"reason"` // 故障原因 0 其它 1 插头故障 2 无电压
}
