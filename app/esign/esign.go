package esign

import (
	"battery/app/model"
	"battery/app/service"
	"battery/library/response"
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

type SignReq struct {
	Action              string `json:"action"`              //标记该通知的业务类型，该通知固定为：SIGN_FLOW_UPDATE
	FlowId              string `json:"flowId"`              //流程id
	AccountId           string `json:"accountId"`           //签署人的accountId
	AuthorizedAccountId string `json:"authorizedAccountId"` //签约主体的账号id（个人/企业）；如签署人本签署，则返回签署人账号id；如签署人代机构签署，则返回机构账号id 。
	Order               int    `json:"order"`               //签署人的签署顺序
	SignTime            string `json:"signTime"`            //签署时间或拒签时间 格式：yyyy-MM-dd HH:mm:ss
	SignResult          int    `json:"signResult"`          //签署结果 2:签署完成 3:失败 4:拒签
	ThirdOrderNo        string `json:"thirdOrderNo"`        //本次签署任务对应指定的第三方业务流水号id，当存在多个第三方业务流水号id时，返回多个，并逗号隔开该参数取值设置签署区的时候设置的thirdOrderNo参数
	ResultDescription   string `json:"resultDescription"`   //拒签或失败时，附加的原因描述
	Timestamp           int64  `json:"timestamp"`           //时间戳
	ThirdPartyUserId    string `json:"thirdPartyUserId"`    //本次签署任务中对应的签署账号唯一标识，和创建当前签署账号时所传入的thirdPartyUserId值一致
}

// Sign 签约完成回调
func (*callbackApi) Sign(r *ghttp.Request) {
	var req SignReq
	if err := r.Parse(&req); err != nil {
		r.Response.Status = http.StatusBadRequest
		r.Exit()
	}
	if req.Action == "SIGN_FLOW_UPDATE" && req.SignResult == 2 {
		sign, err := service.SignService.GetDetailBayFlowId(r.Context(), req.FlowId)

		if err != nil {
			r.Response.Status = http.StatusInternalServerError
			r.Exit()
		}

		if sign.State != model.SignStateDone {
			if service.SignService.Done(r.Context(), req.FlowId) != nil {
				r.Response.Status = http.StatusInternalServerError
				r.Exit()
			}
		}
	}
	response.JsonOkExit(r)
}
