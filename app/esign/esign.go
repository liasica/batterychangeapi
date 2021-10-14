package esign

import (
    "battery/app/dao"
    "battery/app/model"
    "battery/app/service"
    "battery/library/response"
    "context"
    "github.com/gogf/gf/database/gdb"
    "github.com/gogf/gf/net/ghttp"
    "net/http"
)

var CallbackApi = callbackApi{}

type callbackApi struct {
}

type realNameReq struct {
    FlowId     string `json:"flowId"`
    AccountId  string `json:"accountId"`
    Success    bool   `json:"success"`
    ContextId  string `json:"contextId"`
    Verifycode string `json:"verifycode"`
}

func (*callbackApi) RealName(r *ghttp.Request) {
    var req realNameReq
    if err := r.Parse(&req); err == nil {
        var state uint
        if req.Success {
            state = model.AuthStateVerifySuccess
        } else {
            state = model.AuthStateDefaultFailed
        }
        if service.UserService.RealNameAuthVerifyCallBack(r.Context(), req.AccountId, model.RealNameAuthVerifyReq{
            AuthState: state,
        }) == nil {
            r.Response.Status = http.StatusOK
            r.Exit()
        }
    }
    r.Response.Status = http.StatusInternalServerError
    r.Exit()
}

// SignReq 签署人签署完成回调通知结构
// https://open.esign.cn/doc/detail?id=opendoc%2Fpaas_api%2Fawifwu&namespace=opendoc%2Fpaas_api
type SignReq struct {
    Action              string `json:"action"`              // 标记该通知的业务类型，该通知固定为：SIGN_FLOW_UPDATE
    FlowId              string `json:"flowId"`              // 流程id
    AccountId           string `json:"accountId"`           // 签署人的accountId
    AuthorizedAccountId string `json:"authorizedAccountId"` // 签约主体的账号id（个人/企业）；如签署人本签署，则返回签署人账号id；如签署人代机构签署，则返回机构账号id 。
    Order               int    `json:"order"`               // 签署人的签署顺序
    SignTime            string `json:"signTime"`            // 签署时间或拒签时间 格式：yyyy-MM-dd HH:mm:ss
    SignResult          int    `json:"signResult"`          // 签署结果 2:签署完成 3:失败 4:拒签
    ThirdOrderNo        string `json:"thirdOrderNo"`        // 本次签署任务对应指定的第三方业务流水号id，当存在多个第三方业务流水号id时，返回多个，并逗号隔开该参数取值设置签署区的时候设置的thirdOrderNo参数
    ResultDescription   string `json:"resultDescription"`   // 拒签或失败时，附加的原因描述
    Timestamp           int64  `json:"timestamp"`           // 时间戳
    ThirdPartyUserId    string `json:"thirdPartyUserId"`    // 本次签署任务中对应的签署账号唯一标识，和创建当前签署账号时所传入的thirdPartyUserId值一致
}

// SignFinishReq 流程结束回调通知
// https://open.esign.cn/doc/detail?id=opendoc%2Fpaas_api%2Fld5udg&namespace=opendoc%2Fpaas_api
type SignFinishReq struct {
    Action            string `json:"action"`            // 标记该通知的业务类型，该通知固定为：SIGN_FLOW_FINISH
    FlowId            string `json:"flowId"`            // 流程id
    BusinessScence    string `json:"businessScence"`    // 签署文件主题描述
    FlowStatus        string `json:"flowStatus"`        // 任务状态 2-已完成: 所有签署人完成签署； 3-已撤销: 发起方撤销签署任务； 5-已过期: 签署截止日到期后触发； 7-已拒签
    StatusDescription string `json:"statusDescription"` // 当流程异常结束时，附加终止原因描述
    CreateTime        string `json:"createTime"`        // 签署任务发起时间 格式yyyy-MM-dd HH:mm:ss
    EndTime           string `json:"endTime"`           // 签署任务结束时间 格式yyyy-MM-dd HH:mm:ss
    Timestamp         int64  `json:"timestamp"`         // 时间戳
}

// Sign 签约完成回调
func (*callbackApi) Sign(r *ghttp.Request) {
    var req SignFinishReq
    if err := r.Parse(&req); err != nil {
        r.Response.Status = http.StatusBadRequest
        r.Exit()
    }
    // 签署人签署完成回调通知
    if req.Action == "SIGN_FLOW_FINISH" && req.FlowStatus == "2" {
        sign, err := service.SignService.GetDetailBayFlowId(r.Context(), req.FlowId)
        if err != nil {
            r.Response.Status = http.StatusInternalServerError
            r.Exit()
        }
        if sign.State != model.SignStateDone {
            if dao.Sign.DB.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
                if err := service.SignService.Done(ctx, req.FlowId); err != nil {
                    return err
                }
                if sign.GroupId > 0 {
                    if err := service.UserService.GroupUserSignDone(ctx, sign); err != nil {
                        return err
                    }
                }
                return nil
            }) != nil {
                r.Response.Status = http.StatusInternalServerError
                r.Exit()
            }
        }
    }

    response.JsonOkExit(r)
}

// SignState 查询签约结果
// @Summary 获取签约结果
// @Tags    公用
// @Accept  json
// @Produce  json
// @Router  /esign/state/:fileId [GET]
// @Success 200 {object} response.JsonResponse{data=int}  "返回结果"
func (*callbackApi) SignState(r *ghttp.Request) {
    fileId := r.Get("fileId").(string)

    s, err := service.SignService.GetDetailBayFileId(r.Context(), fileId)
    if err != nil {
        r.Response.Status = http.StatusInternalServerError
        r.Exit()
    }

    response.JsonOkExit(r, s.State)
}
