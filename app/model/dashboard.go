// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-24
// Based on apiv2 by liasica, magicrolan@qq.com.

package model

// DashboardOverview 订单概览
type DashboardOverview struct {
    PersonalRiders int64   `json:"personalRiders"` // 个签用户数量
    GroupRiders    int64   `json:"groupRiders"`    // 团签用户数量
    GroupCnt       int64   `json:"groupCnt"`       // 团队数量
    TotalAmount    float64 `json:"totalAmount"`    // 订单金额总计
    Orders         int64   `json:"orders"`         // 总订单量
    PersonalAmount float64 `json:"personalAmount"` // 个签订单金额总计
    PersonalOrders int64   `json:"personalOrders"` // 个签订单数量
    Deposit        float64 `json:"deposit"`        // 个签押金
    GroupAmount    float64 `json:"groupAmount"`    // 团签订单金额总计
    GroupOrders    int64   `json:"groupOrders"`    // 团签订单数量
}

// DashboardNewlyReq 新增统计请求体
type DashboardNewlyReq struct {
    DateBetween
    CityId uint `json:"cityId"` // 城市ID
}

// DashboardOrderNewly 新增订单统计
type DashboardOrderNewly struct {
    Date    string  `json:"date"`    // 日期
    New     int64   `json:"new"`     // 新增订单
    Renewal int64   `json:"renewal"` // 续签订单
    Amount  float64 `json:"amount"`  // 新增订单金额

    // Riders      int64   `json:"riders"`      // 新增骑手
}
