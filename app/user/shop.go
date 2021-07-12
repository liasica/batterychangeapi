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
// @summary 骑手-店铺列表
// @tags    骑手
// @Produce  json
// @param 	pageIndex query integer  true "当前页码"
// @param 	pageLimit query integer  true "每页行数"
// @param 	cityId query integer  true "当前城市ID"
// @param 	lng query number  true "经度"
// @param 	lat query number  true "纬度"
// @param 	name query string   false "店铺名称"
// @router  /rapi/shop [GET]
// @success 200 {object} response.JsonResponse{data=model.ShopListUserRep}  "返回结果"
func (*shopApi) List(r *ghttp.Request) {
	var req model.ShopListUserReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	response.JsonOkExit(r, service.ShopService.ListUser(r.Context(), req))
}
