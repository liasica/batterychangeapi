// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-15
// Based on apiv2 by liasica, magicrolan@qq.com.

package debug

import (
    "battery/app/model"
    "battery/app/service"
    "github.com/gogf/gf/net/ghttp"
    "log"
)

type groupDebug struct {
}

var Group = new(groupDebug)

func (*groupDebug) WeekStat(r *ghttp.Request) {
    log.Println(service.GroupDailyStatService.GenerateWeek(r.Context(), 1, model.BatteryType60))
}
