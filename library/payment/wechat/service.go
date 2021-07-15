package wechat

import (
	"battery/app/model"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"io/ioutil"
	"log"
	"net/http"
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

func (s *service) client() *core.Client {
	pkStr, err := ioutil.ReadFile(g.Cfg().GetString("payment.wechat.pkFile"))
	if err != nil {
		panic(fmt.Sprintf("wxpay err : %s", err))
	}
	block, _ := pem.Decode(pkStr)
	privateKey, _err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if _err != nil {
		panic(fmt.Sprintf("wxpay err : %s", _err))
	}
	var (
		mchID                      = g.Cfg().GetString("payment.wechat.mchId")    // 商户号
		mchCertificateSerialNumber = g.Cfg().GetString("payment.wechat.serialNo") // 商户证书序列号
		mchPrivateKey              = privateKey.(*rsa.PrivateKey)                 // 商户私钥
		mchAPIv3Key                = g.Cfg().GetString("payment.wechat.apiV3Key") // 商户APIv3密钥
	)
	ctx := context.Background()
	opts := []core.ClientOption{
		// 一次性设置 签名/验签/敏感字段加解密，并注册 平台证书下载器，自动定时获取最新的平台证书
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	fmt.Println(123123, client, err)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
		return nil
	}
	return client
}

// App 发起支付
func (s *service) App(ctx context.Context, req model.Prepay) (resp *app.PrepayResponse, err error) {
	a := app.AppApiService{
		Client: s.client(),
	}
	resp, _, err = a.Prepay(ctx, app.PrepayRequest{
		Appid:       core.String(g.Cfg().GetString("payment.wechat.appId")),
		Mchid:       core.String(g.Cfg().GetString("payment.wechat.mchId")),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.No),
		Amount: &app.Amount{
			Total: core.Int64(decimal.NewFromFloat(req.Amount).Mul(decimal.NewFromInt(100)).IntPart()),
		},
		NotifyUrl: core.String(req.NotifyUrl),
	})
	return resp, err
}

// ParseNotify 解析微信通知数据
func (s *service) ParseNotify(ctx context.Context, request *http.Request, content interface{}) (*notify.Request, error) {
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(g.Cfg().GetString("payment.wechat.mchId"))
	handler := notify.NewNotifyHandler(g.Cfg().GetString("payment.wechat.apiV3Key"), verifiers.NewSHA256WithRSAVerifier(certVisitor))
	return handler.ParseNotifyRequest(ctx, request, content)
}

// App 发起退款
func (s *service) Refund(ctx context.Context, transactionId, outTradeNo, outRefundNo, reason string, refundAmount, orderAmount float64) (string, error) {
	r := refunddomestic.RefundsApiService{
		Client: s.client(),
	}
	resp, res, err := r.Create(ctx, refunddomestic.CreateRequest{
		TransactionId: core.String(transactionId),
		OutTradeNo:    core.String(outTradeNo),
		OutRefundNo:   core.String(outRefundNo),
		Reason:        core.String(reason),
		Amount: &refunddomestic.AmountReq{
			Refund: core.Int64(int64(refundAmount * 100)),
			Total:  core.Int64(int64(orderAmount * 100)),
		},
	})
	if err != nil {
		return "", err
	}
	if resp.Status == refunddomestic.STATUS_SUCCESS.Ptr() || resp.Status == refunddomestic.STATUS_PROCESSING.Ptr() {
		return *resp.RefundId, nil
	} else {
		g.Log().Error("退款失败:", res, resp)
		return "", errors.New("退款失败")
	}
}
