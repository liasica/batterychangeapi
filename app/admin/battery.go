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
// @Param   entity body model.BatteryRecordListReq true "请求参数"
// @Produce json
// @Router  /admin/battery/record [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.BatteryRecordListItem}}  "返回结果"
func (*batteryApi) TransferRecord(r *ghttp.Request) {
    var req = new(model.BatteryRecordListReq)
    _ = request.ParseRequest(r, req)
    total, items := service.ShopBatteryRecordService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// Exception
// @Summary 电池异常记录
// @Tags    管理
// @Accept  json
// @Param   entity body model.ExceptionListReq true "请求参数"
// @Produce json
// @Router  /admin/battery/exception [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.ExceptionListItem}}  "返回结果"
func (*batteryApi) Exception(r *ghttp.Request) {
    var req = new(model.ExceptionListReq)
    _ = request.ParseRequest(r, req)
    g.Dump(req)
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
// @Success 200 {object} response.JsonResponse  "返回结果"
func (*batteryApi) ExceptionFix(r *ghttp.Request) {
    id := r.GetInt("id")
    _, _ = dao.Exception.Where("id = ?", id).Data(g.Map{"state": model.ExceptionStateProcessed}).Update()
    response.JsonOkExit(r)
}
