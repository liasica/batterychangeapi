package alipay

import (
	"battery/app/model"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"path/filepath"
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

// client 初始化支付宝
// params: 0是否加载公钥证书
func (s *service) client(params ...interface{}) *alipay.Client {
	appId := g.Cfg().GetString("payment.alipay.appId")
	client, err := alipay.New(
		appId,
		g.Cfg().GetString("payment.alipay.privateKey"),
		true,
	)
	if err != nil {
		panic("[alipay] client error")
	}
	if len(params) > 0 {
		if l := params[0].(bool); l {
			// 加载公钥证书
			p := filepath.Join("config", "alipay", appId)
			_ = client.LoadAppPublicCertFromFile(p + "/appCertPublicKey_2021002155655488.cer") // 加载应用公钥证书
			_ = client.LoadAliPayRootCertFromFile(p + "/alipayRootCert.cer")                   // 加载支付宝根证书
			_ = client.LoadAliPayPublicCertFromFile(p + "/alipayCertPublicKey_RSA2.cer")       // 加载支付宝公钥证书
		}
	}
	return client
}

// App 发起支付
func (s *service) App(ctx context.Context, prepay model.Prepay) (string, error) {
	trade := alipay.TradeAppPay{}
	trade.NotifyURL = prepay.NotifyUrl
	trade.Subject = prepay.Description
	trade.OutTradeNo = prepay.No
	trade.TotalAmount = fmt.Sprintf("%.2f", prepay.Amount)
	return s.client().TradeAppPay(trade)
}

// GetTradeNotification 获取支付回调通知数据
func (s *service) GetTradeNotification(ctx context.Context, r *http.Request) (*alipay.TradeNotification, error) {
	return s.client(true).GetTradeNotification(r)
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
	g.Log().Info("阿里退款响应：", res, err)
	if err != nil {
		return "", err
	}
	if res.IsSuccess() {
		return res.Content.TradeNo, nil
	}
	return "", errors.New("退款失败")
}
