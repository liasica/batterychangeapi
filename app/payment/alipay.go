package payment

import (
    "battery/app/model"
    "battery/library/payment/alipay"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "github.com/gogf/gf/os/gtime"
    alipayV3 "github.com/smartwalle/alipay/v3"
    "net/http"
)

var AlipayApi = alipayApi{}

type alipayApi struct {
}

// ComboOrderNewSuccessCallback 新购套餐成功回调
func (*alipayApi) ComboOrderNewSuccessCallback(r *ghttp.Request) {
    res, err := alipay.Service().GetTradeNotification(r.Context(), r.Request)
    if err != nil || res == nil {
        g.Log().Error(err)
        r.Response.Status = http.StatusBadRequest
        r.Response.Write("error")
        r.Exit()
    }
    if res.TradeStatus == alipayV3.TradeStatusSuccess {
        if comboOrderNewSuccess(r.Context(), gtime.New(res.GmtPayment), res.OutTradeNo, res.TradeNo, model.PayTypeAliPay) != nil {
            r.Response.Status = http.StatusInternalServerError
            r.Response.Write("error")
            r.Exit()
        }
    }
    r.Response.Status = http.StatusOK
    r.Response.Write("success")
    r.Exit()
}

// ComboOrderRenewalSuccessCallback 续购套餐成功回调
func (*alipayApi) ComboOrderRenewalSuccessCallback(r *ghttp.Request) {
    res, err := alipay.Service().GetTradeNotification(r.Context(), r.Request)
    if err != nil || res == nil {
        g.Log().Error(err)
        r.Response.Status = http.StatusBadRequest
        r.Response.Write("error")
        r.Exit()
    }
    if res.TradeStatus == alipayV3.TradeStatusSuccess {
        if comboOrderRenewalSuccess(r.Context(), gtime.New(res.GmtPayment), res.OutTradeNo, res.TradeNo, model.PayTypeAliPay) != nil {
            r.Response.Status = http.StatusInternalServerError
            r.Response.Write("error")
            r.Exit()
        }
    }
    r.Response.Status = http.StatusOK
    r.Response.Write("success")
    r.Exit()
}

// ComboOrderPenaltySuccessCallback 违约金支付成功回调
func (*alipayApi) ComboOrderPenaltySuccessCallback(r *ghttp.Request) {
    res, err := alipay.Service().GetTradeNotification(r.Context(), r.Request)
    if err != nil || res == nil {
        g.Log().Error(err)
        r.Response.Status = http.StatusBadRequest
        r.Response.Write("error")
        r.Exit()
        return
    }
    if res.TradeStatus == alipayV3.TradeStatusSuccess {
        if comboOrderPenaltySuccess(r.Context(), gtime.New(res.GmtPayment), res.OutTradeNo, res.TradeNo, model.PayTypeAliPay) != nil {
            r.Response.Status = http.StatusInternalServerError
            r.Response.Write("error")
            r.Exit()
        }
    }
    r.Response.Status = http.StatusOK
    r.Response.Write("success")
    r.Exit()
}
