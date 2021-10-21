package api

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/afs"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var SmsApi = smsApi{}

type smsApi struct {
}

// Send
// @Summary 公用-发送短信
// @Tags    公用
// @Accept  json
// @Produce json
// @Param   entity  body model.SmsSendReq true "注册数据"
// @Router  /api/sms [POST]
// @Success 200 {object} response.JsonResponse  "返回结果"
func (*smsApi) Send(r *ghttp.Request) {
    var req model.SmsSendReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    if afs.Service().Verify(req.Sig, req.SessionId, req.Token, r.GetRemoteIp()) != true {
        response.Json(r, response.RespCodeArgs, "请稍后再试！")
    }
    if err := service.SmsServer.Send(r.Context(), req); err != nil {
        response.JsonErrExit(r, response.RespCodeSystemError)
    }
    response.JsonOkExit(r)
}
