package user

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var ComboApi = comboApi{}

type comboApi struct {
}

// List 个签用户套餐列表
// @Summary 骑手-个签用户套餐列表
// @Tags    骑手
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	cityId query integer  true "当前城市ID"
// @Router  /rapi/combo [GET]
// @Success 200 {object} response.JsonResponse{data=model.ComboListUserRep}  "返回结果"
func (*comboApi) List(r *ghttp.Request) {
    var req model.ComboListUserReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r, service.ComboService.ListUser(r.Context(), req))
}

// Detail 个签用户套餐详情
// @Summary 骑手-个签用户套餐详情
// @Tags    骑手
// @Accept  json
// @Produce  json
// @Param 	id path integer  true "套餐ID"
// @Router  /rapi/combo/:id [GET]
// @Success 200 {object} response.JsonResponse{data=model.ComboRiderListRepItem}  "返回结果"
func (*comboApi) Detail(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    combo, err := service.ComboService.Detail(r.Context(), uint(req.Id))
    if err != nil {
        response.JsonErrExit(r, response.RespCodeNotFound)
    }
    response.JsonOkExit(r, model.ComboRiderListRepItem{
        Id:          combo.Id,
        Name:        combo.Name,
        Days:        combo.Days,
        BatteryType: combo.BatteryType,
        Amount:      combo.Amount,
        Deposit:     combo.Deposit,
        Desc:        combo.Desc,
    })
}
