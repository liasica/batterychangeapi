// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-17
// Based on apiv2 by liasica, magicrolan@qq.com.

package admin

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
)

type batteryApi struct {
}

var BatteryApi = new(batteryApi)

// TransferRecord
// @Summary 电池调拨记录
// @Tags    管理
// @Accept  json
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Param 	shopId query integer false "门店ID"
// @Param 	type query integer true "1出库 2入库"
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Produce json
// @Router  /admin/battery/record [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.BatteryRecordListItem}} "返回结果"
func (*batteryApi) TransferRecord(r *ghttp.Request) {
    var req = new(model.BatteryRecordListReq)
    _ = request.ParseRequest(r, req)
    total, items := service.ShopBatteryRecordService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// Allocate
// @Summary 电池调拨
// @Tags    管理
// @Accept  json
// @Param   entity body model.BatteryAllocateReq true "请求参数"
// @Param 	batteryType query string true "电池型号" ENUMS(60,72)
// @Param 	from query integer true "调出自 0平台 其他店铺ID"
// @Param 	to query integer true "调入至 0平台 其他店铺ID"
// @Param 	num query integer true "数量"
// @Produce json
// @Router  /admin/battery/record [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*batteryApi) Allocate(r *ghttp.Request) {
    var req = new(model.BatteryAllocateReq)
    _ = request.ParseRequest(r, req)
    if err := service.ShopBatteryRecordService.Allocate(r.Context(), req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    response.JsonOkExit(r)
}

// Exception
// @Summary 电池异常记录
// @Tags    管理
// @Accept  json
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Param 	shopId query integer false "门店ID"
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Produce json
// @Router  /admin/battery/exception [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.ExceptionListItem}} "返回结果"
func (*batteryApi) Exception(r *ghttp.Request) {
    var req = new(model.ExceptionListReq)
    _ = request.ParseRequest(r, req)
    total, items := service.ExceptionService.PageList(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// ExceptionFix
// @Summary 处理电池异常
// @Tags    管理
// @Accept  json
// @Param   id path int true "记录ID"
// @Produce json
// @Router  /admin/battery/exception/{id} [PUT]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*batteryApi) ExceptionFix(r *ghttp.Request) {
    id := r.GetInt("id")
    _, _ = dao.Exception.Where("id = ?", id).Data(g.Map{"state": model.ExceptionStateProcessed}).Update()
    response.JsonOkExit(r)
}
