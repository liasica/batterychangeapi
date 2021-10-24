// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-24
// Based on apiv2 by liasica, magicrolan@qq.com.

package admin

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var DashboardApi = new(dashboardApi)

type dashboardApi struct {
}

// OpenCities
// @Summary 已开通城市
// @Tags    管理
// @Accept  json
// @Produce json
// @Router  /admin/dashboard/cities [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.City} "返回结果"
func (d *dashboardApi) OpenCities(r *ghttp.Request) {
    response.JsonOkExit(r, service.DashboardService.OpenCities(r.Context()))
}

// Overview
// @Summary 统计概览
// @Tags    管理
// @Accept  json
// @Produce json
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Router  /admin/dashboard/overview [GET]
// @Success 200 {object} response.JsonResponse{data=model.DashboardOverview} "返回结果"
func (*dashboardApi) Overview(r *ghttp.Request) {
    var req = new(model.DateBetween)
    _ = request.ParseRequest(r, req)
    ctx := r.Context()
    data := new(model.DashboardOverview)
    service.DashboardService.OverviewRiderCount(ctx, req, data)
    service.DashboardService.OverviewGroupCount(ctx, req, data)
    service.DashboardService.OverviewOrderTotal(ctx, req, data)
    response.JsonOkExit(r, data)
}
