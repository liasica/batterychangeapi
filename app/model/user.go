package model

import (
    "github.com/gogf/gf/os/gtime"
    "github.com/gogf/gf/util/gmeta"

    "battery/app/model/internal"
)

// 实名认证状态
const (
    AuthStateVerifyDefault = iota // 未提交
    AuthStateVerifyWait           // 待审核
    AuthStateVerifySuccess        // 审核通过
    AuthStateDefaultFailed        // 审核失败
)

// 换电状态
const (
    BatteryStateDefault = iota // 未开通
    BatteryStateNew            // 新购待领取 （团用户未领取）
    BatteryStateUse            // 租借中
    BatteryStateSave           // 寄存中
    BatteryStateExit           // 已退租
    BatteryStateOverdue        // 已逾期
    BatteryStateExpired        // 已过期
)

// 用户属性
const (
    UserTypePersonal   = 1 + iota // 个签用户
    UserTypeGroupRider            // 团签骑手
    UserTypeGroupBoss             // 团签BOSS
)

type User internal.User

// UserVerifyReq 用户认证列表请求
type UserVerifyReq struct {
    Page
    RealName  string `json:"realName"`                  // 姓名
    Mobile    string `json:"mobile"`                    // 手机号
    Type      uint   `json:"type"`                      // 用户类别
    AuthState uint   `json:"authState" enums:"0,1,2,3"` // 认证状态 0未提交 1待审核 2审核通过 3审核失败
}

// UserRegisterReq 用户注册请求数据
type UserRegisterReq struct {
    Mobile string `validate:"required" v:"required|phone-loose"` // 手机号
    Sms    string `validate:"required" v:"required|length:6,6"`  // 短信验证码
}

// UserLoginReq 用户登录请求数据
type UserLoginReq struct {
    Mobile string `validate:"required" v:"required|phone-loose"` // 手机号
    Sms    string `validate:"required" v:"required|length:6,6"`  // 短信验证码
}

// UserLoginRep 用户登录返回数据
type UserLoginRep struct {
    AccessToken string `validate:"required" json:"accessToken"` // 请求 token
    Type        uint   `validate:"required" json:"type"`        // 用户角色 1 个签骑手 2 团签骑手 3 团签BOSS
    AuthState   uint   `validate:"required" json:"authState"`   // 实名认证状态 0 未提交 ，1 待审核， 2 审核通过，3 审核未通过
}

// UserRealNameAuthReq 实名认证请求数据
type UserRealNameAuthReq struct {
    RealName   string `validate:"required" v:"required|length:2,10"`              // 真实姓名
    IdCardNo   string `validate:"required" v:"required|length:15,18|resident-id"` // 身份证号码
    IdType     string // 证件类型 CRED_PSN_CH_IDCARD 中国大陆居民身份证 CRED_PSN_CH_HONGKONG 香港来往大陆通行证 CRED_PSN_CH_MACAO 澳门来往大陆通行证 CRED_PSN_CH_TWCARD 台湾来往大陆通行证 CRED_PSN_PASSPORT 护照
    IdCardImg1 string // 身份证正面照片
    IdCardImg2 string // 身份证反面照片
    IdCardImg3 string // 身份证手持照片
}

// UserRealNameAuthRep 实名认证响应数据
type UserRealNameAuthRep struct {
    FlowId    string `json:"flowId"`       // 流程ID
    ShortLink string `json:"shortLink"`    // 短地址
    Url       string `validate:"required"` // 骑手实名认证连接地址
}

// RealNameAuthVerifyProfileRep 获取用户实名认证提交资料信息
type RealNameAuthVerifyProfileRep struct {
    Id         uint64 `orm:"id,primary"         json:"id"`         //
    Mobile     string `orm:"mobile,unique"      json:"mobile"`     //
    RealName   string `orm:"realName"           json:"realName"`   // 真实姓名
    IdCardNo   string `orm:"idCardNo"           json:"idCardNo"`   // 身份证号码
    IdCardImg1 string `orm:"idCardImg1"         json:"idCardImg1"` // 正面图
    IdCardImg2 string `orm:"idCardImg2"         json:"idCardImg2"` // 反面图
    IdCardImg3 string `orm:"idCardImg3"         json:"idCardImg3"` // 人像图
    AuthState  uint   `orm:"authState"          json:"authState"`  // 实名认证状态 0 未提交 ，1 待审核， 2 审核通过，3 审核未通过
}

// RealNameAuthVerifyReq 实名认证审核请求数据
type RealNameAuthVerifyReq struct {
    AuthState uint `v:"required|in:2,3"` // 审核结果 2 通过 3 失败
}

// BizProfileRep 用户业务办理店长扫码获取用户信息
type BizProfileRep struct {
    Id           uint64 `json:"id"`
    Mobile       string `json:"mobile"`
    RealName     string `json:"realName"`            // 真实姓名
    IdCardNo     string `json:"idCardNo"`            // 身份证号码
    AuthState    uint   `json:"authState"`           // 实名认证状态 0 未提交 ，1 待审核， 2 审核通过，3 审核未通过
    BatteryState uint   `json:"batteryState"`        // 换电状态：0 未开通，1 新购待领取 （团用户未领取），2 租借中，3 寄存中，4 已退租 5 已逾期
    BatteryType  string `json:"batteryType"`         // 电池类型 60 / 72
    ComboName    string `json:"comboName,omitempty"` // 套餐名称
    GroupId      uint   `json:"groupId"`             // 团体Id，个签用户为 0
    GroupName    string `json:"groupName,omitempty"` // 团体名称
}

