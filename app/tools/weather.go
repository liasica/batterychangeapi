// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-08-05
// Based on apiv2 by liasica, magicrolan@qq.com.

package tools

import (
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

type weather struct{}

var Weather = new(weather)

// Now
// @summary 工具-获取天气
// @tags    公用
// @Accept  json
// @Produce  json
// @param   lng query string true "经度"
// @param   lat query string true "纬度"
// @router  /tools/weather [GET]
// @success 200 {object} response.JsonResponse{data=service.WeatherNow}  "返回结果"
func (*weather) Now(r *ghttp.Request) {
    lng := r.GetQueryString("lng")
    lat := r.GetQueryString("lat")
    w, err := service.CommonService.WeatherNow(lng, lat)
    if err != nil {
        response.JsonErrExit(r)
        return
    }
    response.JsonOkExit(r, w)
}
