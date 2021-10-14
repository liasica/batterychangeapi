package shop

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var ExceptionApi = exceptionApi{}

type exceptionApi struct {
}

// Report
// @summary 店长-异常上报
// @tags    店长
// @Accept  json
// @Produce  json
// @Param   entity  body model.ExceptionReportReq true "请求数据"
// @router  /sapi/exception [POST]
// @success 200 {object} response.JsonResponse "返回结果"
func (*exceptionApi) Report(r *ghttp.Request) {
    var req model.ExceptionReportReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if service.ExceptionService.Create(r.Context(), req) == nil {
        response.JsonOkExit(r)
    }
    response.JsonErrExit(r)
}
