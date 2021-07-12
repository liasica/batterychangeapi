package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dySmsApi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	_ "github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"math/rand"
	"time"
)

var SmsServer = smsService{
	accessKeyId:     g.Cfg().GetString("sms.accessKeyId"),
	accessKeySecret: g.Cfg().GetString("sms.accessKeySecret"),
}

type smsService struct {
	accessKeyId     string
	accessKeySecret string
}

// CreateAliClient  使用AK&SK初始化账号Client
func (s *smsService) CreateAliClient() (_result *dySmsApi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     &s.accessKeyId,
		AccessKeySecret: &s.accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dySmsApi20170525.Client{}
	_result, _err = dySmsApi20170525.NewClient(config)
	return _result, _err
}

// Send 短信发送
func (s *smsService) Send(ctx context.Context, req model.SmsSendReq) error {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	client, err := s.CreateAliClient()
	if err != nil {
		return err
	}
	res, err := client.SendSms(&dySmsApi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(req.Mobile),
		SignName:      tea.String("时光驹"),
		TemplateCode:  tea.String("SMS_206600136"),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	})
	if res != nil && *res.Body.Code == "OK" && err == nil {
		_, err = dao.Sms.Ctx(ctx).Insert(g.Map{
			dao.Sms.Columns.Mobile: req.Mobile,
			dao.Sms.Columns.Code:   code,
		})
	}
	return err
}

// Verify 短信验证
func (s *smsService) Verify(ctx context.Context, req model.SmsVerifyReq) bool {
	var sms model.Sms
	if err := dao.Sms.Ctx(ctx).Where(dao.Sms.Columns.Mobile, req.Mobile).OrderDesc(dao.Sms.Columns.Id).Limit(1).Scan(&sms); err != nil {
		fmt.Println(req, err)
		return false
	}
	_, _ = dao.Sms.Ctx(ctx).Where(dao.Sms.Columns.Mobile, sms.Mobile).Delete()       //直接删除已经使用的短信
	if sms.Code == req.Code && sms.CreatedAt.Add(time.Minute*2).After(gtime.Now()) { //两分钟过期时间
		return true
	}
	return false
}
