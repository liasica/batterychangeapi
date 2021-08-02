package user

import (
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
	"github.com/gogf/gf/net/ghttp"
)

var DistrictsApi = districtsApi{}

type districtsApi struct {
}

// CurrentCity
// @summary 骑手-定位当前城市
// @tags    骑手
// @Accept  json
// @Produce  json
// @param   entity  body model.DistrictsCurrentCityReq true "请求数据"
// @router  /rapi/districts/current_city [GET]
// @success 200 {object} response.JsonResponse{data=model.DistrictsCurrentCityRep}  "返回结果"
func (*districtsApi) CurrentCity(r *ghttp.Request) {
	var req model.DistrictsCurrentCityReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	rep, err := service.DistrictsService.CurrentCity(r.Context(), req)
	if err != nil {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	response.JsonOkExit(r, rep)
}

// OpenCityList
// @summary 骑手-获取已经开放的城市
// @tags    骑手
// @Produce  json
// @router  /rapi/open_city [GET]
// @success 200 {object} response.JsonResponse{data=[]model.OpenCityListRepItem}  "返回结果"
func (*districtsApi) OpenCityList(r *ghttp.Request) {
	cityIds := service.PackagesService.GetCityIds(r.Context())
	if len(cityIds) > 0 {
		rep := make([]model.OpenCityListRepItem, len(cityIds))
		districtsList := service.DistrictsService.GetByIds(r.Context(), cityIds)
		for key, city := range districtsList {
			rep[key] = model.OpenCityListRepItem{
				Id:   city.Id,
				Name: city.Name,
				AdCode: city.AdCode,
				Lng: city.Lng,
				Lat: city.Lat,
			}
		}
		response.JsonOkExit(r, rep)
	}
	response.JsonOkExit(r, make([]model.DistrictsChildRep, 0))
}
