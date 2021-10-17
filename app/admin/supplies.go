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
    "github.com/gogf/gf/net/ghttp"
)

type batteryApi struct {
}

var BatteryApi = new(batteryApi)

func (*batteryApi) TransferRecord(r *ghttp.Request) {
}
