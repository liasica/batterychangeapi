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

// List
// @Summary 套餐列表
// @Tags    管理
// @Accept  json
// @Param   entity body model.Page true "分页参数"
// @Produce  json
// @Router  /admin/package [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.PackageListItem}}  "返回结果"
func (*packagesApi) List(r *ghttp.Request) {
    var req model.Page
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    total, items := service.PackagesService.ListAdmin(r.Context(), req)
    rep := struct {
        Total int                     `json:"total"`
        Items []model.PackageListItem `json:"items"`
    }{
        Total: total,
    }
    if rep.Total > 0 {
        rep.Items = make([]model.PackageListItem, len(items))
        for key, packages := range items {
            rep.Items[key] = model.PackageListItem{
                Id:          packages.Id,
                Name:        packages.Name,
                Amount:      packages.Amount,
                Price:       packages.Price,
                Earnest:     packages.Earnest,
                Days:        packages.Days,
                BatteryType: packages.BatteryType,
                CityId:      packages.CityId,
                ProvinceId:  packages.ProvinceId,
            }
        }
    }
    response.JsonOkExit(r, rep)
}

// Create
// @Summary 创建套餐
// @Tags    管理
// @Accept  json
// @Param   entity body model.PackageCreateReq true "门店详情"
// @Produce  json
// @Router  /admin/package [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*packagesApi) Create(r *ghttp.Request) {
    var req model.PackageCreateReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    amount, _ := decimal.NewFromFloat(req.Price).Add(decimal.NewFromFloat(req.Earnest)).Float64()
    if _, err := service.PackagesService.Create(r.Context(), model.Packages{
        Name:        req.Name,
        Type:        req.Type,
        BatteryType: req.BatteryType,
        Amount:      amount,
        Price:       req.Price,
        Earnest:     req.Earnest,
        ProvinceId:  req.ProvinceId,
        Days:        req.Days,
        CityId:      req.CityId,
        Desc:        req.Desc,
    }); err != nil {
        response.JsonErrExit(r)
    }
    response.JsonOkExit(r)
}

// Edit
// @Summary 编辑套餐
// @Tags    管理
// @Accept  json
// @Param   id path int true "套餐ID"
// @Param   entity body model.PackageCreateReq true "门店详情"
// @Produce  json
// @Router  /admin/package [PUT]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*packagesApi) Edit(r *ghttp.Request) {

}
