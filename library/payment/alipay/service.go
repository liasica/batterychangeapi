package alipay

import (
	"battery/app/model"
	"context"
	"errors"
	"github.com/gogf/gf/frame/g"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"path/filepath"
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
	appId := g.Cfg().GetString("payment.alipay.appId")
	client, err := alipay.New(
		appId,
		g.Cfg().GetString("payment.alipay.privateKey"),
		true,
	)
	if err != nil {
		panic("[alipay] client error")
	}
	// 加载公钥证书
	p := filepath.Join("config", "alipay", appId)
	err = client.LoadAppPublicCertFromFile(p + "/appCertPublicKey_2021002155655488.cer") // 加载应用公钥证书
	if err != nil {
		panic("[alipay] LoadAppPublicCertFromFile error")
	}
	err = client.LoadAliPayRootCertFromFile(p + "/alipayRootCert.cer")                   // 加载支付宝根证书
	if err != nil {
		panic("[alipay] LoadAliPayRootCertFromFile error")
	}
	err = client.LoadAliPayPublicCertFromFile(p + "/alipayCertPublicKey_RSA2.cer")       // 加载支付宝公钥证书
	if err != nil {
		panic("[alipay] LoadAliPayPublicCertFromFile error")
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
