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

var RiderApi = riderApi{}

type riderApi struct {
}

// Verify
// @Summary 认证列表
// @Tags    管理
// @Accept  json
// @Param 	type query integer false "用户类型: 1个签骑手 2团签骑手 3团签BOSS" ENUMS(1,2,3)
// @Param 	authState query integer false "实名认证状态: 0未提交 1待审核 2审核通过 3审核拒绝" ENUMS(0,1,2,3)
// @Param 	realName query string false "骑手姓名"
// @Param 	mobile query string false "骑手电话"
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Produce json
// @Router  /admin/rider/verify [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.UserVerifyListItem}} "返回结果"
func (*riderApi) Verify(r *ghttp.Request) {
    req := new(model.UserVerifyReq)
    _ = request.ParseRequest(r, req)

    total, items := service.UserService.ListVerifyItems(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// Personal
// @Summary 个签用户列表
// @Tags    管理
// @Accept  json
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Param 	groupId query integer false "团签ID"
// @Param 	realName query string false "成员姓名"
// @Param 	mobile query string false "成员电话"
// @Param 	batteryState query integer false "换电状态, 个签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租 5已逾期; 团签骑手换电状态：0未开通 1新签未领 2租借中 3寄存中 4已退租" ENUMS(0,1,2,3,4,5)
// @Param 	startDate query string false "开始日期"
// @Param 	endDate query string false "结束日期"
// @Produce json
// @Router  /admin/rider/personal [GET]
// @Success 200 {object} response.JsonResponse{data=model.ItemsWithTotal{items=[]model.UserListItem}} "返回结果"
func (*riderApi) Personal(r *ghttp.Request) {
    req := new(model.UserListReq)
    _ = request.ParseRequest(r, req)

    total, items := service.UserService.ListPersonalItems(r.Context(), req)
    response.ItemsWithTotal(r, total, items)
}

// Note
// @Summary 创建骑手跟进
// @Tags    管理
// @Accept  json
// @Param   entity body model.UserNotePostReq true "请求参数"
// @Produce json
// @Router  /admin/rider/note [POST]
// @Success 200 {object} response.JsonResponse "返回结果"
func (*riderApi) Note(r *ghttp.Request) {
    req := new(model.UserNotePostReq)
    _ = request.ParseRequest(r, req)
    service.UserNoteService.Create(r.Context(), req)
    response.JsonOkExit(r)
}

// NoteList
// @Summary 获取骑手跟进
// @Tags    管理
// @Accept  json
// @Param   userId path int true "骑手ID"
// @Param 	pageIndex query integer true "当前页码"
// @Param 	pageLimit query integer true "每页行数"
// @Produce json
// @Router  /admin/rider/note/{userId} [GET]
// @Success 200 {object} response.JsonResponse{data=[]model.UserNoteListItem} "返回结果"
func (*riderApi) NoteList(r *ghttp.Request) {
    userId := r.GetUint64("userId")
    page := new(model.Page)
    _ = request.ParseRequest(r, page)
    total, items := service.UserNoteService.List(r.Context(), userId, page)
    response.ItemsWithTotal(r, total, items)
}
