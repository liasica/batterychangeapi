package user

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var ShopApi = shopApi{}

type shopApi struct {
}

// List
// @Summary 骑手-门店列表
// @Tags    骑手
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	cityId query integer  true "当前城市ID"
// @Param 	lng query number  true "经度"
// @Param 	lat query number  true "纬度"
// @Param 	name query string   false "门店名称"
// @Router  /rapi/shop [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopListUserRep}  "返回结果"
func (*shopApi) List(r *ghttp.Request) {
    var req model.ShopListUserReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r, service.ShopService.ListUser(r.Context(), req))
}
