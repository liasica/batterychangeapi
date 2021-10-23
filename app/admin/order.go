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
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

type orderApi struct {
}

var OrderApi = new(orderApi)

// List
// @Summary 订单列表
// @Tags    管理
// @Accept  json
// @Produce json
// @Param 	no query string true "订单编号"
// @Param 	shopId query integer false "门店ID"
// @Param 	userId query integer false "骑手ID"
// @Param 	comboId query integer false "套餐ID"
// @Param 	cityId query integer false "城市ID"
// @Param 	realName query string false "骑手姓名"
// @Param 	mobile query string false "骑手电话"
// @Param 	type query integer false "订单类型 1新签 2续签 3违约金" ENUMS(1,2,3)
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Router  /admin/order [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.OrderListItem}} "返回结果"
func (*orderApi) List(r *ghttp.Request) {
    req := new(model.OrderListReq)
    _ = request.ParseRequest(r, req)

    total, items := service.ComboOrderService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}
