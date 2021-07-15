package admin

import (
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
	"github.com/gogf/gf/net/ghttp"
	"github.com/shopspring/decimal"
)

var PackagesApi = packagesApi{}

type packagesApi struct {
}

type packageListItem struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Price       float64 `json:"price"`
	BatteryType uint    `json:"batteryType"`
	Days        uint    `json:"days"`
	Earnest     float64 `json:"earnest"`
	CityName    string  `json:"cityName"`
}

func (*packagesApi) List(r *ghttp.Request) {
	var req model.Page
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	total, items := service.PackagesService.ListAdmin(r.Context(), req)
	rep := struct {
		Total int               `json:"total"`
		Items []packageListItem `json:"items"`
	}{
		Total: total,
	}
	if rep.Total > 0 {
		cityIds := make([]uint, len(items))
		for key, packages := range items {
			cityIds[key] = packages.CityId
		}
		cityIdName := service.DistrictsService.MapIdName(r.Context(), cityIds)
		rep.Items = make([]packageListItem, len(items))
		for key, packages := range items {
			rep.Items[key] = packageListItem{
				Id:          packages.Id,
				Name:        packages.Name,
				Amount:      packages.Amount,
				Price:       packages.Price,
				Earnest:     packages.Earnest,
				Days:        packages.Days,
				BatteryType: packages.BatteryType,
				CityName:    cityIdName[packages.CityId],
			}
		}
	}
	response.JsonOkExit(r, rep)
}

type packagesCreateReq struct {
	BatteryType uint    `v:"required|in:60,72" json:"batteryType"`                                         // 60 / 72
	Name        string  `v:"required" json:"name"`                                                         // 名称
	Days        uint    `v:"required|integer|min:1" json:"days"`                                           // 套餐时长天数
	Price       float64 `v:"required|regex:'/(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})$)/'" json:"price"`    // 套餐价格
	Earnest     float64 `v:"required|regex:'/(^[1-9]\d*(\.\d{1,2})?$)|(^0(\.\d{1,2})?$)/'" json:"earnest"` // 保证金
	ProvinceId  uint    `v:"required|integer|min:1" json:"provinceId"`                                     // 省级行政编码
	CityId      uint    `v:"required|integer|min:1" json:"cityId"`                                         // 市级行政编码
}

func (*packagesApi) Create(r *ghttp.Request) {
	var req packagesCreateReq
	if err := r.Parse(&req); err != nil {
		response.Json(r, response.RespCodeArgs, err.Error())
	}
	amount, _ := decimal.NewFromFloat(req.Price).Add(decimal.NewFromFloat(req.Earnest)).Float64()
	if _, err := service.PackagesService.Create(r.Context(), model.Packages{
		Name:        req.Name,
		Type:        1,
		BatteryType: req.BatteryType,
		Amount:      amount,
		Price:       req.Price,
		Earnest:     req.Earnest,
		ProvinceId:  req.ProvinceId,
		Days:        req.Days,
		CityId:      req.CityId,
	}); err != nil {
		response.JsonErrExit(r)
	}
	response.JsonOkExit(r)
}
