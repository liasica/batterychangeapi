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

// PackageOrderNewSuccessCallback 新购套餐成功回调
func (*alipayApi) PackageOrderNewSuccessCallback(r *ghttp.Request) {
	res, err := alipay.Service().GetTradeNotification(r.Context(), r.Request)
	if err != nil || res == nil {
		g.Log().Error(err)
		r.Response.Status = http.StatusBadRequest
		r.Response.Write("error")
		r.Exit()
	}
	if res.TradeStatus == alipayV3.TradeStatusSuccess {
		//TODO 查询校验
		if packageOrderNewSuccess(r.Context(), gtime.New(res.GmtPayment), res.OutTradeNo, res.TradeNo, model.PayTypeAliPay) != nil {
			r.Response.Status = http.StatusInternalServerError
			r.Response.Write("error")
			r.Exit()
		}
	}
	r.Response.Status = http.StatusOK
	r.Response.Write("success")
	r.Exit()
}

// PackageOrderRenewalSuccessCallback 续购套餐成功回调
func (*alipayApi) PackageOrderRenewalSuccessCallback(r *ghttp.Request) {
	res, err := alipay.Service().GetTradeNotification(r.Context(), r.Request)
	if err != nil || res == nil {
		g.Log().Error(err)
		r.Response.Status = http.StatusBadRequest
		r.Response.Write("error")
		r.Exit()
	}
	if res.TradeStatus == alipayV3.TradeStatusSuccess {
		//TODO 查询校验
		if packageOrderRenewalSuccess(r.Context(), gtime.New(res.GmtPayment), res.OutTradeNo, res.TradeNo, model.PayTypeAliPay) != nil {
			r.Response.Status = http.StatusInternalServerError
			r.Response.Write("error")
			r.Exit()
		}
	}
	r.Response.Status = http.StatusOK
	r.Response.Write("success")
	r.Exit()
}

// PackageOrderPenaltySuccessCallback 违约金支付成功回调
func (*alipayApi) PackageOrderPenaltySuccessCallback(r *ghttp.Request) {
	res, err := alipay.Service().GetTradeNotification(r.Context(), r.Request)
	if err != nil || res == nil {
		g.Log().Error(err)
		r.Response.Status = http.StatusBadRequest
		r.Response.Write("error")
		r.Exit()
		return
	}
	if res.TradeStatus == alipayV3.TradeStatusSuccess {
		//TODO 查询校验
		if packageOrderPenaltySuccess(r.Context(), gtime.New(res.GmtPayment), res.OutTradeNo, res.TradeNo, model.PayTypeAliPay) != nil {
			r.Response.Status = http.StatusInternalServerError
			r.Response.Write("error")
			r.Exit()
		}
	}
	r.Response.Status = http.StatusOK
	r.Response.Write("success")
	r.Exit()
}
