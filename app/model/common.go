package model

import "github.com/gogf/gf/os/gtime"

const (
	ContextAdminKey       = "CurrentAdmin"
	ContextRiderKey       = "CurrentRider"
	ContextShopManagerKey = "CurrentShopManager"
)

// ContextAdmin 后台管理员上下文
type ContextAdmin struct {
	Id       uint
	Username string
}

// ContextShop 店长上下文
type ContextShop struct {
	Id       uint
	Username string
	ShopId   int
}

// ContextRider 骑手上下文
type ContextRider struct {
	Id                       uint64
	Mobile                   string
	GroupId                  uint
	PackagesId               uint
	PackagesOrderId          uint64
	BatteryType              uint
	BatteryState             uint
	AuthState                uint
	SignState                uint
	EsignAccountId           string
	Qr                       string
	ExpirationAt             *gtime.Time
	BizBatterySecondsStartAt *gtime.Time
	BizBatteryRenewalCnt     uint
	BizBatteryRenewalSeconds uint
}

// ContextShopManager 店长上下文
type ContextShopManager struct {
	Id     uint64
	ShopId uint
	Mobile string
}

// Page 分页参数
type Page struct {
	PageIndex int `validate:"required" v:"required|integer|min:1"`         //当前页号
	PageLimit int `validate:"required" v:"required|integer|min:1|max:100"` //页大小
}

// IdReq ID参数
type IdReq struct {
	Id uint64 `validate:"required" v:"required|integer|min:0"` //ID
}

// UploadImageRep 图片上传返回信息
type UploadImageRep struct {
	Path string `json:"path"` //图片路径
}
