package shop

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
    "sort"
)

var BatteryApi = batteryApi{}

type batteryApi struct {
}

// Overview
// @Summary 门店-电池数量概览
// @Tags    门店
// @Accept  json
// @Produce json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
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
// @Param 	startTime query string  false "开始时间"
// @Param 	endTime query string  false "结束时间"
// @Router  /sapi/battery/record  [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.ShopBatteryRecordListWithDateGroup} "返回结果"
func (*batteryApi) Record(r *ghttp.Request) {
    var req model.ShopBatteryRecordListReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    recordList := service.ShopBatteryRecordService.ShopList(
        r.Context(),
        r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId,
        req.Type,
        req.StartDate,
        req.EndDate,
    )
    if len(recordList) == 0 {
        response.JsonOkExit(r, make([]model.ShopBatteryRecordListWithDateGroup, 0))
    }
    var rep []model.ShopBatteryRecordListWithDateGroup
    var dates []string
    tmp := make(map[string]model.ShopBatteryRecordListWithDateGroup)
    layout := "Y-m-d"
    for _, record := range recordList {
        date := record.Date.Format(layout)
        dates = append(dates, date)
        list, ok := tmp[date]
        if !ok {
            list = model.ShopBatteryRecordListWithDateGroup{
                Date:     date,
                InTotal:  0,
                OutTotal: 0,
                Items:    []model.ShopBatteryRecordListItem{},
            }
        }
        switch record.Type {
        case model.ShopBatteryRecordTypeIn:
            list.InTotal++
        case model.ShopBatteryRecordTypeOut:
            list.OutTotal++
        }
        name := record.UserName
        if record.BizType == 0 {
            name = "平台调拨"
        }
        list.Items = append(list.Items, model.ShopBatteryRecordListItem{
            BizType:     record.BizType,
            UserName:    name,
            Num:         record.Num,
            BatteryType: record.BatteryType,
            Date:        record.Date.Format(layout),
        })
        tmp[date] = list
    }

    sort.Strings(dates)
    for _, date := range dates {
        rep = append(rep, tmp[date])
    }

    response.JsonOkExit(r, rep)
}
