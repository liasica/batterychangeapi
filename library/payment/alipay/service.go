package alipay

import (
	"battery/app/model"
	"context"
	"errors"
	"github.com/gogf/gf/frame/g"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"strconv"
)

var serv *service

type service struct {
}

func Service() *service {
	if serv == nil {
		serv = &service{}
	}
	return serv
}

func (s *service) client() *alipay.Client {
	client, err := alipay.New(
		g.Cfg().GetString("payment.alipay.appId"),
		g.Cfg().GetString("payment.alipay.privateKey"),
		true,
	)
	if err != nil {
		panic("alipay error")
	}
	return client
}

// App 发起支付
func (s *service) App(ctx context.Context, prepay model.Prepay) (string, error) {
	trade := alipay.TradeAppPay{}
	trade.NotifyURL = prepay.NotifyUrl
	trade.Subject = prepay.Description
	trade.OutTradeNo = prepay.No
	trade.TotalAmount = strconv.FormatFloat(prepay.Amount, 'f', 2, 64)
	return s.client().TradeAppPay(trade)
}

// GetTradeNotification 获取支付回调通知数据
func (s *service) GetTradeNotification(ctx context.Context, r *http.Request) (*alipay.TradeNotification, error) {
	return s.client().GetTradeNotification(r)
}

// Refund 退款
func (s *service) Refund(ctx context.Context, tradeNo, outTradeNo, outRequestNo, refundAmount, refundReason string) (string, error) {
	res, err := s.client().TradeRefund(alipay.TradeRefund{
		TradeNo:      tradeNo,
		OutTradeNo:   outTradeNo,
		OutRequestNo: outRequestNo,
		RefundAmount: refundAmount,
		RefundReason: refundReason,
	})
	if err != nil {
		return "", err
	}
	if res.IsSuccess() {
		return res.Content.TradeNo, nil
	}
	return "", errors.New("退款失败")
}
