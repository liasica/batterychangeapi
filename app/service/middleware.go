// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at https://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-07-14
// Based on apiv2 by liasica, magicrolan@qq.com.

package service

import (
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
    "log"
    "net/http"
)

type middlewareService struct{}

var Middleware = new(middlewareService)

// CORS 全局CORS处理
func (s *middlewareService) CORS(r *ghttp.Request) {
    r.Response.CORSDefault()
    r.Middleware.Next()
}

// ErrorHandle 全局错误处理
func (s *middlewareService) ErrorHandle(r *ghttp.Request) {
    r.Middleware.Next()
    log.Println("服务器故障啦", r.GetError())
    // 服务器故障抛出错误覆盖
    if r.Response.Status >= http.StatusInternalServerError {
        r.Response.Status = http.StatusOK
        r.Response.ClearBuffer()
        // r.Response.Writeln("服务器故障了，请稍后重试或联系管理员！")
        response.Json(r, response.RespCodeSystemError, "服务器故障了，请稍后重试或联系管理员！")
    }
}
