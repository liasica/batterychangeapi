package alipay

import (
	"battery/app/model"
	"context"
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
		g.Cfg().GetBool("payment.alipay.isProd"))
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
