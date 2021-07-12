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
