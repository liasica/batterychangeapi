package shop

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var BatteryApi = batteryApi{}

type batteryApi struct {
}

// Overview
// @Summary 门店-电池数量概览
// @Tags    门店
// @Accept  json
// @Produce json
// @Router  /sapi/battery [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopBatteryRecordStatRep} "返回结果"
func (*batteryApi) Overview(r *ghttp.Request) {
    mgr := r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager)
    response.JsonOkExit(r, service.ShopBatteryRecordService.GetBatteryNumber(r.Context(), mgr.ShopId))
}

// Record
// @Summary 门店-电池调拨记录
// @Tags    门店
// @Accept  json
// @Produce json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	type query integer true "1入库 2出库"
// @Param 	startDate query string  false "开始时间"
// @Param 	endDate query string  false "结束时间"
// @Router  /sapi/battery/record  [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.ShopBatteryRecordListWithDateGroup} "返回结果"
func (*batteryApi) Record(r *ghttp.Request) {
    req := new(model.ShopBatteryRecordListReq)
    _ = request.ParseRequest(r, req)
    req.ShopId = r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId
    response.JsonOkExit(r, service.ShopBatteryRecordService.RecordShopFilter(r.Context(), req))
}
