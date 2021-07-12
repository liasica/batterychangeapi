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
// @param 	pageIndex query integer  true "当前页码"
// @param 	pageLimit query integer  true "每页行数"
// @param 	cityId query integer  true "当前城市ID"
// @router  /rapi/packages [GET]
// @success 200 {object} response.JsonResponse{data=model.PackagesListUserRep}  "返回结果"
func (*packagesApi) List(r *ghttp.Request) {
	var req model.PackagesListUserReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	response.JsonOkExit(r, service.PackagesService.ListUser(r.Context(), req))
}
