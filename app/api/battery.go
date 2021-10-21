// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-21
// Based on apiv2 by liasica, magicrolan@qq.com.

package api

import (
    "battery/app/model"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var BatteryApi = new(batteryApi)

type batteryApi struct {
}

func (b *batteryApi) Battery(r *ghttp.Request) {
    response.JsonOkExit(r, []model.BatteryTypeItem{
        {Type: model.BatteryType60, Name: "60伏电池"},
        {Type: model.BatteryType72, Name: "72伏电池"},
    })
}
