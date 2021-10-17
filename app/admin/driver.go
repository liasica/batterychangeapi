// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-16
// Based on apiv2 by liasica, magicrolan@qq.com.

package admin

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/request"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var DriverApi = driverApi{}

type driverApi struct {
}

// Verify
// @Summary 认证列表
// @Tags    管理
// @Accept  json
// @Param   entity body model.UserListReq true "请求参数"
// @Produce  json
// @Router  /admin/driver/verify [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.PackageListItem}}  "返回结果"
func (*driverApi) Verify(r *ghttp.Request) {
    req := new(model.UserListReq)
    _ = request.ParseRequest(r, req)

    total, items := service.UserService.ListUsers(r.Context(), req)
    result := model.ItemsWithTotal{
        Total: total,
    }
    for _, item := range items {
        result.Items = append(result.Items, model.UserVerifyListItem{
            RealName:   item.RealName,
            Mobile:     item.Mobile,
            Type:       item.Type,
            AuthState:  item.AuthState,
            IdCardNo:   item.IdCardNo,
            IdCardImg1: item.IdCardImg1,
            IdCardImg2: item.IdCardImg2,
            IdCardImg3: item.IdCardImg3,
        })
    }
    response.JsonOkExit(r, result)
}
