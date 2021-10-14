package user

import (
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "github.com/gogf/gf/net/ghttp"
)

var MessageApi = messageApi{}

type messageApi struct {
}

// List
// @summary 骑手-消息列表
// @description type 为 100， 101， 102， 104 时需要跳转详情页面
// @tags    骑手-消息
// @Accept  json
// @Produce  json
// @Param 	pageIndex query integer  true "当前页码"
// @Param 	pageLimit query integer  true "每页行数"
// @router  /rapi/message [GET]
// @success 200 {object} response.JsonResponse{data=[]model.Message} "返回结果"
func (*messageApi) List(r *ghttp.Request) {
    var page model.Page
    if err := r.Parse(&page); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
    response.JsonOkExit(r, service.MessageService.ListUser(r.Context(), user.Id, page))
}

type MessageReadReg struct {
    MessageIds []uint64 `json:"messageIds"` // 消息ID数组
}

// Read
// @summary 骑手-消息已读标记
// @tags    骑手-消息
// @Accept  json
// @Produce  json
// @Param   entity  body MessageReadReg true "请求数据"
// @router  /rapi/message/read [PUT]
// @success 200 {object} response.JsonResponse{data=[]model.Message} "返回结果"
func (*messageApi) Read(r *ghttp.Request) {
    var req MessageReadReg
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    user := r.Context().Value(model.ContextRiderKey).(*model.ContextRider)
    response.JsonOkExit(r, service.MessageService.Read(r.Context(), user.Id, 1, req.MessageIds))
}

// Detail
// @summary 骑手-消息详情
// @tags    骑手-消息
// @Accept  json
// @Produce  json
// @Param 	id path integer  true "消息ID"
// @router  /rapi/message/:id [GET]
// @success 200 {object} response.JsonResponse{data=model.Message} "返回结果"
func (*messageApi) Detail(r *ghttp.Request) {
    var req model.IdReq
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    msg, _ := service.MessageService.Detail(r.Context(), req.Id)
    response.JsonOkExit(r, msg)
}
