package shop

import (
    "github.com/gogf/gf/net/ghttp"

    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
)

var (
    ManagerApi = managerApi{}
)

type managerApi struct {
}

// Login
// @Summary 店长-登录
// @Tags    店长
// @Accept  json
// @Produce  json
// @Param   entity  body model.ShopManagerLoginReq true "登录数据"
// @Router  /sapi/login [POST]
// @Success 200 {object} response.JsonResponse{data=model.ShopManagerLoginRep}  "返回结果"
func (*managerApi) Login(r *ghttp.Request) {
    var req model.ShopManagerLoginReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if rep, err := service.ShopManagerService.Login(r.Context(), req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    } else {
        response.JsonOkExit(r, rep)
    }
}

// Qr
// @Summary 店长-获取店本店二维码
// @Tags    店长
// @Accept  json
// @Produce  json
// @Router  /sapi/qr [GET]
// @Success 200 {object} response.JsonResponse 二维码结果 data 为二维码数据，需本地生成图片
func (*managerApi) Qr(r *ghttp.Request) {
    manager := r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    shop, _ := service.ShopService.Detail(r.Context(), manager.ShopId)
    response.JsonOkExit(r, shop.Qr)
}

// Profile
// @Summary 店长-获取门店信息
// @Tags    店长
// @Accept  json
// @Produce  json
// @Router  /sapi/shop/profile [GET]
// @Success 200 {object} response.JsonResponse{data=model.Shop}  "返回结果"
func (*managerApi) Profile(r *ghttp.Request) {
    manager := r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    rep, err := service.ShopService.Detail(r.Context(), manager.ShopId)
    if err != nil {
        response.JsonErrExit(r, response.RespCodeSystemError)
    }
    response.JsonOkExit(r, rep)
}

// ShopState
// @Summary 店长-修改门店状态
// @Tags    店长
// @Accept  json
// @Produce  json
// @Param   entity  body model.ShopManagerChangeStateReq true "请求数据"
// @Router  /sapi/shop/state [PUT]
// @Success 200 {object} response.JsonResponse{}  "返回结果"
func (*managerApi) ShopState(r *ghttp.Request) {
    var req model.ShopManagerChangeStateReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if err := service.ShopService.State(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId, req.State); err != nil {
        response.JsonErrExit(r, response.RespCodeSystemError)
    }
    response.JsonOkExit(r)
}

// PushToken
// @Summary 店长-上报推送token
// @Tags    店长
// @Accept  json
// @Produce  json
// @Param   entity  body model.PushTokenReq true "登录数据"
// @Router  /sapi/device  [PUT]
// @Success 200 {object} response.JsonResponse  "返回结果"
func (*managerApi) PushToken(r *ghttp.Request) {
    var req model.PushTokenReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if service.ShopManagerService.PushToken(r.Context(), req) != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}

// ResetMobile
// @Summary 店长-修改手机号码
// @Tags    店长
// @Accept  json
// @Produce  json
// @Param   entity  body model.ShopManagerResetMobileReq true "登录数据"
// @Router  /sapi/mobile  [PUT]
// @Success 200 {object} response.JsonResponse  "返回结果"
func (*managerApi) ResetMobile(r *ghttp.Request) {
    var req model.ShopManagerResetMobileReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if !service.SmsServer.Verify(r.Context(), model.SmsVerifyReq{
        Mobile: req.Mobile,
        Code:   req.Sms,
    }) {
        response.Json(r, response.RespCodeArgs, "手机号或验证码错误，修改失败")
    }
    if service.ShopManagerService.ResetMobile(r.Context(), req) != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}
