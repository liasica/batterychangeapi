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
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "github.com/gogf/gf/os/gtime"
    "github.com/shopspring/decimal"
)

type groupDebug struct {
}

var Group = new(groupDebug)

func (*groupDebug) Settlement(r *ghttp.Request) {
    g.Dump(gtime.Now().AddDate(0, 0, -1).Format("Y-m-d"))
    g.Dump(gtime.Now().After(gtime.NewFromStr("2021-10-19").AddDate(0, 0, 1)))
    g.Dump(decimal.NewFromFloat(23.56333).Mul(decimal.NewFromInt(int64(182))).Round(2).Float64())
}
