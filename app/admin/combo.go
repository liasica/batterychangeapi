package admin

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
    "github.com/shopspring/decimal"
)

var ComboApi = comboApi{}

type comboApi struct {
}

// List
// @Summary 套餐列表
// @Tags    管理
// @Accept  json
// @Param   entity body model.Page true "分页参数"
// @Produce  json
// @Router  /admin/combo [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.ComboListItem}}  "返回结果"
func (*comboApi) List(r *ghttp.Request) {
    var req model.Page
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    total, items := service.ComboService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// Create
// @Summary 创建套餐
// @Tags    管理
// @Accept  json
// @Param   entity body model.ComboReq true "门店详情"
// @Produce  json
// @Router  /admin/combo [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*comboApi) Create(r *ghttp.Request) {
    var req model.ComboReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    amount, _ := decimal.NewFromFloat(req.Price).Add(decimal.NewFromFloat(req.Deposit)).Float64()
    if _, err := service.ComboService.Create(r.Context(), model.Combo{
        Name:        req.Name,
        Type:        req.Type,
        BatteryType: req.BatteryType,
        Amount:      amount,
        Price:       req.Price,
        Deposit:     req.Deposit,
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
// @Param   entity body model.ComboReq true "门店详情"
// @Produce  json
// @Router  /admin/combo/{id} [PUT]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*comboApi) Edit(r *ghttp.Request) {
    var req model.ComboReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if req.Disable {
        _, _ = dao.Combo.Where("id", r.GetInt("id")).Delete()
    } else {
        data := r.GetMap()
        _, a1 := data["amount"]
        _, e1 := data["deposit"]
        _, p1 := data["price"]
        if a1 || e1 || p1 {
            response.Json(r, response.RespCodeArgs, "套餐价格禁止修改")
        }
        data["deletedAt"] = nil
        _, _ = dao.Combo.Data(data).Where("id", r.GetInt("id")).Unscoped().Update()
    }
    response.JsonOkExit(r)
}
