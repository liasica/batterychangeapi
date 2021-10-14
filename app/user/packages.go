package user

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var PackagesApi = packagesApi{}

type packagesApi struct {
}

// List 个签用户套餐列表
// @summary 骑手-个签用户套餐列表
// @tags    骑手
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	cityId query integer  true "当前城市ID"
// @router  /rapi/packages [GET]
// @success 200 {object} response.JsonResponse{data=model.PackagesListUserRep}  "返回结果"
func (*packagesApi) List(r *ghttp.Request) {
    var req model.PackagesListUserReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r, service.PackagesService.ListUser(r.Context(), req))
}

// Detail 个签用户套餐详情
// @summary 骑手-个签用户套餐详情
// @tags    骑手
// @Accept  json
// @Produce  json
// @Param 	id path integer  true "套餐ID"
// @router  /rapi/packages/:id [GET]
// @success 200 {object} response.JsonResponse{data=model.PackagesRiderListRepItem}  "返回结果"
func (*packagesApi) Detail(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    packages, err := service.PackagesService.Detail(r.Context(), uint(req.Id))
    if err != nil {
        response.JsonErrExit(r, response.RespCodeNotFound)
    }
    response.JsonOkExit(r, model.PackagesRiderListRepItem{
        Id:          packages.Id,
        Name:        packages.Name,
        Days:        packages.Days,
        BatteryType: packages.BatteryType,
        Amount:      packages.Amount,
        Earnest:     packages.Earnest,
        Desc:        packages.Desc,
    })
}
