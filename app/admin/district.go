package admin

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var DistrictsApi = districtsApi{}

type districtsApi struct {
}

func (*districtsApi) Child(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    rep := service.DistrictsService.Child(req.Id)
    response.JsonOkExit(r, rep)
}

// List
// @summary 城市列表（三级联动）
// @tags    管理
// @Accept  json
// @Produce  json
// @router  /admin/districts [GET]
// @success 200 {object} response.JsonResponse{data=[]service.DistrictEl}  "返回结果"
func (*districtsApi) List(r *ghttp.Request) {
    response.JsonOkExit(r, service.DistrictsService.ListCityTree())
}
