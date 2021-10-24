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

// Newly
// @Summary 新增订单统计
// @Tags    管理
// @Accept  json
// @Produce json
// @Param 	cityId query int false "城市ID"
// @Param 	startDate query string true "开始日期"
// @Param 	endDate query string true "结束日期"
// @Router  /admin/dashboard/newly [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.DashboardOrderNewly} "返回结果"
func (*dashboardApi) Newly(r *ghttp.Request) {
    var req = new(model.DashboardNewlyReq)
    _ = request.ParseRequest(r, req)
    if req.StartDate.AddDate(0, 0, 60).Before(req.EndDate) {
        response.JsonErrExit(r, response.RespCodeArgs, "时间范围太大")
    }
    ctx := r.Context()
    data := service.DashboardService.NewlyOrders(ctx, req)
    response.JsonOkExit(r, data)
}

// Business
// @Summary 业务统计
// @Tags    管理
// @Accept  json
// @Produce json
// @Param 	cityId query int false "城市ID"
// @Param 	startDate query string true "开始日期"
// @Param 	endDate query string true "结束日期"
// @Router  /admin/dashboard/business [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.DashboardBusiness} "返回结果"
func (*dashboardApi) Business(r *ghttp.Request) {
    var req = new(model.DashboardBusinessReq)
    _ = request.ParseRequest(r, req)
    if req.StartDate.AddDate(0, 0, 60).Before(req.EndDate) {
        response.JsonErrExit(r, response.RespCodeArgs, "时间范围太大")
    }
    ctx := r.Context()
    data := service.DashboardService.Business(ctx, req)
    response.JsonOkExit(r, data)
}
