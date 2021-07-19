package payment

import (
	"battery/app/model"
	"battery/library/payment/wechat"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"net/http"
)

var WechatApi = wechatApi{}

type wechatApi struct {
}

// PackageOrderNewSuccessCallback 新购套餐支付成功回调
func (*wechatApi) PackageOrderNewSuccessCallback(r *ghttp.Request) {
	var content wechatPaySuccessNotifyContent
	if _, err := wechat.Service().ParseNotify(r.Context(), r.Request, &content); err != nil {
		g.Log().Error(err.Error())
		r.Response.Status = http.StatusBadRequest
		r.Exit()
	}

	if content.TradeState == "SUCCESS" {
		//TODO 查询校验
		if err := packageOrderNewSuccess(r.Context(), content.SuccessTime, content.OutTradeNo, content.TransactionId, model.PayTypeWechat); err != nil {
			g.Log().Error(err.Error())
			r.Response.Status = http.StatusInternalServerError
			r.Exit()
		}
	}
	r.Response.Status = http.StatusOK
	r.Response.Write(wechatPaySuccessResponse{
		Code:    "SUCCESS",
		Message: "成功",
	})
	r.Exit()
}

// PackageOrderRenewalSuccessCallback 续购套餐支付成功回调
func (*wechatApi) PackageOrderRenewalSuccessCallback(r *ghttp.Request) {
	var content wechatPaySuccessNotifyContent
	if _, err := wechat.Service().ParseNotify(r.Context(), r.Request, &content); err != nil {
		g.Log().Error(err.Error())
		r.Response.Status = http.StatusBadRequest
		r.Exit()
	}
	if content.TradeState == "SUCCESS" {
		//TODO 查询校验
		if err := packageOrderRenewalSuccess(r.Context(), content.SuccessTime, content.OutTradeNo, content.TransactionId, model.PayTypeWechat); err != nil {
			g.Log().Error(err.Error())
			r.Response.Status = http.StatusInternalServerError
			r.Exit()
		}
	}
	r.Response.Status = http.StatusOK
	r.Response.Write(wechatPaySuccessResponse{
		Code:    "SUCCESS",
		Message: "成功",
	})
	r.Exit()
}

// PackageOrderPenaltySuccessCallback 续购套餐支付成功回调
func (*wechatApi) PackageOrderPenaltySuccessCallback(r *ghttp.Request) {
	var content wechatPaySuccessNotifyContent
	if _, err := wechat.Service().ParseNotify(r.Context(), r.Request, &content); err != nil {
		r.Response.Status = http.StatusBadRequest
		r.Exit()
	}
	if content.TradeState == "SUCCESS" {
		//TODO 查询校验
		if err := packageOrderPenaltySuccess(r.Context(), content.SuccessTime, content.OutTradeNo, content.TransactionId, model.PayTypeWechat); err != nil {
			r.Response.Status = http.StatusInternalServerError
			r.Exit()
		}
	}
	r.Response.Status = http.StatusOK
	r.Response.Write(wechatPaySuccessResponse{
		Code:    "SUCCESS",
		Message: "成功",
	})
	r.Exit()
}

type wechatPaySuccessResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type wechatPaySuccessNotifyContent struct {
	TransactionId string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid           string `json:"mchid"`
	TradeState      string `json:"trade_state"`
	BankType        string `json:"bank_type"`
	PromotionDetail []struct {
		Amount              int    `json:"amount"`
		WechatpayContribute int    `json:"wechatpay_contribute"`
		CouponId            string `json:"coupon_id"`
		Scope               string `json:"scope"`
		MerchantContribute  int    `json:"merchant_contribute"`
		Name                string `json:"name"`
		OtherContribute     int    `json:"other_contribute"`
		Currency            string `json:"currency"`
		StockId             string `json:"stock_id"`
		GoodsDetail         []struct {
			GoodsRemark    string `json:"goods_remark"`
			Quantity       int    `json:"quantity"`
			DiscountAmount int    `json:"discount_amount"`
			GoodsId        string `json:"goods_id"`
			UnitPrice      int    `json:"unit_price"`
		} `json:"goods_detail"`
	} `json:"promotion_detail"`
	SuccessTime *gtime.Time `json:"success_time"`
	Payer       struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	Appid          string `json:"appid"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach"`
	SceneInfo      struct {
		DeviceId string `json:"device_id"`
	} `json:"scene_info"`
}
