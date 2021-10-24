// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-24
// Based on apiv2 by liasica, magicrolan@qq.com.

package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
    "fmt"
    "github.com/gogf/gf/database/gdb"
    "sort"
)

var DashboardService = new(dashboardService)

type dashboardService struct {
}

// OpenCities 获取已开通城市
func (*dashboardService) OpenCities(ctx context.Context) (items []model.City) {
    c := dao.Shop.Columns
    _ = dao.Shop.Ctx(ctx).WithAll().Fields(fmt.Sprintf("%s AS id", c.CityId)).Group(c.CityId).Scan(&items)
    for k, item := range items {
        items[k].Name = item.Detail.Name
    }
    return
}

// queryDateBetween 封装时间选择
func (d *dashboardService) queryDateBetween(query *gdb.Model, req *model.DateBetween, field string) *gdb.Model {
    if !req.StartDate.IsZero() {
        query = query.WhereGTE(field, req.StartDate)
    }
    if !req.EndDate.IsZero() {
        query = query.WhereLTE(field, req.EndDate)
    }
    return query
}

// OverviewRiderCount 骑手数量概览
func (d *dashboardService) OverviewRiderCount(ctx context.Context, req *model.DateBetween, data *model.DashboardOverview) {
    // 查询个签用户数量
    c := dao.User.Columns
    query := dao.User.Ctx(ctx)
    _ = d.queryDateBetween(query, req, c.CreatedAt).Fields(`SUM(IF(groupId > 0, 0, 1)) AS personalRiders, SUM(IF(groupId > 0, 1, 0)) AS groupRiders`).Scan(data)
}

// OverviewGroupCount 团队数量
func (d *dashboardService) OverviewGroupCount(ctx context.Context, req *model.DateBetween, data *model.DashboardOverview) {
    c := dao.Group.Columns
    query := dao.Group.Ctx(ctx)
    cnt, _ := d.queryDateBetween(query, req, c.CreatedAt).Count()
    data.GroupCnt = int64(cnt)
}

// OverviewOrderTotal 订单概览
func (d *dashboardService) OverviewOrderTotal(ctx context.Context, req *model.DateBetween, data *model.DashboardOverview) {
    c := dao.ComboOrder.Columns
    query := dao.ComboOrder.Ctx(ctx)
    var orders []*model.ComboOrder
    _ = d.queryDateBetween(query, req, c.CreatedAt).WhereNot(c.Type, model.ComboOrderTypePenalty).Where(c.PayState, model.PayStateSuccess).Scan(&orders)
    for _, order := range orders {
        data.TotalAmount += order.Amount
        data.Deposit += order.Deposit
        data.Orders++
        if order.GroupId == 0 {
            data.PersonalOrders++
            data.PersonalAmount += order.Amount
        } else {
            data.GroupOrders++
            data.GroupAmount += order.Amount
        }
    }
    return
}

// NewlyRiders 新增骑手
func (d *dashboardService) NewlyRiders(ctx context.Context, req *model.DateBetween) (data map[string]model.DashboardNewly) {
    // data = make(map[string]model.DashboardNewly)
    // c := dao.User.Columns
    // query := dao.User.Ctx(ctx)
    // var rows []model.DashboardNewly
    // _ = d.queryDateBetween(query, req, c.CreatedAt).Fields(`DATE(createdAt) AS date, COUNT(1) AS riders`).Scan(&data)
    // for _, row := range rows {
    //     data[row.Date] = row
    // }
    return
}

// NewlyOrders 新增订单
func (d *dashboardService) NewlyOrders(ctx context.Context, req *model.DashboardNewlyReq) (items []model.DashboardNewly) {
    items = make([]model.DashboardNewly, 0)
    c := dao.ComboOrder.Columns
    query := dao.ComboOrder.Ctx(ctx)
    if req.CityId > 0 {
        query = query.Where(c.CityId, req.CityId)
    }
    _ = d.queryDateBetween(query, &req.DateBetween, c.CreatedAt).
        Fields(`DATE(createdAt) AS date, COUNT(1) AS orders, SUM(amount) AS orderAmount, cityId`).
        WhereNot(c.Type, model.ComboOrderTypePenalty).
        Where(c.PayState, model.PayStateSuccess).
        Group("date").
        Scan(&items)
    // 按时间正序排列
    sort.Slice(items, func(i, j int) bool {
        return items[i].Date < items[j].Date
    })
    return
}

// Newly 新增统计
func (d *dashboardService) Newly(ctx context.Context, req *model.DateBetween) (items []model.DashboardNewly) {
    return
}
