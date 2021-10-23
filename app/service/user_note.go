// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-22
// Based on apiv2 by liasica, magicrolan@qq.com.

package service

import (
    "battery/app/dao"
    "battery/app/model"
    "context"
)

var UserNoteService = new(userNoteService)

type userNoteService struct {
}

// Create 创建骑手跟进信息
func (*userNoteService) Create(ctx context.Context, req *model.UserNotePostReq) {
    sysUser := ctx.Value(model.ContextAdminKey).(*model.ContextAdmin)
    _, _ = dao.UserNote.Ctx(ctx).Data(model.UserNote{
        UserId:      req.UserId,
        Content:     req.Content,
        SysUserId:   sysUser.Id,
        SysUserName: sysUser.Username,
    }).Save()
}

// List 获取骑手跟进列表
func (*userNoteService) List(ctx context.Context, userId uint64, page *model.Page) (total int, items []model.UserNoteListItem) {
    query := dao.UserNote.Ctx(ctx).OrderDesc(dao.UserNote.Columns.CreatedAt).Where(dao.UserNote.Columns.UserId, userId)
    total, _ = query.Count()
    _ = query.Page(page.PageIndex, page.PageLimit).Scan(&items)
    return
}
