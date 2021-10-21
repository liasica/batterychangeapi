package shop

import (
    "github.com/gogf/gf/net/ghttp"

    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
)

var AssetApi = assetApi{}

type assetApi struct {
}

// BatteryStat
// @Summary 店长-资产-电池统计
// @Tags    店长-资产
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Router  /sapi/asset/battery_stat [GET]
// @Success 200 {object} response.JsonResponse{data=model.ShopBatteryRecordStatRep} "返回结果"
func (*assetApi) BatteryStat(r *ghttp.Request) {
    shop, _ := service.ShopService.Detail(r.Context(), r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId)
    response.JsonOkExit(r, model.ShopBatteryRecordStatRep{
        InTotal:  shop.BatteryInCnt60 + shop.BatteryInCnt72,
        OutTotal: shop.BatteryOutCnt60 + shop.BatteryOutCnt72,
    })
}

// BatteryList
// @Summary 店长-资产-电池出入库列表
// @Tags    店长-资产
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @Param 	type query integer  true "1 入库 2 出库"
// @Param 	startTime query string  false "开始时间"
// @Param 	endTime query string  false "结束时间"
// @Router  /sapi/asset/battery_list  [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.ShopBatteryRecordListRep} "返回结果"
func (*assetApi) BatteryList(r *ghttp.Request) {
    var req model.ShopBatteryRecordListReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    recordList := service.ShopBatteryRecordService.ShopList(r.Context(),
        r.Context().Value(model.ContextShopManagerKey).(*model.ContextShopManager).ShopId, req.Type,
        req.StartDate,
        req.EndDate,
    )
    if len(recordList) == 0 {
        response.JsonOkExit(r, make([]model.ShopBatteryRecordListRep, 0))
    }
    rep := make([]model.ShopBatteryRecordListRep, len(recordList))
    days := make([]int, len(recordList))
    for key, record := range recordList {
        days[key] = record.Day
    }
    // 按天统计
    daysTotal := service.ShopBatteryRecordService.ShopDaysTotal(r.Context(), days, req.Type)
    daysTotalCnt := make(map[int]uint)
    for _, day := range daysTotal {
        daysTotalCnt[day.Day] = day.Cnt
    }
    for key, record := range recordList {
        rep[key] = model.ShopBatteryRecordListRep{
            BizType:     record.BizType,
            UserName:    record.UserName,
            Num:         record.Num,
            BatteryType: record.BatteryType,
            At:          record.CreatedAt,
            DayCnt:      daysTotalCnt[record.Day],
        }
    }
    response.JsonOkExit(r, rep)
}
