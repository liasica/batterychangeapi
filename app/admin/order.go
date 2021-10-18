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
// @Param   entity body model.OrderListReq true "门店列表请求"
// @Router  /admin/order [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.OrderListItem}}  "返回结果"
func (*orderApi) List(r *ghttp.Request) {
    req := new(model.OrderListReq)
    _ = request.ParseRequest(r, req)

    total, items := service.ComboOrderService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}
