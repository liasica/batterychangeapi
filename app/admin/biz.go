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

var BizApi = new(bizApi)

type bizApi struct {
}

// List
// @Summary 业务列表
// @Tags    管理
// @Accept  json
// @Produce json
// @Param   entity body model.BizListReq true "门店列表请求"
// @Router  /admin/biz [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.BizEntity}}  "返回结果"
func (*bizApi) List(r *ghttp.Request) {
    var req = new(model.BizListReq)
    _ = request.ParseRequest(r, req)

    total, items := service.UserBizService.ListAdmin(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}
