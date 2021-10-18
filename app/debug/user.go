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
    "battery/app/service"
    "battery/library/snowflake"
    "encoding/base64"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "log"
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
    mobile := "18501358302"
    user := model.User{
        GroupId:  2,
        RealName: "ContactName",
        Mobile:   mobile,
        Type:     model.UserTypeGroupBoss,
    }
    log.Println(base64.StdEncoding.EncodeToString(snowflake.Service().Generate().Bytes()))
    log.Println(service.UserService.AddOrSetGroupUser(r.Context(), user))
}