// PushTokenReq 上报用户的推送token
type PushTokenReq struct {
    DeviceType  int    `validate:"required" v:"required|in:1,2" json:"deviceType"` // 1 android  2 ios
    DeviceToken string `validate:"required" v:"required" json:"deviceToken"`       // token 推送平台用户ID
}

// UserProfileRep 骑手端用户信息概况
type UserProfileRep struct {
    Name      string `json:"name"`                          // 姓名
    Mobile    string `json:"mobile"`                        // 手机号码
    Type      uint   `validate:"required" json:"type"`      // 用户角色 1 个签骑手 2 团签骑手 3 团签BOSS
    AuthState uint   `validate:"required" json:"authState"` // 实名认证状态 0 未提交 ，1 待审核， 2 审核通过，3 审核未通过
    Qr        string `json:"qr"`                            // 用户二维码数据，需要本地生成图片
    User      struct {
        BatteryState    uint        `json:"batteryState"`    // 个签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租 5已逾期
        ComboId         uint        `json:"comboId"`         // 个签骑手所购套餐ID
        ComboName       string      `json:"comboName"`       // 个签骑手所购套餐名称
        BatteryReturnAt *gtime.Time `json:"batteryReturnAt"` // 个签骑手套餐到期时间
    } `json:"user,omitempty"` // 个签用户套餐信息， 其它类型用户忽略

    GroupUser struct {
        BatteryState uint   `json:"batteryState"` // 团签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租
        BatteryType  string `json:"batteryType"`  // 电池型号 60 / 72  未开通为 0
    } `json:"groupUser,omitempty"` // 团签用户骑手信息， 其它类型用户忽略

    GroupBoss UserGroupStatRep `json:"groupBoss,omitempty"` // 团队BOSS信息， 其它类型用户忽略
}

// UserVerifyListItem 骑手实名列表项
type UserVerifyListItem struct {
    Id         uint   `json:"id"`                        // 骑手ID
    RealName   string `json:"realName"`                  // 姓名
    Mobile     string `json:"mobile"`                    // 手机号
    Type       uint   `json:"type"`                      // 用户类别
    AuthState  uint   `json:"authState" enums:"0,1,2,3"` // 认证状态 0未提交 1待审核 2审核通过 3审核失败
    IdCardNo   string `json:"idCardNo"`                  // 身份证
    IdCardImg1 string `json:"idCardImg1"`                // 身份证人像面
    IdCardImg2 string `json:"idCardImg2"`                // 身份证国徽面
    IdCardImg3 string `json:"idCardImg3"`                // 手持身份证
}

// UserListReq 用户列表请求
type UserListReq struct {
    Page
    GroupId      uint        `json:"groupId" swaggerignore:"true"` // 团队ID
    RealName     string      `json:"realName"`                     // 姓名
    Mobile       string      `json:"mobile"`                       // 手机号
    BatteryState uint        `json:"batteryState"`                 // 换电状态 个签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租 5已逾期; 团签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租
    StartDate    *gtime.Time `json:"startDate"`                    // 开始日期 eg: 2021-10-17
    EndDate      *gtime.Time `json:"endDate"`                      // 结束日期 eg: 2021-10-17
}

// UserListItem 用户列表项
type UserListItem struct {
    gmeta.Meta `json:"-" orm:"table:user" swaggerignore:"true"`

    Id                    uint        `json:"id"`                    // ID
    GroupId               uint        `json:"groupId"`               // 团队ID
    RealName              string      `json:"realName"`              // 姓名
    Mobile                string      `json:"mobile"`                // 手机号
    BatteryState          uint        `json:"batteryState"`          // 状态 个签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租 5已逾期; 团签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租
    BatteryType           string      `json:"batteryType"`           // 套餐电池型号 60/72
    ComboId               uint        `json:"comboId"`               // 套餐ID
    ComboOrderId          uint64      `json:"comboOrderId"`          // 办理套餐订单ID
    BatteryReturnAt       *gtime.Time `json:"batteryReturnAt"`       // 个人用户应归还电池时间， 小于当前时间即逾期
    BatterySaveAt         *gtime.Time `json:"batterySaveAt"`         // 个签用户电池寄存时间
    BizBatteryRenewalCnt  uint        `json:"bizBatteryRenewalCnt"`  // 积累换次数
    BizBatteryRenewalDays uint        `json:"bizBatteryRenewalDays"` // 累计换电自然天数
    StartDate             *gtime.Time `json:"startDate"`             // 开始日期 eg: 2021-10-17
    EndDate               *gtime.Time `json:"endDate"`               // 结束日期 eg: 2021-10-17
    CreatedAt             *gtime.Time `json:"createdAt"`             // 注册时间

    GroupName string `json:"groupName"` // 团队名称
    ComboName string `json:"comboName"` // 当前套餐名称
    ComboType uint   `json:"comboType"` // 当前套餐类型
    Days      uint   `json:"days"`      // 剩余天数

    ComboDetail *Combo `json:"-" orm:"with:id=comboId"`
    Group       *Group `json:"-" orm:"with:id=groupId"`
    // BizItems    []*UserBiz `json:"-" orm:"with:userId=id, order:createdAt desc"`
}

type UserSignFileRepItem struct {
    FileName string `json:"fileName"` // 文件名称
    FileUrl  string `json:"fileUrl"`  // 文件地址
}

type UserSignFileRep []*UserSignFileRepItem
