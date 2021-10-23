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
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Param 	shopId query integer false "门店ID"
// @Param 	userId query integer false "骑手ID"
// @Param 	shopId query integer false "店铺ID"
// @Param 	realName query string false "骑手姓名"
// @Param 	mobile query string false "手机号"
// @Param 	userTpe query integer true "用户类型: 1个签 2团签" ENUMS(1,2)
// @Param 	bizTpe query integer true "业务类型: 1新签 2换电 3寄存 4恢复计费 5退租" ENUMS(1,2,3,4,5)
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Router  /admin/biz [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.BizEntity}} "返回结果"
func (*bizApi) List(r *ghttp.Request) {
    var req = new(model.BizListReq)
    _ = request.ParseRequest(r, req)

    total, items := service.UserBizService.Filter(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}
