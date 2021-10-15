package model

import (
    "encoding/json"
    "github.com/gogf/gf/os/gtime"
)

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
    Id                           uint64
    Mobile                       string
    GroupId                      uint
    PackagesId                   uint
    PackagesOrderId              uint64
    BatteryType                  uint
    BatteryState                 uint
    AuthState                    uint
    EsignAccountId               string
    Qr                           string
    BatteryReturnAt              *gtime.Time
    BizBatterySecondsStartAt     *gtime.Time
    BizBatteryRenewalCnt         uint
    BizBatteryRenewalDays        uint
    BizBatteryRenewalDaysStartAt *gtime.Time
}

// ContextShopManager 店长上下文
type ContextShopManager struct {
    Id     uint64
    ShopId uint
    Mobile string
}

// Page 分页参数
type Page struct {
    PageIndex int `validate:"required" v:"required|integer|min:1"`         // 当前页号
    PageLimit int `validate:"required" v:"required|integer|min:1|max:100"` // 页大小
}

// ItemsWithTotal 带有总数的列表数据
type ItemsWithTotal struct {
    Total int           `json:"total"` // 总数
    Items []interface{} `json:"items"` // 数据
}

// IdReq ID参数
type IdReq struct {
    Id uint64 `validate:"required" v:"required|integer|min:0"` // ID
}

// UploadRep 图片上传返回信息
type UploadRep struct {
    Path string `json:"path"` // 文件路径
}

// ImageBase64Req base64图片请求
type ImageBase64Req struct {
    Base64Content string `json:"base64Content" validate:"required" v:"required"` // 图片内容
}

type ArrayString []string

func (required *ArrayString) UnmarshalValue(value interface{}) error {
    if value != nil {
        switch value.(type) {
        case string:
            data := []byte(value.(string))
            return json.Unmarshal(data, &required)
        case []interface{}:
            data := value.([]interface{})
            fields := make([]string, len(data))
            for key, field := range data {
                fields[key] = field.(string)
            }
            *required = fields
        }
    }
    return nil
}

type ArrayUint64 []uint64

func (arr *ArrayUint64) UnmarshalValue(value interface{}) error {
    if value != nil {
        switch value.(type) {
        case string:
            return json.Unmarshal([]byte(value.(string)), &arr)
        case []interface{}:
            data := value.([]interface{})
            fields := make([]uint64, len(data))
            for key, field := range data {
                fields[key] = field.(uint64)
            }
            *arr = fields
        }
    }
    return nil
}
