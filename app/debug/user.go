// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-08-02
// Based on apiv2 by liasica, magicrolan@qq.com.

package debug

import (
    "battery/app/dao"
    "battery/app/model"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
)

type userDebug struct {
}

var User = new(userDebug)

func (*userDebug) Reset(r *ghttp.Request) {
    p := r.GetQueryString("phone")
    columns := dao.User.Columns
    data := g.Map{
        columns.GroupId:      0,
        columns.Type:         1,
        columns.BatteryState: 0,
        columns.BatteryType:  0,
        columns.ComboOrderId: 0,
        columns.ComboId:      0,
    }
    auth := r.GetQueryBool("resetauth")

    if auth {
        data[columns.AuthState] = 0
    }

    _, _ = dao.User.Data(data).Where("mobile = ?", p).Update()
}

func (*userDebug) GroupTest(r *ghttp.Request) {
    var orders []model.ComboOrder
    var users []model.User
    userMap := make(map[uint64]model.User)
    _ = dao.ComboOrder.Ctx(r.Context()).Scan(&orders)
    _ = dao.User.Ctx(r.Context()).Scan(&users)

    for _, user := range users {
        userMap[user.Id] = user
    }

    for k, order := range orders {
        orders[k].GroupId = userMap[order.UserId].GroupId
        orders[k].PayState -= 1
    }

    _, _ = dao.ComboOrder.Ctx(r.Context()).Data(orders).FieldsEx("month").Save()
}
